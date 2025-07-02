package entities

import "errors"

// Domain-specific errors
var (
	ErrProductNameRequired = errors.New("product name is required")
	ErrProductPriceInvalid = errors.New("product price must be greater than zero")
	ErrInvalidProductID    = errors.New("invalid product ID format")
	ErrProductNotFound     = errors.New("product not found")
)
