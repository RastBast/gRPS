package domain

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

type Status string

const (
	StatusCreated    Status = "created"
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusCreated, StatusNew, StatusProcessing, StatusCompleted, StatusCanceled:
		return true
	default:
		return false
	}
}

// Validate проверяет корректность полей заказа
func (o Order) Validate() error {
	if o.Total <= 0 {
		return fmt.Errorf("сумма заказа %d не может быть меньше или равна 0", o.Total)
	}

	// strings.TrimSpace защитит от любого количества пробелов и пустых строк именно его надо юзать
	if strings.TrimSpace(o.CustomerName) == "" {
		return fmt.Errorf("имя заказчика %q не может быть пустым или состоять только из пробелов", o.CustomerName)
	}

	return nil
}

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
