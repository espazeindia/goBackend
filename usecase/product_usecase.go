package usecase

import (
	"context"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
)

// ProductUseCase handles business logic for product operations
type ProductUseCase struct {
	productRepo repositories.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repositories.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

// ProductResponse represents the response structure for product data
type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// PaginatedProductResponse represents paginated product response
type PaginatedProductResponse struct {
	Products    []ProductResponse `json:"products"`
	Total       int64             `json:"total"`
	Limit       int64             `json:"limit"`
	Offset      int64             `json:"offset"`
	HasNext     bool              `json:"has_next"`
	HasPrevious bool              `json:"has_previous"`
}

// GetAllProducts retrieves all products with pagination
func (uc *ProductUseCase) GetAllProducts(ctx context.Context, limit, offset int64) (*PaginatedProductResponse, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	products, total, err := uc.productRepo.GetAllProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = *uc.toProductResponse(product)
	}

	hasNext := offset+limit < total
	hasPrevious := offset > 0

	return &PaginatedProductResponse{
		Products:    productResponses,
		Total:       total,
		Limit:       limit,
		Offset:      offset,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}, nil
}

// toProductResponse converts a product entity to response format
func (uc *ProductUseCase) toProductResponse(product *entities.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.GetID(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
