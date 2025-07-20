package usecase

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryUseCaseInterface struct {
	inventoryRepo repositories.InventoryRepository
}

func NewInventoryUseCase(inventoryRepo repositories.InventoryRepository) *InventoryUseCaseInterface {
	return &InventoryUseCaseInterface{inventoryRepo: inventoryRepo}
}

func (u *InventoryUseCaseInterface) GetAllInventory(ctx context.Context, inventoryRequest entities.GetAllInventoryRequest) ([]entities.GetAllInventoryResponse, error) {
	inventory, err := u.inventoryRepo.GetAllInventory(ctx, inventoryRequest)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (u *InventoryUseCaseInterface) AddInventory(ctx context.Context, inventoryRequest entities.AddInventoryRequest) error {
	InventoryProductId := primitive.NewObjectID()

	inventoryId, err := u.inventoryRepo.CreateInventory(ctx, inventoryRequest)
	if err != nil {
		return err
	}
	var InventoryId primitive.ObjectID
	// Convert the returned inventory ID string back to ObjectID for CreateInventoryProduct
	if inventoryId != "" {
		InventoryId, err = primitive.ObjectIDFromHex(inventoryId)
		if err != nil {
			return err
		}
	}

	err = u.inventoryRepo.CreateInventoryProduct(ctx, InventoryId, InventoryProductId, inventoryRequest)
	if err != nil {
		return err
	}

	return nil
}

func (u *InventoryUseCaseInterface) UpdateInventory(ctx context.Context, inventoryRequest entities.UpdateInventoryRequest) error {
	err := u.inventoryRepo.UpdateInventory(ctx, inventoryRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *InventoryUseCaseInterface) DeleteInventory(ctx context.Context, inventoryRequest entities.DeleteInventoryRequest) error {
	err := u.inventoryRepo.DeleteInventory(ctx, inventoryRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *InventoryUseCaseInterface) GetInventoryById(ctx context.Context, inventoryRequest entities.GetInventoryByIdRequest) (*entities.GetInventoryByIdResponse, error) {
	inventory, err := u.inventoryRepo.GetInventoryById(ctx, inventoryRequest)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}
