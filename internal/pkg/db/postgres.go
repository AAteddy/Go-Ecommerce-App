package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// NewPostgresDB establishes a connection to the PostgreSQL database without running migrations.
func NewPostgresDB(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to PostgreSQL database")
	}
	return db, nil
}
