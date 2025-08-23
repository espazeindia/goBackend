package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"espazeBackend/utils"
	"fmt"
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

func (r *LoginRepositoryMongoDB) LoginOperationalGuy(ctx context.Context, loginRequest *entities.OperationalGuyLoginRequest) (*entities.OperationalGuyLoginResponse, error) {
	collection := r.db.Collection("operational_guys")

	// Find user by email
	filter := bson.M{"email": loginRequest.Email}
	var operationalGuy entities.OperationalGuy
	err := collection.FindOne(ctx, filter).Decode(&operationalGuy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &entities.OperationalGuyLoginResponse{
				Success: false,
				Error:   "Invalid credentials",
				Message: "Invalid Email Credentials",
			}, nil
		}
		return &entities.OperationalGuyLoginResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to authenticate user",
		}, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(operationalGuy.Password), []byte(loginRequest.Password))
	if err != nil {
		return &entities.OperationalGuyLoginResponse{
			Success: false,
			Error:   "Invalid credentials",
			Message: "Invalid Password Credentials",
		}, nil
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(operationalGuy.OperationalGuyID, operationalGuy.Name, "operations")
	if err != nil {
		return &entities.OperationalGuyLoginResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	// Update last login time and set isFirstLogin to false
	now := time.Now()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"operationalGuyID": operationalGuy.OperationalGuyID},
		bson.M{"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		}},
	)
	if err != nil {
		fmt.Println("Error updating last login time:", err)
	}

	return &entities.OperationalGuyLoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil
}

func (r *LoginRepositoryMongoDB) RegisterOperationalGuy(ctx context.Context, registrationRequest *entities.OperationalGuyRegistrationRequest) (*entities.OperationalGuyRegistrationResponse, error) {
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

func (r *LoginRepositoryMongoDB) VerifyOTP(ctx context.Context, phoneNumber *string, otp *int64) (*entities.MessageResponse, error) {
	sellerCollection := r.db.Collection("sellers")

	var existingUser entities.Seller
	err := sellerCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		// User already exists
		return &entities.MessageResponse{
			Success: false,
			Error:   "No Seller Found",
			Message: "No Seller is associated to this phone number ",
		}, nil
	} else if err != mongo.ErrNoDocuments && err != nil {
		// Database error
		return &entities.MessageResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}
	now := time.Now()
	fiveMinutes := 5 * time.Minute

	if int(*otp) != existingUser.OTP {

		return &entities.MessageResponse{
			Success: false,
			Error:   "WRONG OTP",
			Message: "OTP is incorrect",
		}, err
	}

	if now.After(existingUser.OTPGeneratedAt.Add(fiveMinutes)) {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP Expired",
			Message: "OTP has expired try using the RESEND OTP",
		}, err
	}

	token, err := utils.GenerateJWTToken(existingUser.SellerID, existingUser.Name, "seller")
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	_, err = sellerCollection.UpdateOne(
		ctx,
		bson.M{"_id": existingUser.SellerID},
		bson.M{"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		}},
	)
	if err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil

}

func (r *LoginRepositoryMongoDB) VerifyPin(ctx context.Context, phoneNumber *string, pin *int64) (*entities.MessageResponse, error) {
	sellerCollection := r.db.Collection("sellers")

	var existingUser entities.Seller
	err := sellerCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		// User already exists
		return &entities.MessageResponse{
			Success: false,
			Error:   "No Seller Found",
			Message: "No Seller is associated to this phone number ",
		}, nil
	} else if err != mongo.ErrNoDocuments && err != nil {
		// Database error
		return &entities.MessageResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}
	now := time.Now()

	if int(*pin) != existingUser.PIN {

		return &entities.MessageResponse{
			Success: false,
			Error:   "WRONG Pin",
			Message: "Pin is incorrect",
		}, err
	}

	token, err := utils.GenerateJWTToken(existingUser.SellerID, existingUser.Name, "seller")
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	_, err = sellerCollection.UpdateOne(
		ctx,
		bson.M{"_id": existingUser.SellerID},
		bson.M{"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		}},
	)
	if err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil

}

func (r *LoginRepositoryMongoDB) GetOTP(ctx context.Context, phoneNumber string) (*entities.MessageResponse, error) {
	sellerCollection := r.db.Collection("sellers")

	var existingUser entities.Seller
	err := sellerCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		otp, err := utils.GenerateOTP()
		if err != nil {
			return &entities.MessageResponse{
				Success: false,
				Error:   "OTP generation failed",
				Message: "Failed to generate OTP",
			}, err
		}

		now := time.Now()
		newUser := entities.Seller{
			Name:               "new user",
			PhoneNumber:        phoneNumber,
			Address:            "dummy",
			OTP:                otp,
			OTPGeneratedAt:     now,
			NumberOfRetriesOTP: 0,
			PIN:                -1,
			NumberOfRetriesPIN: 0,
			LastLoginAt:        now,
			StoreID:            "",
		}
		// Insert user into database
		response, err := sellerCollection.InsertOne(ctx, newUser)
		if err != nil {
			return &entities.MessageResponse{
				Success: false,
				Error:   "Registration failed",
				Message: "Failed to create seller account",
			}, err
		}
		_, ok := response.InsertedID.(primitive.ObjectID)
		if !ok {
			return &entities.MessageResponse{
				Success: false,
				Error:   "Registration failed",
				Message: "Failed to create seller account",
			}, err
		}

		return &entities.MessageResponse{
			Success: true,
			Message: fmt.Sprint("Otp Sent Successfully ", otp),
		}, nil
	} else if err != mongo.ErrNoDocuments && err != nil {
		// Database error
		return &entities.MessageResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP generation failed",
			Message: "Failed to generate OTP",
		}, err
	}
	now := time.Now()

	docs := bson.M{}

	docs["otp"] = otp
	docs["otpGeneratedAt"] = now
	filter := bson.M{"phoneNumber": existingUser.PhoneNumber}
	update := bson.M{"$set": docs}

	response, err := sellerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP storage failed",
			Message: "Failed to update OTP",
		}, err
	}

	if response.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP storage failed",
			Message: "Failed to update OTP",
		}, err
	}

	return &entities.MessageResponse{
		Success: true,
		Message: fmt.Sprint("Otp Sent Successfully ", otp),
	}, nil

}
func (r *LoginRepositoryMongoDB) VerifyOTPForCustomer(ctx context.Context, phoneNumber *string, otp *int64) (*entities.MessageResponse, error) {
	customerCollection := r.db.Collection("customers")

	var existingUser entities.Customer
	err := customerCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		// User already exists
		return &entities.MessageResponse{
			Success: false,
			Error:   "No User Found",
			Message: "No User is associated to this phone number ",
		}, nil
	} else if err != mongo.ErrNoDocuments && err != nil {
		// Database error
		return &entities.MessageResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}
	now := time.Now()
	fiveMinutes := 5 * time.Minute

	if int(*otp) != existingUser.OTP {

		return &entities.MessageResponse{
			Success: false,
			Error:   "WRONG OTP",
			Message: "OTP is incorrect",
		}, err
	}

	if now.After(existingUser.OTPGeneratedAt.Add(fiveMinutes)) {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP Expired",
			Message: "OTP has expired try using the RESEND OTP",
		}, err
	}

	token, err := utils.GenerateJWTToken(existingUser.CustomerID, existingUser.Name, "customer")
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	_, err = customerCollection.UpdateOne(
		ctx,
		bson.M{"_id": existingUser.CustomerID},
		bson.M{"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		}},
	)
	if err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}
	if existingUser.Name == "new user" {
		return &entities.MessageResponse{
			Success:     true,
			Message:     "Login successful",
			Token:       token,
			IsOnboarded: false,
		}, nil
	}
	return &entities.MessageResponse{
		Success:     true,
		Message:     "Login successful",
		Token:       token,
		IsOnboarded: true,
	}, nil
}

func (r *LoginRepositoryMongoDB) VerifyPinForCustomer(ctx context.Context, phoneNumber *string, pin *int64) (*entities.MessageResponse, error) {
	customerCollection := r.db.Collection("customers")

	var existingUser entities.Customer
	err := customerCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		// User already exists
		return &entities.MessageResponse{
			Success: false,
			Error:   "No User Found",
			Message: "No User is associated to this phone number ",
		}, nil
	} else if err != mongo.ErrNoDocuments && err != nil {
		// Database error
		return &entities.MessageResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}
	now := time.Now()

	if int(*pin) != existingUser.PIN {

		return &entities.MessageResponse{
			Success: false,
			Error:   "WRONG Pin",
			Message: "Pin is incorrect",
		}, err
	}

	token, err := utils.GenerateJWTToken(existingUser.CustomerID, existingUser.Name, "customer")
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}

	_, err = customerCollection.UpdateOne(
		ctx,
		bson.M{"_id": existingUser.CustomerID},
		bson.M{"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		}},
	)
	if err != nil {
		// Log the error but don't fail the login
		// You might want to add proper logging here
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil

}

func (r *LoginRepositoryMongoDB) GetOTPForCustomer(ctx context.Context, phoneNumber string) (*entities.MessageResponse, error) {
	customerCollection := r.db.Collection("customers")

	var existingUser entities.Customer
	err := customerCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		otp, err := utils.GenerateOTP()
		if err != nil {
			return &entities.MessageResponse{
				Success: false,
				Error:   "OTP generation failed",
				Message: "Failed to generate OTP",
			}, err
		}

		now := time.Now()
		newUser := entities.Customer{
			Name:               "new user",
			PhoneNumber:        phoneNumber,
			OTP:                otp,
			OTPGeneratedAt:     now,
			NumberOfRetriesOTP: 0,
			PIN:                -1,
			NumberOfRetriesPIN: 0,
			LastLoginAt:        now,
		}
		// Insert user into database
		response, err := customerCollection.InsertOne(ctx, newUser)
		if err != nil {
			return &entities.MessageResponse{
				Success: false,
				Error:   "Registration failed",
				Message: "Failed to create user account",
			}, err
		}
		_, ok := response.InsertedID.(primitive.ObjectID)
		if !ok {
			return &entities.MessageResponse{
				Success: false,
				Error:   "Registration failed",
				Message: "Failed to create user account",
			}, err
		}

		return &entities.MessageResponse{
			Success: true,
			Message: fmt.Sprint("Otp Sent Successfully ", otp),
		}, nil
	} else if err != mongo.ErrNoDocuments && err != nil {
		// Database error
		return &entities.MessageResponse{
			Success: false,
			Error:   "Database error",
			Message: "Failed to check user existence",
		}, err
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP generation failed",
			Message: "Failed to generate OTP",
		}, err
	}
	now := time.Now()

	docs := bson.M{}

	docs["otp"] = otp
	docs["otpGeneratedAt"] = now
	filter := bson.M{"phoneNumber": existingUser.PhoneNumber}
	update := bson.M{"$set": docs}

	response, err := customerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP storage failed",
			Message: "Failed to update OTP",
		}, err
	}

	if response.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Error:   "OTP storage failed",
			Message: "Failed to update OTP",
		}, err
	}

	return &entities.MessageResponse{
		Success: true,
		Message: fmt.Sprint("Otp Sent Successfully ", otp),
	}, nil

}

func (r *LoginRepositoryMongoDB) CustomerBasicSetup(ctx context.Context, requestData *entities.CustomerBasicSetupRequest) (*entities.MessageResponse, error) {
	customerCollection := r.db.Collection("customers")
	locationCollection := r.db.Collection("locations")
	objectId, err := primitive.ObjectIDFromHex(requestData.UserId)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in ObjectIdFromHex",
			Message: "User Id is invalid",
		}, err
	}
	docs := bson.M{}
	if requestData.Name != "" {
		docs["name"] = requestData.Name
	}
	docs["pin"] = requestData.PIN

	result, err := customerCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": docs})
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
	locationData := &entities.Location{
		UserID:          requestData.UserId,
		LocationAddress: requestData.Address,
		Self:            true,
		BuildingType:    "home",
		Coordinates:     "0,0",
	}
	insertedData, err := locationCollection.InsertOne(ctx, locationData)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in db",
			Message: "Database Error",
		}, err
	}
	_, ok := insertedData.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Error in db",
			Message: "Database Error",
		}, err
	}
	token, err := utils.GenerateJWTToken(requestData.UserId, requestData.Name, "customer")
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "User Registed Successfully",
		Token:   token,
	}, nil

}

// func (r *LoginRepositoryMongoDB) RegisterCustomer(ctx context.Context, registrationRequest entities.CustomerRegistrationRequest) (entities.CustomerRegistrationResponse, error) {
// 	collection := r.db.Collection("customers")

// 	// Check if user already exists
// 	var existingUser entities.Customer
// 	err := collection.FindOne(ctx, bson.M{"email": registrationRequest.Email}).Decode(&existingUser)
// 	if err == nil {
// 		// User already exists
// 		return entities.CustomerRegistrationResponse{
// 			Success: false,
// 			Error:   "User already exists",
// 			Message: "An account with this email already exists",
// 		}, nil
// 	} else if err != mongo.ErrNoDocuments {
// 		// Database error
// 		return entities.CustomerRegistrationResponse{
// 			Success: false,
// 			Error:   "Database error",
// 			Message: "Failed to check user existence",
// 		}, err
// 	}

// 	// Hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return entities.CustomerRegistrationResponse{
// 			Success: false,
// 			Error:   "Password hashing failed",
// 			Message: "Failed to process registration",
// 		}, err
// 	}

// 	now := time.Now()
// 	newUser := entities.Customer{
// 		Email:        registrationRequest.Email,
// 		Password:     string(hashedPassword),
// 		Name:         registrationRequest.Name,
// 		IsFirstLogin: true,
// 		PhoneNumber:  registrationRequest.PhoneNumber,
// 		Address:      registrationRequest.Address,
// 		DateOfBirth:  registrationRequest.DateOfBirth,
// 		IsVerified:   false,
// 		CreatedAt:    now,
// 		UpdatedAt:    now,
// 	}

// 	// Insert user into database
// 	result, err := collection.InsertOne(ctx, newUser)
// 	if err != nil {
// 		return entities.CustomerRegistrationResponse{
// 			Success: false,
// 			Error:   "Registration failed",
// 			Message: "Failed to create user account",
// 		}, err
// 	}

// 	// Get the inserted user ID
// 	objectID, ok := result.InsertedID.(primitive.ObjectID)
// 	if !ok {
// 		return entities.CustomerRegistrationResponse{
// 			Success: false,
// 			Error:   "Registration failed",
// 			Message: "Failed to get user ID",
// 		}, nil
// 	}

// 	return entities.CustomerRegistrationResponse{
// 		Success: true,
// 		Message: "User registered successfully",
// 		UserID:  objectID.Hex(),
// 	}, nil
// }
