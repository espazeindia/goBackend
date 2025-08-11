package entities

import "time"

type GetProductsForSpecificStoreRequest struct {
	StoreID string `json:"store_id"`
}

type GetProductsForSpecificStoreResponse struct {
	StoreProducts []struct {
		InventoryId              string    `json:"inventory_id"`
		InventoryProductId       string    `json:"inventory_product_id"`
		MetadataProductId        string    `json:"metadata_product_id"`
		ProductVisibility        bool      `json:"product_visibility"`
		ProductPrice             float64   `json:"product_price"`
		MetadataName             string    `json:"metadata_name"`
		MetadataDescription      string    `json:"metadata_description"`
		MetadataImage            string    `json:"metadata_image"`
		MetadataCategoryId       string    `json:"metadata_category_id"`
		MetadataSubcategoryId    string    `json:"metadata_subcategory_id"`
		MetadataMrp              float64   `json:"metadata_mrp"`
		ProductQuantity          int       `json:"product_quantity"`
		ProductExpiryDate        time.Time `json:"product_expiry_date"`
		ProductManufacturingDate time.Time `json:"product_manufacturing_date"`
	}
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
