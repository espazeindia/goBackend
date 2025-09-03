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
