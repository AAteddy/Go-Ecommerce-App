package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// PostgresInventoryRepository implements InventoryRepository for PostgreSQL.
type PostgresInventoryRepository struct {
	db *gorm.DB
}

// NewPostgresInventoryRepository creates a new instance of PostgresInventoryRepository.
func NewPostgresInventoryRepository(db *gorm.DB) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{db}
}

// implement InventoryRepository.Save inserts a new inventory record into the database.
func (r *PostgresInventoryRepository) Save(ctx context.Context, inventory *domain.Inventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

// implement InventoryRepository.FindByID retrieves an inventory record by its ID.
func (r *PostgresInventoryRepository) FindByID(ctx context.Context, id string) (*domain.Inventory, error) {
	var inventory domain.Inventory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to find inventory by the id")
	}

	return &inventory, nil
}

// implement InventoryRepository.FindByProductID retrieves an inventory record by its associated product ID.
func (r *PostgresInventoryRepository) FindByProductID(ctx context.Context, productID string) (*domain.Inventory, error) {
	var inventory domain.Inventory
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to find inventory by the Product ID")
	}

	return &inventory, nil
}
