package db

import (
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// Migrate runs migrations for the Inventory domain model and adds foreign key constraints.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Inventory{}); err != nil {
		return errors.Wrap(err, "failed to migrate inventory model")
	}

	// Add foreign key constraints
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.table_constraints WHERE constraint_name = 'fk_inventory_product' AND table_name = 'inventory'").Scan(&count)
	if count == 0 {
		if err := db.Exec(`
			ALTER TABLE inventory
			ADD CONSTRAINT fk_inventory_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT;
		`).Error; err != nil {
			return errors.Wrap(err, "failed to add foreign key constraints for inventories")
		}
	}

	return nil
}
