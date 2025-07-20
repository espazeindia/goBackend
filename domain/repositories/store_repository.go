package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type StoreRepository interface {
	GetAllStores(ctx context.Context, request entities.GetAllStoresRequest) (*entities.PaginatedStoresResponse, error)
	GetStoreById(ctx context.Context, storeId string) (*entities.Store, error)
	CreateStore(ctx context.Context, store *entities.Store) error
	UpdateStore(ctx context.Context, storeId string, store *entities.Store) error
	DeleteStore(ctx context.Context, storeId string) error
	GetStoreBySellerId(ctx context.Context, sellerId string) (*entities.Store, error)
	GetStoresByWarehouseId(ctx context.Context, warehouseId string) ([]entities.Store, error)
	UpdateStoreRacks(ctx context.Context, storeId string, occupiedRacks int) error
}
