package repository

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
)

type InventoryRepository interface {
	Save(ctx context.Context, inventory *domain.Inventory) error
	FindByID(ctx context.Context, id string) (*domain.Inventory, error)
	FindByProductID(ctx context.Context, productID string) (*domain.Inventory, error)
}
