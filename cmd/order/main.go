package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	order "github.com/RastBast/GRPS/internal/order/domain"
)

type OrderGRPCServer struct {
	// pb.UnimplementedOrderServiceServer
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ловим Ctrl+C для graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Галя у нас отмена")
		cancel()
	}()

	orders := make(chan order.Order, 5) // создаем очередь ордеров (буфер 5)

	var wg sync.WaitGroup // Создаем счетчик горутин

	wg.Add(2) // добавляем 2 рутины в очередь

	// *producer*
	go func() { // создаем горутину чтобы получать ордера
		defer wg.Done()
		defer close(orders) // закрываем канал когда producer вышел (и по cancel тоже)

		for i := 0; i < 100; i++ {
			// ИСПРАВЛЕНО: добавлена открывающая скобка { и префиксы order. у статусов
			ord := order.Order{
				ID:           int64(i + 1),
				CustomerName: fmt.Sprintf("client-%d", i),
				Total:        1000 * int64(i+1),
				Status:       order.StatusNew,
				CreatedAt:    time.Now(),
			}
			select {
			case <-ctx.Done(): // Проверяем, не прервали ли программу
				log.Println("Продюсер остановлен сигналом отмены")
				return
			case orders <- ord: // ИСПРАВЛЕНО: отправляем переменную ord (имя изменено, чтобы не конфликтовать с именем пакета)
			}
		}
	}()

	counter := &order.CounterReporter{}
	var rep order.Reporter = counter

	ypo := order.LogReporter{}
	var i order.Reporter = ypo
	go func() { // это потребитель он получает ордера через канал а потом их читает
		defer wg.Done()
		for {
			select {
			case o, ok := <-orders: // Перебираем канал чтобы достать новые ордера
				if !ok {
					log.Println("канал закрыт, потребитель завершил работу")
					return
				}
				o.Status = order.StatusCompleted // ИСПРАВЛЕНО: добавлен префикс пакета order.
				i.Report(o)
				rep.Report(o)
			case <-ctx.Done():
				log.Println("Потребитель прекратил обработку из-за завершения работы")
				return
			case <-time.After(1 * time.Second):
				log.Println("Время ожидания ордера истекло")
			}
		}
	}()
	wg.Wait() // Ждет пока 2 горутины закончат работу продюсер и потребитель

	log.Println("корректное завершение работы выполнено")
}
