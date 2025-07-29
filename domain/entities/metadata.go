package entities

import (
	"time"
)

type Metadata struct {
	MetadataProductID     string    `json:"product_id" bson:"_id,omitempty"`
	MetadataHSNCode       string    `json:"hsn_code" bson:"hsn_code"`
	MetadataName          string    `json:"name" bson:"metadata_name"`
	MetadataDescription   string    `json:"description" bson:"metadata_description"`
	MetadataImage         string    `json:"image" bson:"metadata_image"`
	MetadataCategoryID    string    `json:"category_id" bson:"metadata_category_id"`
	MetadataSubcategoryID string    `json:"subcategory_id" bson:"metadata_subcategory_id"`
	MetadataMRP           float64   `json:"mrp" bson:"metadata_mrp"`
	MetadataCreatedAt     time.Time `json:"created_at" bson:"metadata_created_at"`
	MetadataUpdatedAt     time.Time `json:"updated_at" bson:"metadata_updated_at"`
}

type Review struct {
	MetadataProductID string `json:"metadata_product_id" bson:"_id"`
	TotalStars        int    `json:"total_stars" bson:"total_stars"`
	TotalReviews      int    `json:"total_reviews" bson:"total_reviews"`
}

type MetadataResponse struct {
	ID            string  `json:"id" bson:"_id,omitempty"`
	HsnCode       string  `json:"hsn_code"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Image         string  `json:"image"`
	CategoryID    string  `json:"category_id"`
	SubcategoryID string  `json:"subcategory_id"`
	MRP           float64 `json:"mrp"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	TotalStars    int     `json:"total_stars"`
	TotalReviews  int     `json:"total_reviews"`
}

type CreateMetadataResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Error    string    `json:"error" binding:"omitempty"`
	Id       string    `json:"id" binding:"omitempty"`
	Metadata *Metadata `json:"metadata" binding:"omitempty"`
}

// CreateMetadataRequest represents the request structure for creating metadata
type CreateMetadataRequest struct {
	Name          string  `json:"name" binding:"required"`
	HsnCode       string  `json:"hsn_code" binding:"required"`
	Description   string  `json:"description"`
	Image         string  `json:"image"`
	CategoryID    string  `json:"category_id"`
	SubcategoryID string  `json:"subcategory_id"`
	MRP           float64 `json:"mrp" binding:"required"`
}

// UpdateMetadataRequest represents the request structure for updating metadata
type UpdateMetadataRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Image         string  `json:"image"`
	CategoryID    string  `json:"category_id"`
	SubcategoryID string  `json:"subcategory_id"`
	MRP           float64 `json:"mrp"`
}

// PaginatedMetadataResponse represents paginated metadata response
type PaginatedMetadataResponse struct {
	Metadata    []*Metadata `json:"metadata"`
	Total       int64       `json:"total"`
	Limit       int64       `json:"limit"`
	Offset      int64       `json:"offset"`
	HasNext     bool        `json:"has_next"`
	HasPrevious bool        `json:"has_previous"`
	TotalPages  int64       `json:"total_pages"`
}

type AddReviewRequest struct {
	MetadataProductID string `json:"metadata_product_id" binding:"required"`
	Rating            int    `json:"rating" binding:"required"`
}
