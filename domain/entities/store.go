package entities

import "time"

type Store struct {
	StoreID       string    `json:"store_id" bson:"store_id"`
	SellerID      string    `json:"seller_id" bson:"seller_id"`
	WarehouseID   string    `json:"warehouse_id" bson:"warehouse_id"`
	StoreName     string    `json:"store_name" bson:"store_name"`
	StoreAddress  string    `json:"store_address" bson:"store_address"`
	StoreContact  string    `json:"store_contact" bson:"store_contact"`
	NumberOfRacks int       `json:"number_of_racks" bson:"number_of_racks"`
	OccupiedRacks int       `json:"occupied_racks" bson:"occupied_racks"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

// Request structures
type GetAllStoresRequest struct {
	WarehouseID string `json:"warehouse_id" binding:"required"`
	Limit       int64  `json:"limit"`
	Offset      int64  `json:"offset"`
	Search      string `json:"search"`
}

type GetStoreByIdRequest struct {
	StoreID string `json:"store_id" binding:"required"`
}

type CreateStoreRequest struct {
	SellerID      string `json:"seller_id" binding:"required"`
	WarehouseID   string `json:"warehouse_id" binding:"required"`
	StoreName     string `json:"store_name" binding:"required"`
	StoreAddress  string `json:"store_address" binding:"required"`
	StoreContact  string `json:"store_contact" binding:"required"`
	NumberOfRacks int    `json:"number_of_racks" binding:"required,min=1"`
}

type UpdateStoreRequest struct {
	StoreName     string `json:"store_name"`
	StoreAddress  string `json:"store_address"`
	StoreContact  string `json:"store_contact"`
	NumberOfRacks int    `json:"number_of_racks" binding:"min=1"`
	OccupiedRacks int    `json:"occupied_racks" binding:"min=0"`
}

type DeleteStoreRequest struct {
	StoreID string `json:"store_id" binding:"required"`
}

type GetStoreBySellerIdRequest struct {
	SellerID string `json:"seller_id" binding:"required"`
}

// Response structures
type GetAllStoresResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Stores  []Store `json:"stores"`
	Total   int64   `json:"total"`
	Limit   int64   `json:"limit"`
	Offset  int64   `json:"offset"`
}

type GetStoreResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Store   Store  `json:"store"`
}

type CreateStoreResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Store   Store  `json:"store"`
}

type UpdateStoreResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Store   Store  `json:"store"`
}

type DeleteStoreResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetStoreBySellerIdResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Store   Store  `json:"store"`
}

// Paginated response
type PaginatedStoresResponse struct {
	Stores      []Store `json:"stores"`
	Total       int64   `json:"total"`
	Limit       int64   `json:"limit"`
	Offset      int64   `json:"offset"`
	HasNext     bool    `json:"has_next"`
	HasPrevious bool    `json:"has_previous"`
}
