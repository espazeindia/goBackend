package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type LocationRepository interface {
	GetLocationForUserID(context context.Context, userId string) (*entities.MessageResponse, error)
	CreateLocation(location *entities.Location) error
	GetLocationByAddress(address string) (*entities.Location, error)
}
