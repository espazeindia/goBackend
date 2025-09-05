package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"espazeBackend/utils"
	"fmt"

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
	collection := r.db.Collection("sellers")

	objectId, err := primitive.ObjectIDFromHex(sellerIdString)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "Seller Id is invalid",
		}, err
	}

	fmt.Print(sellerIdString, objectId)

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
	if request.Address != "" {
		docs["address"] = request.Address
	}
	if request.Pan != "" {
		docs["pan"] = request.Pan
	}
	docs["pin"] = request.PIN

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": docs})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to update user account",
		}, err
	}
	if result.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get seller  1",
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

func (r *OnboardingRepositoryMongoDB) GetBasicDetail(ctx context.Context, userIdString string) (*entities.MessageResponse, error) {
	collection := r.db.Collection("Seller")

	var sellerData *entities.Seller
	err := collection.FindOne(ctx, bson.M{"_id": userIdString}).Decode(&sellerData)
	if err == mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Error:   "No user found",
			Message: "Seller id does not exist in DB",
		}, err
	}

	sellerDetails := &entities.SellerBasicDetail{

		Name:        sellerData.Name,
		Address:     sellerData.Address,
		Gstin:       sellerData.Gstin,
		Pan:         sellerData.Pan,
		CompanyName: sellerData.CompanyName,
		PIN:         sellerData.PIN,
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Seller details fetched successfully",
		Data:    sellerDetails,
	}, nil

}
