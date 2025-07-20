package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationRepositoryMongoDB struct {
	collection *mongo.Collection
}

func NewLocationRepositoryMongoDB(database *mongo.Database) repositories.LocationRepository {
	return &LocationRepositoryMongoDB{
		collection: database.Collection("locations"),
	}
}

func (r *LocationRepositoryMongoDB) CreateLocation(location *entities.Location) error {
	_, err := r.collection.InsertOne(context.Background(), location)
	return err
}

func (r *LocationRepositoryMongoDB) GetLocationByAddress(address string) (*entities.Location, error) {
	var location entities.Location
	err := r.collection.FindOne(context.Background(), bson.M{"location_address": address}).Decode(&location)
	if err != nil {
		return nil, err
	}
	return &location, nil
}
