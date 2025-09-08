package usecase

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
)

type LoginUseCaseInterface struct {
	loginRepo repositories.LoginRepository
}

func NewLoginUseCase(loginRepo repositories.LoginRepository) *LoginUseCaseInterface {
	return &LoginUseCaseInterface{loginRepo: loginRepo}
}

func (l *LoginUseCaseInterface) LoginOperationalGuy(ctx context.Context, loginRequest *entities.OperationalGuyLoginRequest) (*entities.OperationalGuyLoginResponse, error) {
	return l.loginRepo.LoginOperationalGuy(ctx, loginRequest)
}

func (l *LoginUseCaseInterface) VerifyOTP(ctx context.Context, phoneNumber *string, otp *int64) (*entities.MessageResponse, error) {
	return l.loginRepo.VerifyOTP(ctx, phoneNumber, otp)
}
func (l *LoginUseCaseInterface) VerifyPin(ctx context.Context, phoneNumber *string, pin *int64) (*entities.MessageResponse, error) {
	return l.loginRepo.VerifyPin(ctx, phoneNumber, pin)
}

func (l *LoginUseCaseInterface) GetOTP(ctx context.Context, phoneNumber string) (*entities.MessageResponse, error) {
	return l.loginRepo.GetOTP(ctx, phoneNumber)
}
func (l *LoginUseCaseInterface) VerifyOTPForCustomer(ctx context.Context, phoneNumber *string, otp *int64) (*entities.MessageResponse, error) {
	return l.loginRepo.VerifyOTPForCustomer(ctx, phoneNumber, otp)
}
func (l *LoginUseCaseInterface) VerifyPinForCustomer(ctx context.Context, phoneNumber *string, pin *int64) (*entities.MessageResponse, error) {
	return l.loginRepo.VerifyPinForCustomer(ctx, phoneNumber, pin)
}

func (l *LoginUseCaseInterface) GetOTPForCustomer(ctx context.Context, phoneNumber string) (*entities.MessageResponse, error) {
	return l.loginRepo.GetOTPForCustomer(ctx, phoneNumber)
}

func (l *LoginUseCaseInterface) LoginAdmin(ctx context.Context, loginRequest *entities.AdminLoginRequest) (*entities.AdminLoginResponse, error) {
	return l.loginRepo.LoginAdmin(ctx, loginRequest)
}

func (l *LoginUseCaseInterface) RegisterAdmin(ctx context.Context, registrationRequest *entities.AdminRegistrationRequest) (*entities.AdminRegistrationResponse, error) {
	return l.loginRepo.RegisterAdmin(ctx, registrationRequest)
}

func (l *LoginUseCaseInterface) CustomerBasicSetup(ctx context.Context, requestData *entities.CustomerBasicSetupRequest) (*entities.MessageResponse, error) {
	return l.loginRepo.CustomerBasicSetup(ctx, requestData)
}

// func (l *LoginUseCaseInterface) AddBasicData(ctx context.Context, request *entities.AddBasicData) (*entities.ResponseMessage, error) {
// 	return l.loginRepo.AddBasicData(ctx, request)
// }

// func (l *LoginUseCaseInterface) LoginCustomer(ctx context.Context, loginRequest entities.CustomerLoginRequest) (entities.CustomerLoginResponse, error) {
// 	return l.loginRepo.LoginCustomer(ctx, loginRequest)
// }

// func (l *LoginUseCaseInterface) RegisterCustomer(ctx context.Context, registrationRequest entities.CustomerRegistrationRequest) (entities.CustomerRegistrationResponse, error) {
// 	return l.loginRepo.RegisterCustomer(ctx, registrationRequest)
// }
