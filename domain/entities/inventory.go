package entities

type GetAllInventoryRequest struct {
	SellerID string `json:"seller_id"`
	Limit    int64  `json:"limit"`
	Offset   int64  `json:"offset"`
	Search   string `json:"search"`
}

type GetAllInventoryResponse struct {
	InventoryId              string  `json:"inventory_id"`
	InventoryProductId       string  `json:"inventory_product_id"`
	MetadataProductId        string  `json:"metadata_product_id"`
	ProductVisibility        string  `json:"product_visibility"`
	MetadataName             string  `json:"metadata_name"`
	MetadataDescription      string  `json:"metadata_description"`
	MetadataImage            string  `json:"metadata_image"`
	MetadataCategoryId       string  `json:"metadata_category_id"`
	MetadataSubcategoryId    string  `json:"metadata_subcategory_id"`
	MetadataMrp              float64 `json:"metadata_mrp"`
	ProductQuantity          int     `json:"product_quantity"`
	ProductExpiryDate        string  `json:"product_expiry_date"`
	ProductManufacturingDate string  `json:"product_manufacturing_date"`
	MetadataCreatedAt        string  `json:"metadata_created_at"`
}

type Inventory struct {
	SellerID    string `json:"seller_id" bson:"seller_id"`
	InventoryID string `json:"inventory_id" bson:"inventory_id"`
}

type InventoryProduct struct {
	InventoryProductID       string `json:"inventory_product_id" bson:"inventory_product_id"`
	InventoryID              string `json:"inventory_id" bson:"inventory_id"`
	MetadataProductID        string `json:"metadata_product_id" bson:"metadata_product_id"`
	ProductVisibility        string `json:"product_visibility" bson:"product_visibility"`
	ProductQuantity          int    `json:"product_quantity" bson:"product_quantity"`
	ProductExpiryDate        string `json:"product_expiry_date" bson:"product_expiry_date"`
	ProductManufacturingDate string `json:"product_manufacturing_date" bson:"product_manufacturing_date"`
}

type AddInventoryRequest struct {
	SellerID          string `json:"seller_id"`
	InventoryProducts []struct {
		MetadataProductID        string `json:"metadata_product_id"`
		ProductVisibility        string `json:"product_visibility"`
		ProductQuantity          int    `json:"product_quantity"`
		ProductExpiryDate        string `json:"product_expiry_date"`
		ProductManufacturingDate string `json:"product_manufacturing_date"`
	} `json:"inventory_products"`
}

type UpdateInventoryRequest struct {
	InventoryID              string `json:"inventory_id"`
	InventoryProductID       string `json:"inventory_product_id"`
	ProductVisibility        string `json:"product_visibility"`
	ProductQuantity          int    `json:"product_quantity"`
	ProductExpiryDate        string `json:"product_expiry_date"`
	ProductManufacturingDate string `json:"product_manufacturing_date"`
}

type DeleteInventoryRequest struct {
	InventoryID        string `json:"inventory_id"`
	InventoryProductID string `json:"inventory_product_id"`
}

type GetInventoryByIdRequest struct {
	InventoryID string `json:"inventory_id"`
}

type GetInventoryByIdResponse struct {
	InventoryId              string  `json:"inventory_id"`
	InventoryProductId       string  `json:"inventory_product_id"`
	MetadataProductId        string  `json:"metadata_product_id"`
	ProductVisibility        string  `json:"product_visibility"`
	MetadataName             string  `json:"metadata_name"`
	MetadataDescription      string  `json:"metadata_description"`
	MetadataImage            string  `json:"metadata_image"`
	MetadataCategoryId       string  `json:"metadata_category_id"`
	MetadataSubcategoryId    string  `json:"metadata_subcategory_id"`
	MetadataMrp              float64 `json:"metadata_mrp"`
	ProductQuantity          int     `json:"product_quantity"`
	ProductExpiryDate        string  `json:"product_expiry_date"`
	ProductManufacturingDate string  `json:"product_manufacturing_date"`
	MetadataCreatedAt        string  `json:"metadata_created_at"`
}
