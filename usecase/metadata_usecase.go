package usecase

import (
	"context"
	"time"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
)

// MetadataUseCase handles business logic for metadata operations
type MetadataUseCase struct {
	metadataRepo repositories.MetadataRepository
}

// NewMetadataUseCase creates a new metadata use case
func NewMetadataUseCase(metadataRepo repositories.MetadataRepository) *MetadataUseCase {
	return &MetadataUseCase{
		metadataRepo: metadataRepo,
	}
}

func (uc *MetadataUseCase) toMetadataResponse(metadata *entities.MetadataResponse) *entities.MetadataResponse {
	return &entities.MetadataResponse{
		ID:            metadata.ID,
		HsnCode:       metadata.HsnCode,
		Name:          metadata.Name,
		Description:   metadata.Description,
		Image:         metadata.Image,
		CategoryID:    metadata.CategoryID,
		SubcategoryID: metadata.SubcategoryID,
		MRP:           metadata.MRP,
		CreatedAt:     metadata.CreatedAt,
		UpdatedAt:     metadata.UpdatedAt,
		TotalStars:    metadata.TotalStars,
		TotalReviews:  metadata.TotalReviews,
	}
}

// GetAllMetadata retrieves all metadata with pagination
func (uc *MetadataUseCase) GetAllMetadata(ctx context.Context, limit, offset int64, search string) (*entities.PaginatedMetadataResponse, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	metadata, total, err := uc.metadataRepo.GetAllMetadata(ctx, limit, offset, search)
	if err != nil {
		return nil, err
	}

	hasNext := offset+limit < total
	hasPrevious := offset > 0
	var totalPages int64 = (total + limit - 1) / limit

	return &entities.PaginatedMetadataResponse{
		Metadata:    metadata,
		Total:       total,
		Limit:       limit,
		Offset:      offset,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
		TotalPages:  totalPages,
	}, nil
}

// GetMetadataByID retrieves a metadata by ID
func (uc *MetadataUseCase) GetMetadataByID(ctx context.Context, id string) (*entities.MetadataResponse, error) {
	metadata, err := uc.metadataRepo.GetMetadataByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.toMetadataResponse(metadata), nil
}

// CreateMetadata creates a new metadata
func (uc *MetadataUseCase) CreateMetadata(ctx context.Context, req *entities.CreateMetadataRequest) (*entities.CreateMetadataResponse, error) {
	// Generate a new product ID automatically (like UUID)

	now := time.Now()
	metadata := &entities.Metadata{
		MetadataName:          req.Name,
		MetadataHSNCode:       req.HsnCode,
		MetadataDescription:   req.Description,
		MetadataImage:         req.Image,
		MetadataCategoryID:    req.CategoryID,
		MetadataSubcategoryID: req.SubcategoryID,
		MetadataMRP:           req.MRP,
		MetadataCreatedAt:     now,
		MetadataUpdatedAt:     now,
	}

	response, err := uc.metadataRepo.CreateMetadata(ctx, metadata)
	if err != nil {
		return response, err
	}

	reviewResponse, err := uc.metadataRepo.CreateReview(ctx, response.Id)
	if err != nil {
		return reviewResponse, err
	}

	metadata.MetadataProductID = response.Id

	return response, nil
}

// UpdateMetadata updates an existing metadata
func (uc *MetadataUseCase) UpdateMetadata(ctx context.Context, id string, req *entities.UpdateMetadataRequest) (*entities.MetadataResponse, error) {

	now := time.Now()
	metadata := &entities.Metadata{
		MetadataName:          req.Name,
		MetadataDescription:   req.Description,
		MetadataImage:         req.Image,
		MetadataCategoryID:    req.CategoryID,
		MetadataSubcategoryID: req.SubcategoryID,
		MetadataMRP:           req.MRP,
		MetadataUpdatedAt:     now,
	}

	err := uc.metadataRepo.UpdateMetadata(ctx, id, metadata)
	if err != nil {
		return nil, err
	}

	// For now, return the updated metadata (in a real app, you'd fetch the updated record)
	return &entities.MetadataResponse{
		ID:            id,
		Name:          req.Name,
		Description:   req.Description,
		Image:         req.Image,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		MRP:           req.MRP,
		UpdatedAt:     now.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteMetadata deletes a metadata by ID
func (uc *MetadataUseCase) DeleteMetadata(ctx context.Context, id string) error {
	return uc.metadataRepo.DeleteMetadata(ctx, id)
}

func (uc *MetadataUseCase) AddReview(ctx context.Context, req *entities.AddReviewRequest) error {
	err := uc.metadataRepo.AddReview(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
