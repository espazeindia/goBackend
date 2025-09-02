package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type LoginRepository interface {
	LoginOperationalGuy(ctx context.Context, loginRequest *entities.OperationalGuyLoginRequest) (*entities.OperationalGuyLoginResponse, error)
	RegisterOperationalGuy(ctx context.Context, registrationRequest *entities.OperationalGuyRegistrationRequest) (*entities.OperationalGuyRegistrationResponse, error)
	VerifyOTP(ctx context.Context, phoneNumber *string, otp *int64) (*entities.MessageResponse, error)
	VerifyPin(ctx context.Context, phoneNumber *string, pin *int64) (*entities.MessageResponse, error)
	GetOTP(ctx context.Context, phoneNumber string) (*entities.MessageResponse, error)
	VerifyOTPForCustomer(ctx context.Context, phoneNumber *string, otp *int64) (*entities.MessageResponse, error)
	VerifyPinForCustomer(ctx context.Context, phoneNumber *string, pin *int64) (*entities.MessageResponse, error)
	GetOTPForCustomer(ctx context.Context, phoneNumber string) (*entities.MessageResponse, error)
	LoginAdmin(ctx context.Context, loginRequest *entities.AdminLoginRequest) (*entities.AdminLoginResponse, error)
	RegisterAdmin(ctx context.Context, registrationRequest *entities.AdminRegistrationRequest) (*entities.AdminRegistrationResponse, error)
	CustomerBasicSetup(ctx context.Context, requestData *entities.CustomerBasicSetupRequest) (*entities.MessageResponse, error)
	// AddBasicData(ctx context.Context, request *entities.AddBasicData) (*entities.ResponseMessage, error)
	// LoginCustomer(ctx context.Context, loginRequest entities.CustomerLoginRequest) (entities.CustomerLoginResponse, error)
	// RegisterCustomer(ctx context.Context, registrationRequest entities.CustomerRegistrationRequest) (entities.CustomerRegistrationResponse, error)
}
