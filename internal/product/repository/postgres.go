package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/product/domain"
)

// PostgresProductRepository implements ProductRepository for PostgreSQL.
type PostgresProductRepository struct {
	db *gorm.DB
}

// NewPostgresProductRepository creates a new PostgresProductRepository.
func NewPostgresProductRepository(db *gorm.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

// implement ProducctRepository.Save inserts a new product into the database.
func (r *PostgresProductRepository) Save(ctx context.Context, product *domain.Product) error {
	// Save product to the database
	return r.db.WithContext(ctx).Create(product).Error
}

// implement ProductRepository.FindByID retrieves a product by its ID.
func (r *PostgresProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to find product by id")
	}

	return &product, nil
}

// implement ProductRepository.List retrieves all products from the database.
func (r *PostgresProductRepository) List(ctx context.Context) ([]*domain.Product, error) {
	var products []*domain.Product

	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "failed to list products")
	}

	return products, nil
}

// implement ProductRepository.Update modifies an existing product in the database.
func (r *PostgresProductRepository) Update(ctx context.Context, id string, product *domain.Product) (*domain.Product, error) {
	// find the product by ID first to ensure it exists
	var existing domain.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to find product by id")
	}

	// update the product fields
	existing.Name = product.Name
	existing.Price = product.Price
	existing.Description = product.Description
	existing.Stock = product.Stock

	// save the updated product
	result := r.db.WithContext(ctx).Model(&existing).Updates(existing)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to update product")
	}

	return &existing, nil
}

// implement ProductRepository.Delete removes a product from the database by its ID.
func (r *PostgresProductRepository) Delete(ctx context.Context, id string) error {
	// find the product by ID first to ensure it exists
	var existing domain.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, "failed to find product by id")
	}

	// delete the product
	result := r.db.WithContext(ctx).Delete(&existing)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete product")
	}
	return nil
}
