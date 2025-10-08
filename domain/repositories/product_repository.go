package repositories

import (
	"context"

	"espazeBackend/domain/entities"
)

// ProductRepository defines the interface for product data access operations
type ProductRepository interface {
	GetProductsForSpecificStore(ctx context.Context, store_id string) ([]*entities.GetProductsForSpecificStoreResponse, error)
	GetProductsForAllStores(ctx context.Context, allStores *[]entities.Store) ([]*entities.GetProductsForAllStoresResponse, error)
	FetchSellerId(ctx context.Context, storeID string) (string, error)
	GetAllStores(ctx context.Context, warehouseID string) (*[]entities.Store, error)
	GetProductsForStoreSubcategory(ctx context.Context, storeId, subcategoryId string) ([]*entities.GetProductsForStoreSubcategory, error)
	GetProductsForAllStoresSubcategory(ctx context.Context, warehouseId, subcategoryId string) ([]*entities.GetProductsForStoreSubcategory, error)
	GetBasicDetailsForProduct(ctx context.Context, inventoryProductID string) (*entities.GetBasicDetailsForProductResponse, error)
	GetProductComparisonByStore(ctx context.Context, warehouse_id string, inventoryProductID string) ([]*entities.GetProductComparisonByStoreResult, error)
}
