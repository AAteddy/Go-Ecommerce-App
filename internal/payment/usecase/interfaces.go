package usecase

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *domain.Payment) error
	FindByID(ctx context.Context, id string) (*domain.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	UpdateStatus(ctx context.Context, id string, status string) (*domain.Payment, error)
}
