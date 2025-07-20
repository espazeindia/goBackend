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

type LoginRepositoryMongoDB struct {
	db *mongo.Database
}

func NewLoginRepositoryMongoDB(db *mongo.Database) repositories.LoginRepository {
	return &LoginRepositoryMongoDB{db: db}
}

func (r *LoginRepositoryMongoDB) LoginOperationalGuy(ctx context.Context, loginRequest entities.OperationalGuyLoginRequest) (entities.OperationalGuyLoginResponse, error) {
	collection := r.db.Collection("operational_guys")

	// Find user by email
	var operationalGuy entities.OperationalGuy
	err := collection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&operationalGuy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entities.OperationalGuyLoginResponse{
				Success: false,
				Error:   "Invalid credentials",
				Message: "Email or password is incorrect",
			}, nil
		}
		return entities.OperationalGuyLoginResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to authenticate user",
		}, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(operationalGuy.Password), []byte(loginRequest.Password))
	if err != nil {
		return entities.OperationalGuyLoginResponse{
			Success: false,
			Error:   "Invalid credentials",
			Message: "Email or password is incorrect",
		}, nil
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(operationalGuy.ID, operationalGuy.Email)
	if err != nil {
		return entities.OperationalGuyLoginResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	// Update last login time and set isFirstLogin to false
	now := time.Now()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": operationalGuy.ID},
		bson.M{"$set": bson.M{
			"lastLoginAt":  now,
			"isFirstLogin": false,
			"updatedAt":    now,
		}},
	)
	if err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}

	return entities.OperationalGuyLoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil
}

func (r *LoginRepositoryMongoDB) RegisterOperationalGuy(ctx context.Context, registrationRequest entities.OperationalGuyRegistrationRequest) (entities.OperationalGuyRegistrationResponse, error) {
	collection := r.db.Collection("operational_guys")

	// Check if user already exists
	var existingUser entities.OperationalGuy
	err := collection.FindOne(ctx, bson.M{"email": registrationRequest.Email}).Decode(&existingUser)
	if err == nil {
		// User already exists
		return entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "User already exists",
			Message: "An account with this email already exists",
		}, nil
	} else if err != mongo.ErrNoDocuments {
		// Database error
		return entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Password hashing failed",
			Message: "Failed to process registration",
		}, err
	}

	now := time.Now()
	newUser := entities.OperationalGuy{
		Email:                  registrationRequest.Email,
		Password:               string(hashedPassword),
		Name:                   registrationRequest.Name,
		IsFirstLogin:           true,
		PhoneNumber:            registrationRequest.PhoneNumber,
		Address:                registrationRequest.Address,
		EmergencyContactNumber: registrationRequest.EmergencyContactNumber,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	// Insert user into database
	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to create user account",
		}, err
	}

	// Get the inserted user ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return entities.OperationalGuyRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get user ID",
		}, nil
	}

	return entities.OperationalGuyRegistrationResponse{
		Success: true,
		Message: "User registered successfully",
		UserID:  objectID.Hex(),
	}, nil
}

func (r *LoginRepositoryMongoDB) LoginSeller(ctx context.Context, loginRequest entities.SellerLoginRequest) (entities.SellerLoginResponse, error) {
	collection := r.db.Collection("sellers")

	// Find user by email
	var seller entities.Seller
	err := collection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&seller)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entities.SellerLoginResponse{
				Success: false,
				Error:   "Invalid credentials",
				Message: "Email or password is incorrect",
			}, nil
		}
		return entities.SellerLoginResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to authenticate user",
		}, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(seller.Password), []byte(loginRequest.Password))
	if err != nil {
		return entities.SellerLoginResponse{
			Success: false,
			Error:   "Invalid credentials",
			Message: "Email or password is incorrect",
		}, nil
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(seller.ID, seller.Email)
	if err != nil {
		return entities.SellerLoginResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	// Update last login time and set isFirstLogin to false
	now := time.Now()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": seller.ID},
		bson.M{"$set": bson.M{
			"lastLoginAt":  now,
			"isFirstLogin": false,
			"updatedAt":    now,
		}},
	)
	if err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}

	return entities.SellerLoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil
}

func (r *LoginRepositoryMongoDB) RegisterSeller(ctx context.Context, registrationRequest entities.SellerRegistrationRequest) (entities.SellerRegistrationResponse, error) {
	collection := r.db.Collection("sellers")

	// Check if user already exists
	var existingUser entities.Seller
	err := collection.FindOne(ctx, bson.M{"email": registrationRequest.Email}).Decode(&existingUser)
	if err == nil {
		// User already exists
		return entities.SellerRegistrationResponse{
			Success: false,
			Error:   "User already exists",
			Message: "An account with this email already exists",
		}, nil
	} else if err != mongo.ErrNoDocuments {
		// Database error
		return entities.SellerRegistrationResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.SellerRegistrationResponse{
			Success: false,
			Error:   "Password hashing failed",
			Message: "Failed to process registration",
		}, err
	}

	now := time.Now()
	newUser := entities.Seller{
		Email:        registrationRequest.Email,
		Password:     string(hashedPassword),
		Name:         registrationRequest.Name,
		IsFirstLogin: true,
		PhoneNumber:  registrationRequest.PhoneNumber,
		Address:      registrationRequest.Address,
		BusinessName: registrationRequest.BusinessName,
		BusinessType: registrationRequest.BusinessType,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Insert user into database
	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return entities.SellerRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to create user account",
		}, err
	}

	// Get the inserted user ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return entities.SellerRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get user ID",
		}, nil
	}

	return entities.SellerRegistrationResponse{
		Success: true,
		Message: "User registered successfully",
		UserID:  objectID.Hex(),
	}, nil
}

func (r *LoginRepositoryMongoDB) LoginCustomer(ctx context.Context, loginRequest entities.CustomerLoginRequest) (entities.CustomerLoginResponse, error) {
	collection := r.db.Collection("customers")

	// Find user by email
	var customer entities.Customer
	err := collection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entities.CustomerLoginResponse{
				Success: false,
				Error:   "Invalid credentials",
				Message: "Email or password is incorrect",
			}, nil
		}
		return entities.CustomerLoginResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to authenticate user",
		}, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginRequest.Password))
	if err != nil {
		return entities.CustomerLoginResponse{
			Success: false,
			Error:   "Invalid credentials",
			Message: "Email or password is incorrect",
		}, nil
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(customer.ID, customer.Email)
	if err != nil {
		return entities.CustomerLoginResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	// Update last login time and set isFirstLogin to false
	now := time.Now()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": customer.ID},
		bson.M{"$set": bson.M{
			"lastLoginAt":  now,
			"isFirstLogin": false,
			"updatedAt":    now,
		}},
	)
	if err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}

	return entities.CustomerLoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil
}

func (r *LoginRepositoryMongoDB) RegisterCustomer(ctx context.Context, registrationRequest entities.CustomerRegistrationRequest) (entities.CustomerRegistrationResponse, error) {
	collection := r.db.Collection("customers")

	// Check if user already exists
	var existingUser entities.Customer
	err := collection.FindOne(ctx, bson.M{"email": registrationRequest.Email}).Decode(&existingUser)
	if err == nil {
		// User already exists
		return entities.CustomerRegistrationResponse{
			Success: false,
			Error:   "User already exists",
			Message: "An account with this email already exists",
		}, nil
	} else if err != mongo.ErrNoDocuments {
		// Database error
		return entities.CustomerRegistrationResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.CustomerRegistrationResponse{
			Success: false,
			Error:   "Password hashing failed",
			Message: "Failed to process registration",
		}, err
	}

	now := time.Now()
	newUser := entities.Customer{
		Email:        registrationRequest.Email,
		Password:     string(hashedPassword),
		Name:         registrationRequest.Name,
		IsFirstLogin: true,
		PhoneNumber:  registrationRequest.PhoneNumber,
		Address:      registrationRequest.Address,
		DateOfBirth:  registrationRequest.DateOfBirth,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Insert user into database
	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return entities.CustomerRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to create user account",
		}, err
	}

	// Get the inserted user ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return entities.CustomerRegistrationResponse{
			Success: false,
			Error:   "Registration failed",
			Message: "Failed to get user ID",
		}, nil
	}

	return entities.CustomerRegistrationResponse{
		Success: true,
		Message: "User registered successfully",
		UserID:  objectID.Hex(),
	}, nil
}
