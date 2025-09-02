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

	// objectId, err := primitive.ObjectIDFromHex(request.SellerID)
	// if err != nil {
	// 	return &entities.MessageResponse{
	// 		Success: false,
	// 		Error:   "Error in ObjectIdFromHex",
	// 		Message: "Seller Id is invalid",
	// 	}, err
	// }

	newUser := entities.SellerBasicDetail{
		Name:        request.Name,
		ShopAddress: request.ShopAddress,
		Gstin:       request.Gstin,
		Pan:         request.Pan,
		CompanyName: request.CompanyName,
		ShopName:    request.ShopName,
		PIN:         request.PIN,
	}

	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to create user account",
		}, err
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
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
