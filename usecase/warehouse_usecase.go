package usecase

import (
	"context"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WarehouseUseCase handles business logic for warehouse operations
type WarehouseUseCase struct {
	warehouseRepo repositories.WarehouseRepository
}

// NewWarehouseUseCase creates a new warehouse use case
func NewWarehouseUseCase(warehouseRepo repositories.WarehouseRepository) *WarehouseUseCase {
	return &WarehouseUseCase{
		warehouseRepo: warehouseRepo,
	}
}

func (u *WarehouseUseCase) GetAllWarehouses(ctx context.Context) (*entities.MessageResponse, error) {
	return u.warehouseRepo.GetAllWarehouses(ctx)
}

func (u *WarehouseUseCase) GetWarehouseById(ctx context.Context, id string) (*entities.Warehouse, error) {
	// Validate ObjectID format
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return u.warehouseRepo.GetWarehouseById(ctx, id)
}

func (u *WarehouseUseCase) CreateWarehouse(ctx context.Context, warehouse *entities.CreateWarehouseRequest) (*entities.MessageResponse, error) {
	return u.warehouseRepo.CreateWarehouse(ctx, warehouse)
}

func (u *WarehouseUseCase) UpdateWarehouse(ctx context.Context, id string, warehouse *entities.UpdateWarehouseRequest) (*entities.MessageResponse, error) {
	return u.warehouseRepo.UpdateWarehouse(ctx, id, warehouse)
}

func (u *WarehouseUseCase) DeleteWarehouse(ctx context.Context, id string) error {
	// Validate ObjectID format
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return u.warehouseRepo.DeleteWarehouse(ctx, id)
}
