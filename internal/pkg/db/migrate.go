package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"

	inventoryDomain "github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
	orderDomain "github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	paymentDomain "github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	productDomain "github.com/AAteddy/go-ecommerce-app/internal/product/domain"
	userDomain "github.com/AAteddy/go-ecommerce-app/internal/user/domain"
)

// Migrate initializes the database connection and runs migrations for all domain models.
func Migrate(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to PostgreSQL Database")
	}

	// Step 1: Create tables individually with error handling
	tables := []struct {
		name  string
		model interface{}
	}{
		{"users", &userDomain.User{}},
		{"products", &productDomain.Product{}},
		{"orders", &orderDomain.Order{}},
		{"payments", &paymentDomain.Payment{}},
		{"inventories", &inventoryDomain.Inventory{}},
	}

	log := logging.Init()
	log.Info("Starting database migrations")

	for _, table := range tables {
		log.Info("Creating table", "table", table.name)
		if err := db.AutoMigrate(table.model); err != nil {
			log.Error("Failed to migrate table", "table", table.name, "error", err)
			return nil, errors.Wrap(err, fmt.Sprintf("failed to migrate %s table", table.name))
		}
		log.Info("Successfully migrated table", "table", table.name)
	}

	// Step 2: Verify all tables exist
	tableNames := []string{"users", "products", "orders", "payments", "inventories"}
	for _, tableName := range tableNames {
		var count int64
		if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", tableName).Scan(&count).Error; err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to verify %s table existence", tableName))
		}
		if count == 0 {
			return nil, errors.New(fmt.Sprintf("%s table was not created", tableName))
		}
		log.Info("Verified table exists", "table", tableName)
	}

	// Step 3: Add foreign key constraints only if they don't exist
	constraints := []struct {
		table      string
		constraint string
		sql        string
	}{
		{"orders", "fk_user", `ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT`},
		{"payments", "fk_payment_order", `ALTER TABLE payments ADD CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT`},
		{"inventories", "fk_inventory_product", `ALTER TABLE inventories ADD CONSTRAINT fk_inventory_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT`},
	}

	for _, constraint := range constraints {
		var count int64
		if err := db.Raw(`SELECT COUNT(*) FROM information_schema.table_constraints 
			WHERE constraint_name = ? AND table_name = ?`, constraint.constraint, constraint.table).Scan(&count).Error; err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to check %s constraint existence", constraint.constraint))
		}
		if count == 0 {
			log.Info("Adding foreign key constraint", "constraint", constraint.constraint, "table", constraint.table)
			if err := db.Exec(constraint.sql).Error; err != nil {
				log.Error("Failed to add constraint", "constraint", constraint.constraint, "error", err)
				return nil, errors.Wrap(err, fmt.Sprintf("failed to add %s constraint", constraint.constraint))
			}
			log.Info("Successfully added constraint", "constraint", constraint.constraint)
		} else {
			log.Info("Constraint already exists", "constraint", constraint.constraint)
		}
	}

	log.Info("Database migrations completed successfully")
	return db, nil
}
