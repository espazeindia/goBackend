package mongodb

import (
	"context"
	"errors"
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

func (r *LocationRepositoryMongoDB) GetLocationForUserID(userId string) ([]*entities.Location, error) {
	var addresses []*entities.Location
	filter := bson.M{"user_id": userId}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var location entities.Location
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}
		addresses = append(addresses, &location)
	}

	if len(addresses) == 0 {
		return nil, errors.New("no location found")
	}

	return addresses, nil
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
