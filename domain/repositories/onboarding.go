package repositories

import (
	"context"
	"espazeBackend/domain/entities"
)

type OnboardingRepository interface {
	AddBasicDetail(ctx context.Context, requestData *entities.SellerBasicDetail) (*entities.MessageResponse, error)
}
