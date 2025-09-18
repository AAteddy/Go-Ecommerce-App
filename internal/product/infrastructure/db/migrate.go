package db

import (
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/product/domain"
)

// Migrate runs migrations for the Product domain model.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		return errors.Wrap(err, "failed to migrate product model")
	}

	return nil
}
