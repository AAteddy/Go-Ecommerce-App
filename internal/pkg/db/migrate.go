package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	productDomain "github.com/AAteddy/go-ecommerce-app/internal/product/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/user/domain"
	orderDomain "github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	paymentDomain "github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
	inventoryDomain "github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
)

// Migrate initializes the database connection and runs migrations for all domain models.
func Migrate(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to PostgreSQL Database")
	}

	// AutoMigrate creates/updates tables based on domain models
	if err := db.AutoMigrate(
		&domain.User{},
		&orderDomain.Order{},
		&productDomain.Product{},
		&paymentDomain.Payment{},
		&inventoryDomain.Inventory{},
	); err != nil {
		return nil, errors.Wrap(err, "failed to run database migrations")
	}

	return db, nil
}
