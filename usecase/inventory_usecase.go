package usecase

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
)

type InventoryUseCaseInterface struct {
	inventoryRepo repositories.InventoryRepository
}

func NewInventoryUseCase(inventoryRepo repositories.InventoryRepository) *InventoryUseCaseInterface {
	return &InventoryUseCaseInterface{inventoryRepo: inventoryRepo}
}

func (u *InventoryUseCaseInterface) GetAllInventory(ctx context.Context, seller_id string, offset, limit int64, search, sort string) (*entities.PaginatedInventoryResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	inventory, total, err := u.inventoryRepo.GetAllInventory(ctx, seller_id, offset, limit, search, sort)
	if err != nil {
		return nil, err
	}
	var totalPages int64 = (total + limit - 1) / limit

	return &entities.PaginatedInventoryResponse{
		InventoryProduct: inventory,
		Total:            total,
		Limit:            limit,
		Offset:           offset,
		TotalPages:       totalPages,
	}, nil
}

func (u *InventoryUseCaseInterface) AddInventory(ctx context.Context, inventoryRequest *entities.AddInventoryRequest) (*entities.MessageResponse, error) {
	return u.inventoryRepo.CreateInventory(ctx, inventoryRequest)

}

func (u *InventoryUseCaseInterface) UpdateInventory(ctx context.Context, inventoryRequest entities.UpdateInventoryRequest) (*entities.MessageResponse, error) {
	return u.inventoryRepo.UpdateInventory(ctx, inventoryRequest)

}

func (u *InventoryUseCaseInterface) DeleteInventory(ctx context.Context, inventoryRequest entities.DeleteInventoryRequest) error {
	err := u.inventoryRepo.DeleteInventory(ctx, inventoryRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *InventoryUseCaseInterface) GetInventoryById(ctx context.Context, inventoryRequest string) (*entities.GetInventoryByIdResponse, error) {
	return u.inventoryRepo.GetInventoryById(ctx, inventoryRequest)

}
