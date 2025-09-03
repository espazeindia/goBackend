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

func (r *OnboardingRepositoryMongoDB) AddBasicDetail(ctx context.Context, request *entities.SellerBasicDetail, sellerIdString string) (*entities.MessageResponse, error) {
	collection := r.db.Collection("seller")
	storeCollection := r.db.Collection("store")
	var existingUser entities.Seller
	err := collection.FindOne(ctx, bson.M{"_id": sellerIdString}).Decode(&existingUser)
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

	objectId, err := primitive.ObjectIDFromHex(sellerIdString)
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
	if result.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get seller ID",
		}, nil
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
	if results.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get seller ID",
		}, nil
	}
	token, err := utils.GenerateJWTToken(sellerIdString, request.Name, "seller", true)
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
