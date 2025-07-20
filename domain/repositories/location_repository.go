package repositories

import "espazeBackend/domain/entities"

type LocationRepository interface {
	CreateLocation(location *entities.Location) error
	GetLocationByAddress(address string) (*entities.Location, error)
}
