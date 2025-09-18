package db

import (
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/user/domain"
)

// Migrate runs migrations for the User domain model.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return errors.Wrap(err, "failed to migrate user model")
	}
	return nil
}
