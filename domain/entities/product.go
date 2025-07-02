package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents the core product entity in the domain
type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewProduct creates a new Product instance with validation
func NewProduct(name, description string, price float64) (*Product, error) {
	if name == "" {
		return nil, ErrProductNameRequired
	}
	if price <= 0 {
		return nil, ErrProductPriceInvalid
	}

	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// Update updates the product with new values
func (p *Product) Update(name, description string, price float64) error {
	if name != "" {
		p.Name = name
	}
	if description != "" {
		p.Description = description
	}
	if price > 0 {
		p.Price = price
	}
	p.UpdatedAt = time.Now()
	return nil
}

// GetID returns the product ID as string
func (p *Product) GetID() string {
	return p.ID.Hex()
}

// SetID sets the product ID from string
func (p *Product) SetID(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidProductID
	}
	p.ID = objectID
	return nil
}
