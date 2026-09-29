package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	order "github.com/RastBast/grpc/internal/order/domain"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// SaveOrder сохраняет или обновляет статус заказа в PostgreSQL
func (r *OrderRepository) SaveOrder(ctx context.Context, o order.Order) error {
	query := `
		INSERT INTO orders (id, customer_name, total, status)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE 
		SET status = EXCLUDED.status;
	`

	_, err := r.db.Exec(ctx, query, o.ID, o.CustomerName, o.Total, string(o.Status))
	if err != nil {
		return fmt.Errorf("ошибка сохранения заказа в БД: %w", err)
	}

	return nil
}

// GetOrderByID получает заказ из PostgreSQL по его ID
func (r *OrderRepository) GetOrderByID(ctx context.Context, id int64) (order.Order, error) {
	query := `
		SELECT id, customer_name, total, status 
		FROM orders 
		WHERE id = $1;
	`

	var o order.Order
	var statusStr string

	err := r.db.QueryRow(ctx, query, id).Scan(&o.ID, &o.CustomerName, &o.Total, &statusStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order.Order{}, ErrOrderNotFound
		}
		return order.Order{}, fmt.Errorf("ошибка чтения заказа из БД: %w", err)
	}

	o.Status = order.Status(statusStr)
	return o, nil
}
