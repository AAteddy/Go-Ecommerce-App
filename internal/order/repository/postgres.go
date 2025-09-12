package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// PostgresOrderRepository implements OrderRepository for PostgreSQL.
type PostgresOrderRepository struct {
	db *gorm.DB
}

// NewPostgresOrderRepository creates a new instance of PostgresOrderRepository.
func NewPostgresOrderRepository(db *gorm.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db}
}

// implement OrderRepository.Save inserts a new order into the database.
func (r *PostgresOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// implement OrderRepository.FindByID retrieves an order by its ID.
func (r *PostgresOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to find order by the id")
	}

	return &order, nil
}

// implement OrderRepository.ListByUserID retrieves all orders for a specific user.
func (r *PostgresOrderRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	var orders []*domain.Order
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "failed to list orders")
	}
	return orders, nil
}
