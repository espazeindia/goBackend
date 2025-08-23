package usecase

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
)

type LocationUseCase struct {
	locationRepository repositories.LocationRepository
}

func NewLocationUseCase(locationRepository repositories.LocationRepository) *LocationUseCase {
	return &LocationUseCase{
		locationRepository: locationRepository,
	}
}

func (uc *LocationUseCase) GetLocationForUserID(ctx context.Context, userId string) (*entities.MessageResponse, error) {
	return uc.locationRepository.GetLocationForUserID(ctx, userId)
}

func (uc *LocationUseCase) CreateLocation(location *entities.Location) error {
	return uc.locationRepository.CreateLocation(location)
}

func (uc *LocationUseCase) GetLocationByAddress(address string) (*entities.Location, error) {
	return uc.locationRepository.GetLocationByAddress(address)
}
