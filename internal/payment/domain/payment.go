package domain

import (
	"time"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/google/uuid"
)

// Payment represents the core business entity for a payment.
type Payment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	Amount    float64   `json:"amount" gorm:"not null;type:decimal(10,2)"`
	Status    string    `json:"status" gorm:"not null;size:50"` // e.g., "pending", "completed"
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewPayment creates a new payment with validation.
func NewPayment(orderID uuid.UUID, amount float64) (*Payment, error) {
	if orderID == uuid.Nil || amount <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return &Payment{
		ID:        uuid.New(),
		OrderID:   orderID,
		Amount:    amount,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
