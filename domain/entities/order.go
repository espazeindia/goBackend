package entities

import "time"

//db me datatype
type Orders struct {
	OrderID     string    `json:"order_id"   bson:"order_id"`
	UserID      string    `json:"user_id" bson:"user_id"`
	WarehouseID string    `json:"warehouse_id" bson:"warehouse_id"`
	Address     string    `json:"address"   bson:"address"`
	OrderTotal  int       `json:"order_total"  bson:"order_total"`
	OrderedAt   time.Time `json:"ordered_at"  bson:"ordered_at"`
}

type OrderedItems struct {
	OrderID   string `json:"order_id" bson:"order_id"`
	ProductID string `json:"product_id"  bson:"product_id"`
	Quantity  int    `json:"quantity"  bson:"quantity"`
	Price     int    `json:"price"  bson:"price"`
	MRP       int    `json:"mrp"  bson:"mrp"`
	SellerID  string `json:"seller_id"  bson:"seller_id"`
}

// requests and respone types

type GetAllOrdersRequest struct {
	Limit  int `json:"limit" binding:"required,gte=1"`
	Offset int `json:"offset" binding:"gte=0"`
}

type GetAllOrdersReturn struct {
	OrderID     string          `json:"order_id"`
	UserID      string          `json:"user_id"`
	WarehouseID string          `json:"warehouse_id"`
	Address     string          `json:"address"`
	OrderTotal  int             `json:"order_total"`
	OrderedAt   time.Time       `json:"ordered_at"`
	Products    []*OrderedItems `json:"products"`
}

type GetAllOrderPaginated struct {
	Orders      []*GetAllOrdersReturn `json:"orders"`
	Total       int                   `json:"total"`
	Limit       int64                 `json:"limit"`
	Offset      int64                 `json:"offset"`
	HasNext     bool                  `json:"has_next"`
	HasPrevious bool                  `json:"has_previous"`
}

type CreateOrderRequest struct {
	UserID      string `json:"user_id"`
	WarehouseID string `json:"warehouse_id"`
	Address     string `json:"address"`
	OrderTotal  int    `json:"order_total"`
	Products    []*struct {
		ProductID string `json:"product_id"  bson:"product_id"`
		Quantity  int    `json:"quantity"  bson:"quantity"`
		Price     int    `json:"price"  bson:"price"`
		MRP       int    `json:"mrp"  bson:"mrp"`
		SellerID  string `json:"seller_id"  bson:"seller_id"`
	} `json:"products"`
}
