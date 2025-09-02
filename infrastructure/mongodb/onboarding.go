package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"espazeBackend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OnboardingRepositoryMongoDB struct {
	db *mongo.Database
}

func NewOnboardingRepositoryMongoDB(db *mongo.Database) repositories.OnboardingRepository {
	return &OnboardingRepositoryMongoDB{db: db}
}

func (r *OnboardingRepositoryMongoDB) AddBasicDetail(ctx context.Context, request *entities.SellerBasicDetail) (*entities.MessageResponse, error) {
	collection := r.db.Collection("seller")
	storeCollection := r.db.Collection("store")
	var existingUser entities.Seller
	err := collection.FindOne(ctx, bson.M{"id": request.SellerID}).Decode(&existingUser)
	if err == nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "User already exists",
			Message: "An account with this phone number already exists",
		}, nil
	} else if err != mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}

	objectId, err := primitive.ObjectIDFromHex(request.SellerID)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "Seller Id is invalid",
		}, err
	}

	docs := bson.M{}
	if request.Name != "" {
		docs["name"] = request.Name
	}
	if request.Gstin != "" {
		docs["gstin"] = request.Gstin
	}
	if request.CompanyName != "" {
		docs["companyName"] = request.CompanyName
	}
	if request.Pan != "" {
		docs["pan"] = request.Pan
	}
	docs["pin"] = request.PIN

	result, err := collection.UpdateByID(ctx, bson.M{"_id": objectId}, bson.M{"$set": docs})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to create user account",
		}, err
	}

	docStore := bson.M{}
	if request.ShopName != "" {
		docs["shopeName"] = request.ShopName
	}
	if request.ShopAddress != "" {
		docs["shopAddress"] = request.ShopAddress
	}
	results, err := storeCollection.UpdateByID(ctx, bson.M{"_id": objectId}, bson.M{"$set": docStore})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to create user account",
		}, err
	}

	_, ok := results.UpsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get seller ID",
		}, nil
	}

	_, ok = result.UpsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get seller ID",
		}, nil
	}

	token, err := utils.GenerateJWTToken(request.SellerID, request.Name, "seller", true)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Basic Details Saved Successfully",
		Token:   token,
	}, nil

}
