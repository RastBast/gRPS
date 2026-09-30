package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	order "github.com/RastBast/grpc/internal/order/domain"
)

// Ошибка оредр не найден
var ErrOrderNotFound = errors.New("order not found")

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder создаёт новый заказ. ID заказа генерирует сама PostgreSQL
// (колонка id объявлена как BIGSERIAL), поэтому в запросе его не передаём,
// а забираем сгенерированное значение через RETURNING.
func (r *OrderRepository) CreateOrder(ctx context.Context, o order.Order) (order.Order, error) {
	if !o.Status.IsValid() {
		return order.Order{}, fmt.Errorf("недопустимый статус заказа: %q", o.Status)
	}

	// ВАЖНО: текст запроса — константа, значения передаются ОТДЕЛЬНО через
	// конкатенацию строк с пользовательскими данными — это открывает SQL-инъекции.
	query := `
		INSERT INTO orders (customer_name, total, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`

	err := r.db.QueryRow(ctx, query, o.CustomerName, o.Total, string(o.Status)).
		Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		return order.Order{}, fmt.Errorf("ошибка создания заказа в БД: %w", err)
	}
	return o, nil
}

// UpdateOrderStatus обновляет статус уже существующего заказа по его ID.
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, o order.Order) error {
	if !o.Status.IsValid() {
		return fmt.Errorf("недопустимый статус заказа: %q", o.Status)
	}

	query := `
		UPDATE orders
		SET status = $2
		WHERE id = $1;
	`

	tag, err := r.db.Exec(ctx, query, o.ID, string(o.Status))
	if err != nil {
		return fmt.Errorf("ошибка обновления заказа в БД: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}

// GetOrderByID получает заказ из PostgreSQL по его ID.
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
