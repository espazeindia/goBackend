package repositories

import (
	"context"

	"espazeBackend/domain/entities"
)

type MetadataRepository interface {
	GetAllMetadata(ctx context.Context, limit, offset int64, search string) ([]*entities.GetAllMetadata, int64, error)
	GetMetadataByID(ctx context.Context, id string) (*entities.MetadataResponse, error)
	CreateMetadata(ctx context.Context, metadata *entities.Metadata) (*entities.MetadataApiResponse, error)
	UpdateMetadata(ctx context.Context, id string, metadata *entities.Metadata) (*entities.MetadataApiResponse, error)
	DeleteMetadata(ctx context.Context, id string) (*entities.MetadataApiResponse, error)
	AddReview(ctx context.Context, req *entities.AddReviewRequest) error
	CreateReview(ctx context.Context, id string) (*entities.MetadataApiResponse, error)
}
