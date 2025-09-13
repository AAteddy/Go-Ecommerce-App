package domain

import (
	"time"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/google/uuid"
)

// Inventory represents the core business entity for product stock.
type Inventory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewInventory creates a new inventory entry with validation.
func NewInventory(productID uuid.UUID, quantity int) (*Inventory, error) {
	if productID == uuid.Nil || quantity < 0 {
		return nil, errors.ErrInvalidInput
	}

	return &Inventory{
		ID:        uuid.New(),
		ProductID: productID,
		Quantity:  quantity,
	}, nil
}
