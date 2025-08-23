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

func (r *LocationRepositoryMongoDB) GetLocationForUserID(context context.Context, userId string) (*entities.MessageResponse, error) {
	var addresses []*entities.Location
	filter := bson.M{"user_id": userId}
	cursor, err := r.collection.Find(context, filter)
	if err == mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Message: "No Address for this user",
			Error:   "No document for this user id ",
		}, err
	} else if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "No Address for this user",
			Error:   err.Error(),
		}, err
	}
	err = cursor.All(context, &addresses)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Db Error",
			Error:   err.Error(),
		}, err
	}
	defer cursor.Close(context)
	return &entities.MessageResponse{
		Success: true,
		Data:    addresses,
	}, nil

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
