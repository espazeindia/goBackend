package repositories

import (
	"context"
	"espazeBackend/domain/entities"
	"time"
)

type OrderRepository interface {
	GetAllOrders(ctx context.Context, requestData *entities.GetAllOrdersRequest) ([]*entities.GetAllOrdersReturn, int, error)
	CreateNewOrder(ctx context.Context, requestOrder *entities.CreateOrderRequest, OrderID string, OrderedAt time.Time) error
	CreateNewOrderProducts(ctx context.Context, requestOrder *entities.CreateOrderRequest, OrderID string) error
	GetOrderByOrderID(ctx context.Context, orderId *string) (*entities.GetAllOrdersReturn, error)
	GetOrderByUserID(ctx context.Context, userId *string) ([]*entities.GetAllOrdersReturn, error)
	GetOrderBySellerID(ctx context.Context, sellerId *string) ([]*entities.GetAllOrdersReturn, error)
}
