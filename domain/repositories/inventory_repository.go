package repositories

import (
	"context"
	"espazeBackend/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryRepository interface {
	GetAllInventory(ctx context.Context, inventoryRequest entities.GetAllInventoryRequest) ([]entities.GetAllInventoryResponse, error)
	CreateInventory(ctx context.Context, inventoryRequest entities.AddInventoryRequest) (string, error)
	CreateInventoryProduct(ctx context.Context, inventoryId, inventoryProductId primitive.ObjectID, inventoryRequest entities.AddInventoryRequest) error
	UpdateInventory(ctx context.Context, inventoryRequest entities.UpdateInventoryRequest) error
	DeleteInventory(ctx context.Context, inventoryRequest entities.DeleteInventoryRequest) error
	GetInventoryById(ctx context.Context, inventoryRequest entities.GetInventoryByIdRequest) (*entities.GetInventoryByIdResponse, error)
}
