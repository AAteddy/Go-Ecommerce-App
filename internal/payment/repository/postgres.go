package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// PostgresPaymentRepository implements PaymentRepository for PostgreSQL.
type PostgresPaymentRepository struct {
	db *gorm.DB
}

// NewPostgresPaymentRepository creates a new PostgresPaymentRepository.
func NewPostgresPaymentRepository(db *gorm.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

// implement PaymentRepository.Save inserts a new payment into the database.
func (r *PostgresPaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// implement PaymentRepository.FindByID retrieves a payment by its ID.
func (r *PostgresPaymentRepository) FindByID(ctx context.Context, id string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to find payment by the id")
	}

	return &payment, nil
}

// implement PaymentRepository.FindByOrderID retrieves a payment by its associated order ID.
func (r *PostgresPaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to find payment by the order id")
	}

	return &payment, nil
}
