package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// Product represents a product in the e-commerce catalog.
type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewProduct validates and creates a new product.
func NewProduct(name, description string, price float64) (*Product, error) {
	// Validate input
	if name == "" || price <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
	}, nil

}
