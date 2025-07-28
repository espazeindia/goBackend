package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type LoginRepository interface {
	LoginOperationalGuy(ctx context.Context, loginRequest *entities.OperationalGuyLoginRequest) (*entities.OperationalGuyLoginResponse, error)
	RegisterOperationalGuy(ctx context.Context, registrationRequest *entities.OperationalGuyRegistrationRequest) (*entities.OperationalGuyRegistrationResponse, error)
	VerifyOTP(ctx context.Context, phoneNumber *string, otp *int64) (*entities.SellerVerifyOTPResponse, error)
	RegisterSeller(ctx context.Context, registrationRequest *entities.SellerRegistrationRequest) (*entities.SellerRegistrationResponse, error)
	// LoginCustomer(ctx context.Context, loginRequest entities.CustomerLoginRequest) (entities.CustomerLoginResponse, error)
	// RegisterCustomer(ctx context.Context, registrationRequest entities.CustomerRegistrationRequest) (entities.CustomerRegistrationResponse, error)
}
