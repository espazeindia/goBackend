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

func (l *LoginUseCaseInterface) LoginOperationalGuy(ctx context.Context, loginRequest entities.OperationalGuyLoginRequest) (entities.OperationalGuyLoginResponse, error) {
	return l.loginRepo.LoginOperationalGuy(ctx, loginRequest)
}

func (l *LoginUseCaseInterface) RegisterOperationalGuy(ctx context.Context, registrationRequest entities.OperationalGuyRegistrationRequest) (entities.OperationalGuyRegistrationResponse, error) {
	return l.loginRepo.RegisterOperationalGuy(ctx, registrationRequest)
}

func (l *LoginUseCaseInterface) LoginSeller(ctx context.Context, loginRequest entities.SellerLoginRequest) (entities.SellerLoginResponse, error) {
	return l.loginRepo.LoginSeller(ctx, loginRequest)
}

func (l *LoginUseCaseInterface) RegisterSeller(ctx context.Context, registrationRequest entities.SellerRegistrationRequest) (entities.SellerRegistrationResponse, error) {
	return l.loginRepo.RegisterSeller(ctx, registrationRequest)
}

func (l *LoginUseCaseInterface) LoginCustomer(ctx context.Context, loginRequest entities.CustomerLoginRequest) (entities.CustomerLoginResponse, error) {
	return l.loginRepo.LoginCustomer(ctx, loginRequest)
}

func (l *LoginUseCaseInterface) RegisterCustomer(ctx context.Context, registrationRequest entities.CustomerRegistrationRequest) (entities.CustomerRegistrationResponse, error) {
	return l.loginRepo.RegisterCustomer(ctx, registrationRequest)
}
