package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *LocationRepositoryMongoDB) CreateLocation(ctx context.Context, locationRequest *entities.CreateLocationRequest) (*entities.MessageResponse, error) {
	locationData := &entities.Location{
		UserID:          locationRequest.UserID,
		LocationAddress: locationRequest.LocationAddress,
		Coordinates:     "0,0",
		Self:            locationRequest.Self,
		BuildingType:    locationRequest.BuildingType,
	}
	if locationRequest.Self {
		locationData = &entities.Location{
			UserID:          locationRequest.UserID,
			LocationAddress: locationRequest.LocationAddress,
			Coordinates:     "0,0",
			Self:            locationRequest.Self,
			Name:            locationRequest.Name,
			PhoneNumber:     locationRequest.PhoneNumber,
			BuildingType:    locationRequest.BuildingType,
		}
	}
	result, err := r.collection.InsertOne(context.Background(), locationData)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Db error",
			Error:   err.Error(),
		}, err
	}
	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Message: "Db error",
			Error:   "Error in getting inserted id",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Location Created SuccessFully",
	}, nil
}

func (r *LocationRepositoryMongoDB) GetLocationByAddress(address string) (*entities.Location, error) {
	var location entities.Location
	err := r.collection.FindOne(context.Background(), bson.M{"location_address": address}).Decode(&location)
	if err != nil {
		return nil, err
	}
	return &location, nil
}
