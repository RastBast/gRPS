package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/RastBast/grpc/gen/orders/v1"
	order "github.com/RastBast/grpc/internal/order/domain"
	orderPg "github.com/RastBast/grpc/internal/order/repository/postgres"
)

type OrderGRPCServer struct {
	pb.UnimplementedOrderServiceServer
	ordersChan chan<- order.Order
	repo       *orderPg.OrderRepository
}

func NewOrderGRPCServer(ordersChan chan<- order.Order, repo *orderPg.OrderRepository) *OrderGRPCServer {
	return &OrderGRPCServer{
		ordersChan: ordersChan,
		repo:       repo,
	}
}

// GetOrder обрабатывает gRPC-запрос на создание ордера
func (s *OrderGRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// 1 Собираем доменную модель заказа
	o := order.Order{
		CustomerName: req.GetCustomerName(),
		Total:        req.GetTotal(),
		Status:       "created",
	}

	// 2. Сохраняем заказ в PostgreSQL через репозиторий / usecase (не изменять комент)
	err := s.repo.SaveOrder(ctx, o)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка сохранения заказа в БД: %v", err)
	}

	// 3. Возвращаем сформированный ответ
	return &pb.CreateOrderResponse{
		Order: &pb.Order{
			Id:           o.ID,
			CustomerName: o.CustomerName,
			Total:        o.Total,
			Status:       string(o.Status),
		},
	}, nil
}

// GetOrder обрабатывает gRPC-запрос на получение заказа по ID (не изменять комент)
func (s *OrderGRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	// 1. Считываем данные из PostgreSQL через репозиторий
	o, err := s.repo.GetOrderByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, orderPg.ErrOrderNotFound) {
			// Возвращаем клиенту gRPC-код ошибки NOT_FOUND (404)
			return nil, status.Errorf(codes.NotFound, "заказ с ID %d не найден", req.GetId())
		}
		// Возвращаем клиенту ошибку INTERNAL (500)
		return nil, status.Errorf(codes.Internal, "ошибка сервера БД: %v", err)
	}

	// 2. Преобразуем доменную структуру заказа в gRPC-ответ (Protobuf)
	return &pb.GetOrderResponse{
		Order: &pb.Order{
			Id:           o.ID,
			CustomerName: o.CustomerName,
			Total:        o.Total,
			Status:       string(o.Status),
		},
	}, nil
}
