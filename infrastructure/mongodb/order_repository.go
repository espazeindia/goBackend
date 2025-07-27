package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepositoryMongoDB struct {
	Database *mongo.Database
}

func NewOrderRepositoryMongoDB(database *mongo.Database) repositories.OrderRepository {
	return &OrderRepositoryMongoDB{Database: database}
}

func (r *OrderRepositoryMongoDB) GetAllOrders(ctx context.Context, requestData *entities.GetAllOrdersRequest) ([]*entities.GetAllOrdersReturn, int, error) {
	orderCollection := r.Database.Collection("order")
	orderedItemCollection := r.Database.Collection("orderedItems")

	filter := bson.M{}
	total, err := orderCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	option := options.Find().SetLimit(int64(requestData.Limit)).SetSkip(int64(requestData.Offset)).SetSort(bson.D{{Key: "ordered_at", Value: -1}})
	cursor, err := orderCollection.Find(ctx, filter, option)
	if err != nil {
		return nil, 0, err
	}

	var orderDetail []*entities.Orders

	err = cursor.All(ctx, &orderDetail)
	if err != nil {
		return nil, 0, err
	}

	var result []*entities.GetAllOrdersReturn

	for _, order := range orderDetail {
		var allProducts []*entities.OrderedItems

		filter := bson.M{"order_id": order.OrderID}
		cursor, err := orderedItemCollection.Find(ctx, filter)
		if err != nil {
			return nil, 0, err
		}
		err = cursor.All(ctx, &allProducts)
		if err != nil {
			return nil, 0, err
		}

		orderResult := &entities.GetAllOrdersReturn{
			OrderID:     order.OrderID,
			UserID:      order.UserID,
			WarehouseID: order.WarehouseID,
			Address:     order.Address,
			OrderTotal:  order.OrderTotal,
			OrderedAt:   order.OrderedAt,
			Products:    allProducts,
		}

		result = append(result, orderResult)
	}

	return result, int(total), nil

}

func (r *OrderRepositoryMongoDB) CreateNewOrder(ctx context.Context, requestOrder *entities.CreateOrderRequest, OrderID string, OrderedAt time.Time) error {
	orderCollection := r.Database.Collection("order")

	orderDetail := &entities.Orders{
		OrderID:     OrderID,
		UserID:      requestOrder.UserID,
		WarehouseID: requestOrder.WarehouseID,
		Address:     requestOrder.Address,
		OrderTotal:  requestOrder.OrderTotal,
		OrderedAt:   OrderedAt,
	}

	_, err := orderCollection.InsertOne(ctx, orderDetail)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryMongoDB) CreateNewOrderProducts(ctx context.Context, requestOrder *entities.CreateOrderRequest, OrderID string) error {

	orderedItemCollection := r.Database.Collection("orderedItems")

	var allProducts []interface{}

	for _, product := range requestOrder.Products {

		orderProductDetail := &entities.OrderedItems{
			OrderID:   OrderID,
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
			Price:     product.Price,
			MRP:       product.MRP,
			SellerID:  product.SellerID,
		}
		allProducts = append(allProducts, orderProductDetail)
	}

	_, err := orderedItemCollection.InsertMany(ctx, allProducts)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryMongoDB) GetOrderByOrderID(ctx context.Context, orderId *string) (*entities.GetAllOrdersReturn, error) {

	orderCollection := r.Database.Collection("order")
	orderItemCollection := r.Database.Collection("orderedItems")

	filter := bson.M{"order_id": orderId}

	var order *entities.Orders

	err := orderCollection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	products, err := orderItemCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var allProducts []*entities.OrderedItems

	err = products.All(ctx, &allProducts)
	if err != nil {
		return nil, err
	}

	resultOrder := &entities.GetAllOrdersReturn{
		OrderID:     order.OrderID,
		UserID:      order.UserID,
		WarehouseID: order.WarehouseID,
		Address:     order.Address,
		OrderTotal:  order.OrderTotal,
		OrderedAt:   order.OrderedAt,
		Products:    allProducts,
	}

	return resultOrder, nil

}

func (r *OrderRepositoryMongoDB) GetOrderByUserID(ctx context.Context, userId *string) ([]*entities.GetAllOrdersReturn, error) {
	orderCollection := r.Database.Collection("order")
	orderItemCollection := r.Database.Collection("orderedItems")

	filter := bson.M{"user_id": userId}

	var orders []*entities.Orders

	ordersFound, err := orderCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = ordersFound.All(ctx, &orders)
	if err != nil {
		return nil, err
	}

	var result []*entities.GetAllOrdersReturn

	for _, order := range orders {
		var allProducts []*entities.OrderedItems

		filter := bson.M{"order_id": order.OrderID}
		cursor, err := orderItemCollection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, &allProducts)
		if err != nil {
			return nil, err
		}

		orderResult := &entities.GetAllOrdersReturn{
			OrderID:     order.OrderID,
			UserID:      order.UserID,
			WarehouseID: order.WarehouseID,
			Address:     order.Address,
			OrderTotal:  order.OrderTotal,
			OrderedAt:   order.OrderedAt,
			Products:    allProducts,
		}

		result = append(result, orderResult)
	}

	return result, nil
}

func (r *OrderRepositoryMongoDB) GetOrderBySellerID(ctx context.Context, sellerId *string) ([]*entities.GetAllOrdersReturn, error) {
	orderCollection := r.Database.Collection("order")
	orderedItemCollection := r.Database.Collection("orderedItems")

	// Step 1: Find all items sold by the seller
	filter := bson.M{"seller_id": *sellerId}

	var sellerItems []entities.OrderedItems
	cursor, err := orderedItemCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &sellerItems); err != nil {
		return nil, err
	}

	// Step 2: Group items by order ID
	orderItemMap := make(map[string][]*entities.OrderedItems)
	for _, item := range sellerItems {
		orderID := item.OrderID
		orderItemMap[orderID] = append(orderItemMap[orderID], &item)
	}

	// Step 3: Fetch order details in bulk
	orderIDs := make([]string, 0, len(orderItemMap))
	for orderID := range orderItemMap {
		orderIDs = append(orderIDs, orderID)
	}

	orderCursor, err := orderCollection.Find(ctx, bson.M{"order_id": bson.M{"$in": orderIDs}})
	if err != nil {
		return nil, err
	}

	var orders []entities.Orders
	if err = orderCursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	// Step 4: Construct final result
	var result []*entities.GetAllOrdersReturn
	for _, order := range orders {
		orderProducts := orderItemMap[order.OrderID]
		result = append(result, &entities.GetAllOrdersReturn{
			OrderID:     order.OrderID,
			UserID:      order.UserID,
			WarehouseID: order.WarehouseID,
			Address:     order.Address,
			OrderTotal:  order.OrderTotal,
			OrderedAt:   order.OrderedAt,
			Products:    orderProducts,
		})
	}

	return result, nil
}
