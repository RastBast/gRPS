package order

import (
	"log"
	"sync/atomic"
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)

type LogReporter struct{}

type CounterReporter struct {
	Count int64
}

type Reporter interface {
	Report(o Order)
}

type Order struct {
	ID           int64
	CustomerName string
	Total        int64
	CreatedAt    time.Time
	Status       Status
}

func (c *CounterReporter) Report(o Order) {
	atomic.AddInt64(&c.Count, 1)
	log.Println(atomic.LoadInt64(&c.Count), "Номер")
}

func (l LogReporter) Report(o Order) {
	log.Println(o, "Интерфейс")
}
