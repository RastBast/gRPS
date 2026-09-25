package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/RastBast/grpc/gen/orders/v1"
	order "github.com/RastBast/grpc/internal/order/domain"
	server "github.com/RastBast/grpc/internal/presentation/grpc"
)

// 1. НАШ gRPC СЕРВИС
// Это наш "Официант". Ему нужен доступ к каналу, чтобы передавать заказы на кухню.
type OrderGRPCServer struct {
	pb.UnimplementedOrderServiceServer
	ordersChan chan<- order.Order // Канал для отправки заказов повару
}

// 2. МЕТОД ПРИЕМА ЗАКАЗА ПО СЕТИ
// Когда кто-то вызывает метод CreateOrder через интернет, срабатывает эта функция
func (s *OrderGRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Создаем заказ в нашем внутреннем формате
	newOrder := order.Order{
		ID:           time.Now().UnixMilli(), // Временный ID
		CustomerName: req.GetCustomerName(),
		Total:        req.GetTotal(),
		Status:       order.StatusNew,
		CreatedAt:    time.Now(),
	}

	// Кладем заказ на ленту (в канал)
	s.ordersChan <- newOrder
	log.Printf("Принят новый заказ по сети от: %s\n", req.GetCustomerName())

	// Отвечаем клиенту (всё ок)
	return &pb.CreateOrderResponse{
		//TODO: вернуть рельный заказ
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. СОЗДАЕМ КАНАЛ (Наша "Лента на кухню")
	orders := make(chan order.Order, 5)
	var wg sync.WaitGroup

	// 4. ЗАПУСКАЕМ ПОВАРА (Потребитель)
	wg.Add(1)
	counter := &order.CounterReporter{}
	ypo := order.LogReporter{}

	go func() {
		defer wg.Done()
		for {
			select {
			case o, ok := <-orders: // Берем заказ с ленты
				if !ok {
					log.Println("Лента остановлена, повар уходит домой")
					return
				}
				log.Printf("Повар начал готовить заказ %d\n", o.ID)

				// Эмулируем время готовки
				time.Sleep(500 * time.Millisecond)
				o.Status = order.StatusCompleted
				ypo.Report(o)
				counter.Report(o)

			case <-ctx.Done():
				log.Println("Ресторан закрывается, повар уходит")
				return
			}
		}
	}()

	// 5. НАНИМАЕМ ОФИЦИАНТА И ЗАПУСКАЕМ СЕРВЕР
	myOrderService := &OrderGRPCServer{
		ordersChan: orders,
	}

	// Вызываем функцию из твоего server.go
	grpcServer := server.SetupAndRunGRPCServer(myOrderService)

	// 6. ЖДЕМ СИГНАЛА НА ВЫКЛЮЧЕНИЕ (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh // Программа замирает здесь и работает, пока не нажмут Ctrl+C

	log.Println("Галя, у нас отмена! Закрываем ресторан...")

	cancel()                  // 1. Говорим всем горутинам, что пора закругляться
	grpcServer.GracefulStop() // 2. Перестаем принимать новые gRPC запросы (Официант уходит)
	close(orders)             // 3. Останавливаем ленту (больше заказов не будет)
	wg.Wait()                 // 4. Ждем, пока повар доготовит последний заказ

	log.Println("Корректное завершение работы выполнено")
}
