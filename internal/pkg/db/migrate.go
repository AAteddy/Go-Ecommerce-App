package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	productDomain "github.com/AAteddy/go-ecommerce-app/internal/product/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/user/domain"
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
		// &domain.Order{},
		&productDomain.Product{},
		// &domain.Payment{},
		// &domain.Inventory{},
	); err != nil {
		return nil, errors.Wrap(err, "failed to run database migrations")
	}

	return db, nil
}
