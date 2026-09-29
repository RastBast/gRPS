package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"

	pb "github.com/RastBast/grpc/gen/orders/v1"
	commonPg "github.com/RastBast/grpc/internal/common/postgres"
	order "github.com/RastBast/grpc/internal/order/domain"
	orderPg "github.com/RastBast/grpc/internal/order/repository/postgres"
	grpcServer "github.com/RastBast/grpc/internal/presentation/grpc"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, используются системные переменные")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Ошибка: переменная DATABASE_URL не задана в .env")
	}

	// 1. Создаем главный контекст приложения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Подключаемся к PostgreSQL через пул соединений
	pool, err := commonPg.NewPool(ctx)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer pool.Close()

	// 3. Инициализируем репозиторий для работы с заказами
	repo := orderPg.NewOrderRepository(pool)

	// 4. Создаем канал "Лента на кухню"
	orders := make(chan order.Order, 5)
	var wg sync.WaitGroup

	// 5. Запускаем воркера (Повара), который сохраняет заказы в Postgres
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case o, ok := <-orders:
				if !ok {
					log.Println("Лента остановлена, повар уходит домой")
					return
				}
				log.Printf("Повар начал обрабатывать и сохранять заказ %d\n", o.ID)

				// Меняем статус заказа на COMPLETED
				o.Status = order.StatusCompleted

				// Пишем результат напрямую в PostgreSQL
				if err := repo.SaveOrder(ctx, o); err != nil {
					log.Printf("Ошибка сохранения заказа %d в БД: %v\n", o.ID, err)
					continue
				}

				log.Printf("Заказ %d успешно обработан и сохранен в Postgres\n", o.ID)

			case <-ctx.Done():
				log.Println("Ресторан закрывается, повар завершает работу")
				return
			}
		}
	}()

	// 6. Настраиваем и регистрируем gRPC-сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось открыть порт 50051: %v", err)
	}

	gRPCServer := grpc.NewServer()
	orderServiceHandler := grpcServer.NewOrderGRPCServer(orders, repo)
	pb.RegisterOrderServiceServer(gRPCServer, orderServiceHandler)

	reflection.Register(gRPCServer)

	// Запускаем gRPC-сервер в отдельной горутине, чтобы не блокировать main
	go func() {
		log.Println("gRPC сервер запущен на порту :50051")
		if err := gRPCServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Ошибка работы gRPC сервера: %v", err)
		}
	}()

	// 7. Ожидаем сигнал остановки (Ctrl+C или SIGTERM от Docker/K8s)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("Получен сигнал завершения. Начинаем Graceful Shutdown...")

	// 8. Порядок остановки всех систем:
	gRPCServer.GracefulStop() // Перестаем принимать новые сетевые запросы
	cancel()                  // Уведомляем горутины об остановке
	close(orders)             // Закрываем канал
	wg.Wait()                 // Ждем завершения фоновой обработки заказа

	log.Println("Сервер успешно и корректно остановлен")
}
