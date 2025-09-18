package repository

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/product/domain"
)

// ProductRepository defines data access for products.
type ProductRepository interface {
	Save(ctx context.Context, product *domain.Product) error
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Update(ctx context.Context, id string, product *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}
