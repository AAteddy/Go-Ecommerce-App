package db

import (
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// Migrate runs migrations for the Payment domain model and adds foreign key constraints.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Payment{}); err != nil {
		return errors.Wrap(err, "failed to migrate payment model")
	}

	// Add foreign key constraints
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.table_constraints WHERE constraint_name = 'fk_order' AND table_name = 'payments'").Scan(&count)
	if count == 0 {
		if err := db.Exec(`
			ALTER TABLE payments
			ADD CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT;
		`).Error; err != nil {
			return errors.Wrap(err, "failed to add foreign key constraints for payments")
		}
	}

	return nil
}
