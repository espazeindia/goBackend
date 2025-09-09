package usecase

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
)

type OnboardingUseCaseInterface struct {
	onboardingRepo repositories.OnboardingRepository
}

func NewOnboardingUseCase(onboardingRepo repositories.OnboardingRepository) *OnboardingUseCaseInterface {
	return &OnboardingUseCaseInterface{onboardingRepo: onboardingRepo}
}

func (u *OnboardingUseCaseInterface) AddBasicDetail(ctx context.Context, request *entities.SellerBasicDetail, sellerIdString string) (*entities.MessageResponse, error) {
	return u.onboardingRepo.AddBasicDetail(ctx, request, sellerIdString)
}

func (u *OnboardingUseCaseInterface) GetBasicDetail(ctx context.Context, userIdString string) (*entities.MessageResponse, error) {
	return u.onboardingRepo.GetBasicDetail(ctx, userIdString)

}

func (u *OnboardingUseCaseInterface) OnboardingAdmin(ctx context.Context, requestData *entities.AdminOnboaring) (*entities.MessageResponse, error) {
	return u.onboardingRepo.OnboardingAdmin(ctx, requestData)
}

func (u *OnboardingUseCaseInterface) RegisterOperationalGuy(ctx context.Context, registrationRequest *entities.OperationalGuyRegistrationRequest) (*entities.OperationalGuyRegistrationResponse, error) {
	return u.onboardingRepo.RegisterOperationalGuy(ctx, registrationRequest)
}

func (u *OnboardingUseCaseInterface) GetOperationalGuy(ctx context.Context, userIdString string) (*entities.MessageResponse, error) {
	return u.onboardingRepo.GetOperationalGuy(ctx, userIdString)
}

func (u *OnboardingUseCaseInterface) EditOperationalGuy(ctx context.Context, request *entities.OperationalGuyGetRequest, operationsIdString string) (*entities.MessageResponse, error) {
	return u.onboardingRepo.EditOperationalGuy(ctx, request, operationsIdString)
}

func (u *OnboardingUseCaseInterface) OnboardingOperationalGuy(ctx context.Context, request *entities.OperationalOnboarding, operationsIdString string) (*entities.MessageResponse, error) {
	return u.onboardingRepo.OnboardingOperationalGuy(ctx, request, operationsIdString)
}
