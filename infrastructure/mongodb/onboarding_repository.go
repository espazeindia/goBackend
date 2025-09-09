package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"espazeBackend/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	docs["isFirstLogin"] = false
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
	collection := r.db.Collection("sellers")

	objectId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "objectId from Hex error",
			Message: "Error converting sellerId to ObjectId",
		}, err
	}

	var sellerData *entities.Seller
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&sellerData)
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
		Data:    sellerDetails,
	}, nil

}

func (r *OnboardingRepositoryMongoDB) OnboardingAdmin(ctx context.Context, requestData *entities.AdminOnboaring) (*entities.MessageResponse, error) {
	collection := r.db.Collection("admin")
	objectId, err := primitive.ObjectIDFromHex(requestData.AdminId)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "User Id is invalid",
		}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Password hashing failed",
			Message: "Failed to process registration",
		}, err
	}

	docs := bson.M{}
	docs["isFirstLogin"] = false
	docs["password"] = string(hashedPassword)

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": docs})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in db",
			Message: "Database Error",
		}, err
	}
	if result.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "no user ",
			Message: "No User Found",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Password Saved Successfully",
	}, nil
}

func (r *OnboardingRepositoryMongoDB) RegisterOperationalGuy(ctx context.Context, registrationRequest *entities.OperationalGuyRegistrationRequest) (*entities.OperationalGuyRegistrationResponse, error) {
	collection := r.db.Collection("operational_guys")

	var existingUser entities.OperationalGuy
	err := collection.FindOne(ctx, bson.M{"email": registrationRequest.Email}).Decode(&existingUser)
	if err == nil {
		return &entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "User already exists",
			Message: "An account with this email already exists",
		}, nil
	} else if err != mongo.ErrNoDocuments {
		return &entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return &entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Password hashing failed",
			Message: "Failed to process registration",
		}, err
	}

	now := time.Now()
	newUser := entities.OperationalGuy{
		Email:        registrationRequest.Email,
		Password:     string(hashedPassword),
		Name:         registrationRequest.Name,
		PhoneNumber:  registrationRequest.PhoneNumber,
		Address:      registrationRequest.Address,
		Pan:          registrationRequest.Pan,
		WarehouseId:  registrationRequest.WarehouseId,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsFirstLogin: true,
	}

	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return &entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to create user account",
		}, err
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get user ID",
		}, nil
	}

	return &entities.OperationalGuyRegistrationResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

func (r *OnboardingRepositoryMongoDB) GetOperationalGuy(ctx context.Context, userIdString string) (*entities.MessageResponse, error) {
	collection := r.db.Collection("operational_guys")
	collectionWarehouse := r.db.Collection("warehouses")

	objectId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "Operational Guy Id is invalid",
		}, err
	}
	var operationData *entities.OperationalGuy
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&operationData)
	if err == mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Error:   "No user found",
			Message: "Operational Guy id does not exist in DB",
		}, err
	}

	warehouseId, err := primitive.ObjectIDFromHex(operationData.WarehouseId)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "Warehouse Id is invalid",
		}, err
	}
	var warehouseData *entities.Warehouse
	err = collectionWarehouse.FindOne(ctx, bson.M{"_id": warehouseId}).Decode(&warehouseData)
	if err == mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Error:   "No user found",
			Message: "Warehouse Id does not exist in DB",
		}, err
	}

	operationDetails := &entities.OperationalGuyGetRespone{

		Name:          operationData.Name,
		Address:       operationData.Address,
		Email:         operationData.Email,
		Pan:           operationData.Pan,
		Password:      operationData.Password,
		PhoneNumber:   operationData.PhoneNumber,
		WarehouseId:   operationData.WarehouseId,
		WarehouseName: warehouseData.WarehouseName,
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Operational Guy fetched successfully",
		Data:    operationDetails,
	}, nil

}

func (r *OnboardingRepositoryMongoDB) EditOperationalGuy(ctx context.Context, request *entities.OperationalGuyGetRequest, operationsIdString string) (*entities.MessageResponse, error) {
	collection := r.db.Collection("operational_guys")

	objectId, err := primitive.ObjectIDFromHex(operationsIdString)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "Operational Guy Id is invalid",
		}, err
	}
	var operationalGuy *entities.OperationalGuy
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&operationalGuy)
	if err == mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Error:   "operational guy is not present in db",
			Message: "operational guy is not present in db",
		}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Password hashing failed",
			Message: "Failed to process registration",
		}, err
	}

	docs := bson.M{}
	if request.PhoneNumber != "" {
		docs["phoneNumber"] = request.PhoneNumber
	}
	if request.Address != "" {
		docs["address"] = request.Address
	}
	if request.Password != "" {
		docs["password"] = string(hashedPassword)
	}
	docs["isFirstLogin"] = false

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": docs})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to update changes",
		}, err
	}
	if result.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "No Operational Guy found",
		}, nil
	}

	token, err := utils.GenerateJWTToken(operationsIdString, operationalGuy.Name, "operations", true)
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

func (r *OnboardingRepositoryMongoDB) OnboardingOperationalGuy(ctx context.Context, request *entities.OperationalOnboarding, operationsIdString string) (*entities.MessageResponse, error) {
	collection := r.db.Collection("operational_guys")

	objectId, err := primitive.ObjectIDFromHex(operationsIdString)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "Operational Guy Id is invalid",
		}, err
	}
	var operationalGuy *entities.OperationalGuy
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&operationalGuy)
	if err == mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Error:   "operational guy is not present in db",
			Message: "operational guy is not present in db",
		}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Password hashing failed",
			Message: "Failed to process registration",
		}, err
	}

	docs := bson.M{}
	if request.Password != "" {
		docs["password"] = string(hashedPassword)
	}
	docs["isFirstLogin"] = false

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": docs})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to update changes",
		}, err
	}
	if result.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "No Operational Guy found",
		}, nil
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Password Saved Successfully",
	}, nil

}
