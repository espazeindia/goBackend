package repositories

import (
	"context"

	"espazeBackend/domain/entities"
)

// WarehouseRepository defines the interface for warehouse data access operations
type WarehouseRepository interface {
	GetAllWarehouses(ctx context.Context) ([]*entities.Warehouse, error)
	GetWarehouseById(ctx context.Context, id string) (*entities.Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse *entities.Warehouse) error
	UpdateWarehouse(ctx context.Context, id string, warehouse *entities.Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error
}
