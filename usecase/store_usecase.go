package usecase

import (
	"context"
	"errors"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"time"
)

type StoreUseCase struct {
	storeRepository repositories.StoreRepository
}

func NewStoreUseCase(storeRepository repositories.StoreRepository) *StoreUseCase {
	return &StoreUseCase{
		storeRepository: storeRepository,
	}
}

func (u *StoreUseCase) GetAllStores(ctx context.Context, request entities.GetAllStoresRequest) (*entities.GetAllStoresResponse, error) {
	// Validate request
	if request.WarehouseID == "" {
		return nil, errors.New("warehouse_id is required")
	}

	// Set default pagination values
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Offset < 0 {
		request.Offset = 0
	}

	// Get stores from repository
	paginatedResponse, err := u.storeRepository.GetAllStores(ctx, request)
	if err != nil {
		return nil, err
	}

	return &entities.GetAllStoresResponse{
		Success: true,
		Message: "Stores retrieved successfully",
		Stores:  paginatedResponse.Stores,
		Total:   paginatedResponse.Total,
		Limit:   paginatedResponse.Limit,
		Offset:  paginatedResponse.Offset,
	}, nil
}
func (u *StoreUseCase) GetAllStoresForCustomer(ctx context.Context, warehouseId, search string) ([]*entities.Store, error) {
	return u.storeRepository.GetAllStoresForCustomer(ctx, warehouseId, search)
}

func (u *StoreUseCase) GetStoreById(ctx context.Context, storeId string) (*entities.GetStoreResponse, error) {
	// Validate store ID
	if storeId == "" {
		return nil, errors.New("store_id is required")
	}

	// Get store from repository
	store, err := u.storeRepository.GetStoreById(ctx, storeId)
	if err != nil {
		return nil, err
	}

	return &entities.GetStoreResponse{
		Success: true,
		Message: "Store retrieved successfully",
		Store:   *store,
	}, nil
}

func (u *StoreUseCase) CreateStore(ctx context.Context, request *entities.CreateStoreRequest) (*entities.MessageResponse, error) {
	return u.storeRepository.CreateStore(ctx, request)
}

func (u *StoreUseCase) UpdateStore(ctx context.Context, storeId string, request entities.UpdateStoreRequest) (*entities.UpdateStoreResponse, error) {
	// Validate store ID
	if storeId == "" {
		return nil, errors.New("store_id is required")
	}

	// Get existing store
	existingStore, err := u.storeRepository.GetStoreById(ctx, storeId)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if request.StoreName != "" {
		existingStore.StoreName = request.StoreName
	}
	if request.StoreAddress != "" {
		existingStore.StoreAddress = request.StoreAddress
	}
	if request.StoreContact != "" {
		existingStore.StoreContact = request.StoreContact
	}
	if request.NumberOfRacks > 0 {
		existingStore.NumberOfRacks = request.NumberOfRacks
	}
	if request.OccupiedRacks >= 0 {
		// Validate that occupied racks don't exceed total racks
		if request.OccupiedRacks > existingStore.NumberOfRacks {
			return nil, errors.New("occupied_racks cannot exceed number_of_racks")
		}
		existingStore.OccupiedRacks = request.OccupiedRacks
	}

	// Update timestamp
	existingStore.UpdatedAt = time.Now()

	// Save to repository
	err = u.storeRepository.UpdateStore(ctx, storeId, existingStore)
	if err != nil {
		return nil, err
	}

	return &entities.UpdateStoreResponse{
		Success: true,
		Message: "Store updated successfully",
		Store:   *existingStore,
	}, nil
}

func (u *StoreUseCase) DeleteStore(ctx context.Context, storeId string) (*entities.DeleteStoreResponse, error) {
	// Validate store ID
	if storeId == "" {
		return nil, errors.New("store_id is required")
	}

	// Check if store exists
	_, err := u.storeRepository.GetStoreById(ctx, storeId)
	if err != nil {
		return nil, err
	}

	// Delete from repository
	err = u.storeRepository.DeleteStore(ctx, storeId)
	if err != nil {
		return nil, err
	}

	return &entities.DeleteStoreResponse{
		Success: true,
		Message: "Store deleted successfully",
	}, nil
}

func (u *StoreUseCase) GetStoreBySellerId(ctx context.Context, sellerId string) (*entities.GetStoreBySellerIdResponse, error) {
	// Validate seller ID
	if sellerId == "" {
		return nil, errors.New("seller_id is required")
	}

	// Get store from repository
	store, err := u.storeRepository.GetStoreBySellerId(ctx, sellerId)
	if err != nil {
		return nil, err
	}

	return &entities.GetStoreBySellerIdResponse{
		Success: true,
		Message: "Store retrieved successfully",
		Store:   *store,
	}, nil
}

func (u *StoreUseCase) UpdateStoreRacks(ctx context.Context, storeId string, occupiedRacks int) error {
	// Validate store ID
	if storeId == "" {
		return errors.New("store_id is required")
	}

	// Validate occupied racks
	if occupiedRacks < 0 {
		return errors.New("occupied_racks cannot be negative")
	}

	// Get existing store to validate against total racks
	existingStore, err := u.storeRepository.GetStoreById(ctx, storeId)
	if err != nil {
		return err
	}

	// Validate that occupied racks don't exceed total racks
	if occupiedRacks > existingStore.NumberOfRacks {
		return errors.New("occupied_racks cannot exceed number_of_racks")
	}

	// Update racks in repository
	return u.storeRepository.UpdateStoreRacks(ctx, storeId, occupiedRacks)
}
