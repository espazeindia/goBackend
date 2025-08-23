package usecase

import (
	"context"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductUseCase handles business logic for product operations
type ProductUseCase struct {
	productRepo repositories.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repositories.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

func (u *ProductUseCase) GetProductsForSpecificStore(ctx context.Context, store_id string) ([]*entities.GetProductsForSpecificStoreResponse, error) {

	return u.productRepo.GetProductsForSpecificStore(ctx, store_id)

}

func (u *ProductUseCase) GetProductsForAllStores(ctx context.Context, getProductsForAllStoresRequest entities.GetProductsForAllStoresRequest) ([]*entities.GetProductsForAllStoresResponse, error) {
	warehouseID, err := primitive.ObjectIDFromHex(getProductsForAllStoresRequest.WarehouseID)
	if err != nil {
		return nil, err
	}
	allStores, err := u.productRepo.GetAllStores(ctx, warehouseID.Hex())
	if err != nil {
		return nil, err
	}
	products, err := u.productRepo.GetProductsForAllStores(ctx, allStores)
	if err != nil {
		return nil, err
	}
	return products, nil
}
