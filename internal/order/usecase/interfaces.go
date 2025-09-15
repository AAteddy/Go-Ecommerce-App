package usecase

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/order/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	ListByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
}
