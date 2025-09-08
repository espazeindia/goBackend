package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type OnboardingRepository interface {
	AddBasicDetail(ctx context.Context, requestData *entities.SellerBasicDetail, sellerIdString string) (*entities.MessageResponse, error)
	GetBasicDetail(ctx context.Context, userIdString string) (*entities.MessageResponse, error)

	OnboardingAdmin(ctx context.Context, requestData *entities.AdminOnboaring) (*entities.MessageResponse, error)

	RegisterOperationalGuy(ctx context.Context, registrationRequest *entities.OperationalGuyRegistrationRequest) (*entities.OperationalGuyRegistrationResponse, error)
	GetOperationalGuy(ctx context.Context, userIdString string) (*entities.MessageResponse, error)
	EditOperationalGuy(ctx context.Context, request *entities.OperationalGuyGetRequest, operationsIdString string) (*entities.MessageResponse, error)
}
