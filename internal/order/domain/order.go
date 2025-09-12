package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// Order represents an e-commerce order.
type Order struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"`
	ProductIDs []uuid.UUID `json:"product_ids" gorm:"type:uuid[]"`
	Total      float64     `json:"total" gorm:"type:decimal(10,2);not null;default:0.00"`
	Status     string      `json:"status" gorm:"size:50;not null;default:'pending'"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewOrder validates and creates a new order.
func NewOrder(userID uuid.UUID, productIDs []uuid.UUID, total float64) (*Order, error) {
	// Validate input
	if userID == uuid.Nil || len(productIDs) == 0 || total <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return &Order{
		ID:         uuid.New(),
		UserID:     userID,
		ProductIDs: productIDs,
		Total:      total,
		Status:     "pending",
	}, nil
}
