package repository

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
)

type InventoryRepository interface {
	Save(ctx context.Context, inventory *domain.Inventory) error
	FindByID(ctx context.Context, id string) (*domain.Inventory, error)
	FindByProductID(ctx context.Context, productID string) (*domain.Inventory, error)
	List(ctx context.Context) ([]*domain.Inventory, error)
	Update(ctx context.Context, id string, inventory *domain.Inventory) (*domain.Inventory, error)
	Delete(ctx context.Context, id string) error
}
