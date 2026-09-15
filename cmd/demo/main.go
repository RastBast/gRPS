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
)

type Status string // создаем статус

const (
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)

type Order struct { // струтора заказа
	ID           int64
	CustomerName string
	Total        int64
	CreatedAt    time.Time
	Status       Status
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

	orders := make(chan Order, 5) // создаем очередь горутин (буфер 5)

	var wg sync.WaitGroup // Создаем счечик го рутин

	wg.Add(2) // добовляем 2 рутины в очередь (фикс: было 3, а горутин 2)

	// *producer*
	go func() { // создаем го рутину что бы получать ордера
		defer wg.Done()
		defer close(orders) // закрваем канал когда producer вышел (и по cancel тоже)

		for i := 0; i < 10; i++ {
			order := Order{
				ID:           int64(i + 1),
				CustomerName: fmt.Sprintf("client-%d", i),
				Total:        1000 * int64(i+1),
				Status:       StatusNew,
				CreatedAt:    time.Now(),
			}
			select {
			case <-ctx.Done(): // Проверяем, не прервали ли программу
				log.Println("Продюсер остановлен сигналом отмены")
				return
			case orders <- order: // Типо блатные каждый раз меняем данные что бы ордер был не пустой
			}
		}
	}()

	go func() { // это потребилеь он получает оредра через канал а потом их читает
		defer wg.Done()
		for {
			select {
			case o, ok := <-orders: // Переберям канал что бы достань новые оредра
				if !ok {
					log.Println("канал закрыт, потребитель завершил работу")
					return
				}
				o.Status = StatusCompleted         // меняем стаусу заказа
				time.Sleep(100 * time.Millisecond) // Печатам после 100 мс
				log.Println(o)
			case <-ctx.Done():
				log.Println("Потребитель прекратил обработку из-за завершения работы")
				return
			case <-time.After(1 * time.Second):
				log.Println("Время ожидания ордера истекло")
			}
		}
	}()
	wg.Wait() // Ждет пока 2 го рутины закчант работу продюсер и потребитель
	log.Println("корректное завершение работы выполнено")
}
