package db

import (
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// Migrate runs migrations for the Order domain model and adds foreign key constraints.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Order{}); err != nil {
		return errors.Wrap(err, "failed to migrate order model")
	}

	// check if constraints already exist
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.table_constraints WHERE constraint_name IN ('fk_user', 'fk_product') AND table_name = 'orders'").Scan(&count)
	if count == 0 {
		if err := db.Exec(`
			ALTER TABLE orders
			ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;
		`).Error; err != nil {
			return errors.Wrap(err, "failed to add foreign key constraints for orders")
		}
	}

	return nil
}
