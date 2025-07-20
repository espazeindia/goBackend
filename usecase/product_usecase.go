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

func (u *ProductUseCase) GetProductsForSpecificStore(ctx context.Context, getProductsForSpecificStoreRequest entities.GetProductsForSpecificStoreRequest) ([]*entities.GetProductsForSpecificStoreResponse, error) {
	storeID, err := primitive.ObjectIDFromHex(getProductsForSpecificStoreRequest.StoreID)
	if err != nil {
		return nil, err
	}
	sellerId, err := u.productRepo.FetchSellerId(ctx, storeID.Hex())
	if err != nil {
		return nil, err
	}

	products, err := u.productRepo.GetProductsForSpecificStore(ctx, sellerId)
	if err != nil {
		return nil, err
	}
	return products, nil
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
