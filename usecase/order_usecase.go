package usecase

import (
	"context"
	"errors"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderUsecase struct {
	OrderRepository repositories.OrderRepository
}

func NewOrderUsecase(orderRepository repositories.OrderRepository) *OrderUsecase {
	return &OrderUsecase{OrderRepository: orderRepository}
}

func (u *OrderUsecase) GetAllOrders(ctx context.Context, requestData *entities.GetAllOrdersRequest) (*entities.GetAllOrderPaginated, error) {
	if requestData.Limit < 10 {
		requestData.Limit = 10
	}

	if requestData.Offset < 0 {
		requestData.Offset = 0
	}

	orders, total, err := u.OrderRepository.GetAllOrders(ctx, requestData)

	if err != nil {
		return nil, err
	}

	hasNext := requestData.Offset*requestData.Limit < total
	hasPrev := requestData.Offset > 0

	return &entities.GetAllOrderPaginated{
		Orders:      orders,
		Total:       total,
		Limit:       int64(requestData.Limit),
		Offset:      int64(requestData.Offset),
		HasNext:     hasNext,
		HasPrevious: hasPrev,
	}, nil
}

func (u *OrderUsecase) CreateNewOrder(ctx context.Context, requestOrder *entities.CreateOrderRequest) error {
	OrderId := primitive.NewObjectID().Hex()
	OrderedAt := time.Now()

	err := u.OrderRepository.CreateNewOrder(ctx, requestOrder, OrderId, OrderedAt)
	if err != nil {
		return err
	}

	err = u.OrderRepository.CreateNewOrderProducts(ctx, requestOrder, OrderId)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecase) GetOrderByOrderID(ctx context.Context, orderId *string) (*entities.GetAllOrdersReturn, error) {
	order, err := u.OrderRepository.GetOrderByOrderID(ctx, orderId)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no order found for this order ID")
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (u *OrderUsecase) GetOrderByUserID(ctx context.Context, userId *string) ([]*entities.GetAllOrdersReturn, error) {
	orders, err := u.OrderRepository.GetOrderByUserID(ctx, userId)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no order found for this user ID")
	}
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (u *OrderUsecase) GetOrderBySellerID(ctx context.Context, sellerId *string) ([]*entities.GetAllOrdersReturn, error) {
	orders, err := u.OrderRepository.GetOrderBySellerID(ctx, sellerId)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no order found for this seller ID")
	}
	if err != nil {
		return nil, err
	}
	return orders, nil
}
