package repository

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *domain.Payment) error
	FindByID(ctx context.Context, id string) (*domain.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	// FindByUserID(ctx context.Context, userID string) ([]*domain.Payment, error)
}
