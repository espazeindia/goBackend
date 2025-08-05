package usecase

import (
	"context"
	"time"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
)

// CategorySubcategoryUseCase handles business logic for category and subcategory operations
type CategorySubcategoryUseCase struct {
	categorySubcategoryRepo repositories.CategorySubcategoryRepository
}

// NewCategorySubcategoryUseCase creates a new category subcategory use case
func NewCategorySubcategoryUseCase(categorySubcategoryRepo repositories.CategorySubcategoryRepository) *CategorySubcategoryUseCase {
	return &CategorySubcategoryUseCase{
		categorySubcategoryRepo: categorySubcategoryRepo,
	}
}

// Category operations
func (u *CategorySubcategoryUseCase) GetCategories(ctx context.Context, limit, offset int64, search *string) (*entities.PaginatedCategoryResponse, error) {
	if limit < 10 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	categories, total, err := u.categorySubcategoryRepo.GetCategories(ctx, limit, offset, search)
	if err != nil {
		return nil, err
	}
	var totalPages int64 = (total + limit - 1) / limit

	return &entities.PaginatedCategoryResponse{
		Category:   categories,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}, nil
}
func (u *CategorySubcategoryUseCase) GetAllCategories(ctx context.Context) ([]*entities.Category, error) {
	return u.categorySubcategoryRepo.GetAllCategories(ctx)
}

func (u *CategorySubcategoryUseCase) GetAllSubcategories(ctx context.Context, limit, offset int64, search *string) (*entities.PaginatedSubCategoryResponse, error) {
	if limit < 10 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	sub_category, total, err := u.categorySubcategoryRepo.GetAllSubcategories(ctx, limit, offset, search)
	if err != nil {
		return nil, err
	}
	var totalPages int64 = (total + limit - 1) / limit

	return &entities.PaginatedSubCategoryResponse{
		SubCategory: sub_category,
		Total:       total,
		Limit:       limit,
		Offset:      offset,
		TotalPages:  totalPages,
	}, nil
}

func (u *CategorySubcategoryUseCase) GetSubcategoryByCategoryId(ctx context.Context, CategoryID *string, limit, offset int64, search *string) (*entities.PaginatedSubCategoryResponse, error) {
	if limit < 10 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	sub_category, total, err := u.categorySubcategoryRepo.GetSubcategoryByCategoryId(ctx, CategoryID, limit, offset, search)
	if err != nil {
		return nil, err
	}
	var totalPages int64 = (total + limit - 1) / limit

	return &entities.PaginatedSubCategoryResponse{
		SubCategory: sub_category,
		Total:       total,
		Limit:       limit,
		Offset:      offset,
		TotalPages:  totalPages,
	}, nil

}

func (u *CategorySubcategoryUseCase) CreateCategory(ctx context.Context, request *entities.CreateCategoryRequest) (*entities.MessageResponse, error) {
	// Set timestamps
	now := time.Now()
	category := &entities.Category{
		CategoryName:      request.CategoryName,
		CategoryImage:     request.CategoryImage,
		CategoryCreatedAt: now,
		CategoryUpdatedAt: now,
	}

	return u.categorySubcategoryRepo.CreateCategory(ctx, category)
}

func (u *CategorySubcategoryUseCase) CreateSubcategory(ctx context.Context, request *entities.CreateSubcategoryRequest) (*entities.MessageResponse, error) {
	// Set timestamps
	now := time.Now()
	subcategory := &entities.Subcategory{
		SubcategoryName:      request.SubcategoryName,
		SubcategoryImage:     request.SubcategoryImage,
		CategoryID:           request.CategoryID,
		SubcategoryCreatedAt: now,
		SubcategoryUpdatedAt: now,
	}

	return u.categorySubcategoryRepo.CreateSubcategory(ctx, subcategory)
}

func (u *CategorySubcategoryUseCase) UpdateCategory(ctx context.Context, categoryId *string, request *entities.UpdateCategoryRequest) (*entities.MessageResponse, error) {
	return u.categorySubcategoryRepo.UpdateCategory(ctx, categoryId, request)
}

// func (u *CategorySubcategoryUseCase) GetCategoryById(ctx context.Context, categoryID string) (*entities.Category, error) {
// 	return u.categorySubcategoryRepo.GetCategoryById(ctx, categoryID)
// }

// func (u *CategorySubcategoryUseCase) DeleteCategory(ctx context.Context, categoryID string) error {
// 	return u.categorySubcategoryRepo.DeleteCategory(ctx, categoryID)
// }

// func (u *CategorySubcategoryUseCase) GetSubcategoriesByCategoryId(ctx context.Context, categoryID string) ([]*entities.Subcategory, error) {
// 	return u.categorySubcategoryRepo.GetSubcategoriesByCategoryId(ctx, categoryID)
// }

//

// func (u *CategorySubcategoryUseCase) UpdateSubcategory(ctx context.Context, subcategory *entities.Subcategory) error {
// 	// Update timestamp
// 	subcategory.SubcategoryUpdatedAt = time.Now()

// 	return u.categorySubcategoryRepo.UpdateSubcategory(ctx, subcategory)
// }

// func (u *CategorySubcategoryUseCase) DeleteSubcategory(ctx context.Context, subcategoryID string) error {
// 	return u.categorySubcategoryRepo.DeleteSubcategory(ctx, subcategoryID)
// }

// // Enhanced category response with subcategories
// type CategoryWithSubcategories struct {
// 	*entities.Category
// 	Subcategories []*entities.Subcategory `json:"subcategories"`
// }

// func (u *CategorySubcategoryUseCase) GetCategoryWithSubcategories(ctx context.Context, categoryID string) (*CategoryWithSubcategories, error) {
// 	// Get category
// 	category, err := u.categorySubcategoryRepo.GetCategoryById(ctx, categoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get subcategories for this category
// 	subcategories, err := u.categorySubcategoryRepo.GetSubcategoriesByCategoryId(ctx, categoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &CategoryWithSubcategories{
// 		Category:      category,
// 		Subcategories: subcategories,
// 	}, nil
// }

// func (u *CategorySubcategoryUseCase) GetAllCategoriesWithSubcategories(ctx context.Context) ([]*CategoryWithSubcategories, error) {
// 	// Get all categories
// 	categories, err := u.categorySubcategoryRepo.GetAllCategories(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result []*CategoryWithSubcategories

// 	// For each category, get its subcategories
// 	for _, category := range categories {
// 		subcategories, err := u.categorySubcategoryRepo.GetSubcategoriesByCategoryId(ctx, category.CategoryID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, &CategoryWithSubcategories{
// 			Category:      category,
// 			Subcategories: subcategories,
// 		})
// 	}

// 	return result, nil
// }
