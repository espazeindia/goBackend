package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type LocationRepository interface {
	GetLocationForUserID(context context.Context, userId string) (*entities.MessageResponse, error)
	CreateLocation(ctx context.Context, location *entities.CreateLocationRequest) (*entities.MessageResponse, error)
	GetLocationByAddress(address string) (*entities.Location, error)
}
