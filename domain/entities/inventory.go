package entities

import "time"

type Inventory struct {
	InventoryID string `json:"id" bson:"_id,omitempty"`
	SellerID    string `json:"seller_id" bson:"seller_id"`
	StoreId     string `json:"store_id" bson:"store_id"`
}

type InventoryProduct struct {
	InventoryProductID       string    `json:"id" bson:"_id,omitempty"`
	InventoryID              string    `json:"inventory_id" bson:"inventory_id"`
	MetadataProductID        string    `json:"metadata_product_id" bson:"metadata_product_id"`
	ProductVisibility        bool      `json:"product_visibility" bson:"product_visibility"`
	ProductQuantity          int       `json:"product_quantity" bson:"product_quantity"`
	ProductPrice             float64   `json:"product_price" bson:"product_price"`
	ProductExpiryDate        time.Time `json:"product_expiry_date" bson:"product_expiry_date"`
	ProductManufacturingDate time.Time `json:"product_manufacturing_date" bson:"product_manufacturing_date"`
}
type GetAllInventoryRequest struct {
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
	Search string `json:"search"`
}

type GetAllInventoryResponse struct {
	InventoryId              string  `json:"inventory_id" bson:"inventory_id"`
	InventoryProductId       string  `json:"inventory_product_id" bson:"inventory_product_id"`
	MetadataProductId        string  `json:"metadata_product_id" bson:"metadata_product_id"`
	ProductVisibility        bool    `json:"product_visibility" bson:"product_visibility"`
	MetadataName             string  `json:"metadata_name" bson:"metadata_name"`
	MetadataDescription      string  `json:"metadata_description" bson:"metadata_description"`
	MetadataImage            string  `json:"metadata_image" bson:"metadata_image"`
	MetadataCategoryId       string  `json:"metadata_category_id" bson:"metadata_category_id"`
	MetadataSubcategoryId    string  `json:"metadata_subcategory_id" bson:"metadata_subcategory_id"`
	MetadataMrp              float64 `json:"metadata_mrp" bson:"metadata_mrp"`
	ProductQuantity          int     `json:"product_quantity" bson:"product_quantity"`
	ProductPrice             float64 `json:"product_price" bson:"product_price"`
	ProductExpiryDate        string  `json:"product_expiry_date" bson:"product_expiry_date"`
	ProductManufacturingDate string  `json:"product_manufacturing_date" bson:"product_manufacturing_date"`
	MetadataCreatedAt        string  `json:"metadata_created_at" bson:"metadata_created_at"`
	MetadataCategoryName     string  `json:"metadata_category_name" bson:"metadata_category_name"`
	MetadataSubcategoryName  string  `json:"metadata_subcategory_name" bson:"metadata_subcategory_name"`
	MetadataHSNCode          string  `json:"metadata_hsn_code" bson:"metadata_hsn_code"`
}

type PaginatedInventoryResponse struct {
	InventoryProduct []GetAllInventoryResponse `json:"inventory_products"`
	Total            int64                     `json:"total"`
	Limit            int64                     `json:"limit"`
	Offset           int64                     `json:"offset"`
	TotalPages       int64                     `json:"total_pages"`
}

type AddInventoryRequest struct {
	SellerID          string   `json:"seller_id" bson:"omitempty"`
	MetadataProductID []string `json:"metadata_ids"`
}

type UpdateInventoryRequest struct {
	InventoryProductID       string  `json:"inventory_product_id"`
	ProductVisibility        bool    `json:"product_visibility"`
	ProductQuantity          int     `json:"product_quantity"`
	ProductPrice             float64 `json:"product_price"`
	ProductExpiryDate        string  `json:"product_expiry_date"`
	ProductManufacturingDate string  `json:"product_manufacturing_date"`
}

type DeleteInventoryRequest struct {
	InventoryID        string `json:"inventory_id"`
	InventoryProductID string `json:"inventory_product_id"`
}

type GetInventoryByIdRequest struct {
	InventoryID string `json:"inventory_id"`
}

type GetInventoryByIdResponse struct {
	InventoryProductId       string    `json:"inventory_product_id"`
	MetadataProductId        string    `json:"metadata_product_id"`
	ProductVisibility        bool      `json:"product_visibility"`
	MetadataName             string    `json:"metadata_name"`
	MetadataDescription      string    `json:"metadata_description"`
	MetadataImage            string    `json:"metadata_image"`
	MetadataCategoryId       string    `json:"metadata_category_id"`
	MetadataSubcategoryId    string    `json:"metadata_subcategory_id"`
	MetadataMrp              float64   `json:"metadata_mrp"`
	ProductQuantity          int       `json:"product_quantity"`
	ProductPrice             float64   `json:"product_price"`
	ProductExpiryDate        time.Time `json:"product_expiry_date"`
	ProductManufacturingDate time.Time `json:"product_manufacturing_date"`
	MetadataCreatedAt        string    `json:"metadata_created_at"`
	MetadataCategoryName     string    `json:"metadata_category_name" bson:"metadata_category_name"`
	MetadataSubcategoryName  string    `json:"metadata_subcategory_name" bson:"metadata_subcategory_name"`
	MetadataHSNCode          string    `json:"metadata_hsn_code" bson:"metadata_hsn_code"`
	TotalReviews             int       `json:"metadata_total_reviews"`
	TotalStars               int       `json:"metadata_total_stars"`
}
