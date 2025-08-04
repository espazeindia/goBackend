package repositories

import (
	"context"

	"espazeBackend/domain/entities"
)

// CategorySubcategoryRepository defines the interface for category and subcategory data access operations
type CategorySubcategoryRepository interface {
	// Category operations
	GetAllCategories(ctx context.Context, limit, offset int64, search *string) ([]*entities.Category, int64, error)
	GetAllSubcategories(ctx context.Context, limit, offset int64, search *string) ([]*entities.Subcategory, int64, error)
	CreateCategory(ctx context.Context, category *entities.Category) (*entities.MessageResponse, error)

	// GetCategoryById(ctx context.Context, categoryID string) (*entities.Category, error)
	// UpdateCategory(ctx context.Context, category *entities.Category) error
	// DeleteCategory(ctx context.Context, categoryID string) error

	// // Subcategory operations
	// GetSubcategoryById(ctx context.Context, subcategoryID string) (*entities.Subcategory, error)
	// GetSubcategoriesByCategoryId(ctx context.Context, categoryID string) ([]*entities.Subcategory, error)
	// CreateSubcategory(ctx context.Context, subcategory *entities.Subcategory) error
	// UpdateSubcategory(ctx context.Context, subcategory *entities.Subcategory) error
	// DeleteSubcategory(ctx context.Context, subcategoryID string) error
}
