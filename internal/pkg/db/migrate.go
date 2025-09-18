package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	inventoryDomain "github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
	orderDomain "github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	paymentDomain "github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
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
		&orderDomain.Order{},
		&productDomain.Product{},
		&paymentDomain.Payment{},
		&inventoryDomain.Inventory{},
	); err != nil {
		return nil, errors.Wrap(err, "failed to run database migrations")
	}

	// Add foreign key constraints
	// Check if constraints already exist
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.table_constraints WHERE constraint_name IN ('fk_user', 'fk_product', 'fk_order', 'fk_inventory_product') AND table_name IN ('orders', 'payments', 'inventory')").Scan(&count)
	if count == 0 {
		// Add foreign key constraints
		if err := db.Exec(`
			ALTER TABLE orders
			ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
			ADD CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT;
			ALTER TABLE payments
			ADD CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT;
			ALTER TABLE inventory
			ADD CONSTRAINT fk_inventory_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT;
		`).Error; err != nil {
			return nil, errors.Wrap(err, "failed to add foreign key constraints")
		}
	}

	return db, nil
}
