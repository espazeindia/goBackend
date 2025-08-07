package entities

import "time"

type Category struct {
	CategoryID        string    `json:"id" bson:"_id,omitempty"`
	CategoryName      string    `json:"category_name" bson:"category_name"`
	CategoryImage     string    `json:"category_image" bson:"category_image"`
	CategoryCreatedAt time.Time `json:"category_created_at" bson:"category_created_at"`
	CategoryUpdatedAt time.Time `json:"category_updated_at" bson:"category_updated_at"`
}

type Subcategory struct {
	SubcategoryID        string    `json:"id" bson:"_id,omitempty"`
	SubcategoryName      string    `json:"subcategory_name" bson:"subcategory_name"`
	SubcategoryImage     string    `json:"subcategory_image" bson:"subcategory_image"`
	CategoryID           string    `json:"category_id" bson:"category_id"`
	SubcategoryCreatedAt time.Time `json:"subcategory_created_at" bson:"subcategory_created_at"`
	SubcategoryUpdatedAt time.Time `json:"subcategory_updated_at" bson:"subcategory_updated_at"`
}

// Request DTOs
type CreateCategoryRequest struct {
	CategoryName  string `json:"category_name" binding:"required"`
	CategoryImage string `json:"category_image" binding:"required"`
}

type UpdateCategoryRequest struct {
	CategoryName  string `json:"category_name" binding:"required"`
	CategoryImage string `json:"category_image"`
}

type CreateSubcategoryRequest struct {
	SubcategoryName  string `json:"subcategory_name" binding:"required"`
	SubcategoryImage string `json:"subcategory_image" binding:"required"`
	CategoryID       string `json:"category_id" binding:"required"`
}

type UpdateSubcategoryRequest struct {
	SubcategoryName  string `json:"subcategory_name" binding:"required"`
	SubcategoryImage string `json:"subcategory_image"`
	CategoryID       string `json:"category_id" binding:"required"`
}

// Response DTOs
type CategoryWithSubcategoriesResponse struct {
	*Category
	Subcategories []*Subcategory `json:"subcategories"`
}

type CategoriesWithSubcategoriesResponse struct {
	Categories []*CategoryWithSubcategoriesResponse `json:"categories"`
}

type PaginatedCategoryResponse struct {
	Category   []*Category `json:"category"`
	Total      int64       `json:"total"`
	Limit      int64       `json:"limit"`
	Offset     int64       `json:"offset"`
	TotalPages int64       `json:"total_pages"`
}
type PaginatedSubCategoryResponse struct {
	SubCategory []*Subcategory `json:"sub_category"`
	Total       int64          `json:"total"`
	Limit       int64          `json:"limit"`
	Offset      int64          `json:"offset"`
	TotalPages  int64          `json:"total_pages"`
}
