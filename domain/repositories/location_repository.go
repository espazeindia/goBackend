package repositories

import "espazeBackend/domain/entities"

type LocationRepository interface {
	GetLocationForUserID(userId string) ([]*entities.Location, error)
	CreateLocation(location *entities.Location) error
	GetLocationByAddress(address string) (*entities.Location, error)
}
