package entities

import (
	"time"
)

type Metadata struct {
	MetadataProductID     string    `json:"metadata_product_id" bson:"_id,omitempty"`
	MetadataHSNCode       string    `json:"hsn_code" bson:"hsn_code"`
	MetadataName          string    `json:"metadata_name" bson:"metadata_name"`
	MetadataDescription   string    `json:"metadata_description" bson:"metadata_description"`
	MetadataImage         string    `json:"metadata_image" bson:"metadata_image"`
	MetadataCategoryID    string    `json:"metadata_category_id" bson:"metadata_category_id"`
	MetadataSubcategoryID string    `json:"metadata_subcategory_id" bson:"metadata_subcategory_id"`
	MetadataMRP           float64   `json:"metadata_mrp" bson:"metadata_mrp"`
	MetadataCreatedAt     time.Time `json:"metadata_created_at" bson:"metadata_created_at"`
	MetadataUpdatedAt     time.Time `json:"metadata_updated_at" bson:"metadata_updated_at"`
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
	Metadata    []*MetadataResponse `json:"metadata"`
	Total       int64               `json:"total"`
	Limit       int64               `json:"limit"`
	Offset      int64               `json:"offset"`
	HasNext     bool                `json:"has_next"`
	HasPrevious bool                `json:"has_previous"`
}

type AddReviewRequest struct {
	MetadataProductID string `json:"metadata_product_id" binding:"required"`
	Rating            int    `json:"rating" binding:"required"`
}
