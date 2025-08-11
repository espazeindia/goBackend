package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type InventoryRepository interface {
	GetAllInventory(ctx context.Context, seller_id string, offset, limit int64, search, sort string) ([]entities.GetAllInventoryResponse, int64, error)
	CreateInventory(ctx context.Context, inventoryRequest *entities.AddInventoryRequest) (*entities.MessageResponse, error)
	UpdateInventory(ctx context.Context, inventoryRequest entities.UpdateInventoryRequest) error
	DeleteInventory(ctx context.Context, inventoryRequest entities.DeleteInventoryRequest) error
	GetInventoryById(ctx context.Context, inventoryRequest entities.GetInventoryByIdRequest) (*entities.GetInventoryByIdResponse, error)
}
