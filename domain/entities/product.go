package entities

import "time"

type GetProductsForSpecificStoreRequest struct {
	StoreID string `json:"store_id"`
}

type GetProductsForSpecificStoreResponse struct {
	InventoryId              string    `json:"inventory_id"`
	InventoryProductId       string    `json:"inventory_product_id"`
	MetadataProductId        string    `json:"metadata_product_id"`
	ProductVisibility        bool      `json:"visibility"`
	ProductPrice             float64   `json:"price"`
	MetadataName             string    `json:"name"`
	MetadataDescription      string    `json:"description"`
	MetadataImage            string    `json:"image"`
	MetadataCategoryId       string    `json:"category_id"`
	MetadataSubcategoryId    string    `json:"subcategory_id"`
	MetadataMrp              float64   `json:"mrp"`
	ProductQuantity          int       `json:"quantity"`
	ProductExpiryDate        time.Time `json:"expiry_date"`
	ProductManufacturingDate time.Time `json:"manufacturing_date"`
	ProductCategoryName      string    `json:"category_name"`
	ProductSubCategoryName   string    `json:"subcategory_name"`
	TotalStars               string    `json:"TotalStars"`
	TotalReviews             string    `json:"TotalReviews"`
}

type GetProductsForAllStoresRequest struct {
	WarehouseID string `json:"warehouse_id"`
}

type GetProductsForAllStoresResponse struct {
	AllStoresProducts []struct {
		StoreID       string                                `json:"store_id"`
		StoreProducts []GetProductsForSpecificStoreResponse `json:"store_products"`
	}
}

type GetProductsForStoreSubcategory struct {
	MetadataProductId        string    `json:"metadata_product_id"`
	MetadataName             string    `json:"name"`
	MetadataDescription      string    `json:"description"`
	MetadataImage            string    `json:"image"`
	MetadataCategoryId       string    `json:"category_id"`
	MetadataSubcategoryId    string    `json:"subcategory_id"`
	MetadataMrp              float64   `json:"mrp"`
	ProductCategoryName      string    `json:"category_name"`
	ProductSubCategoryName   string    `json:"subcategory_name"`
	TotalStars               int       `json:"TotalStars"`
	TotalReviews             int       `json:"TotalReviews"`
	InventoryId              string    `json:"inventory_id"`
	InventoryProductId       string    `json:"inventory_product_id"`
	ProductPrice             float64   `json:"price"`
	ProductQuantity          int       `json:"quantity"`
	ProductExpiryDate        time.Time `json:"expiry_date"`
	ProductManufacturingDate time.Time `json:"manufacturing_date"`
	MetadataRating           float64   `json:"metadata_rating"`
	StoreName                string    `json:"store_name"`
}

type GetBasicDetailsForProductRequest struct {
	InventoryProductId string `json:"inventory_product_id"`
}
