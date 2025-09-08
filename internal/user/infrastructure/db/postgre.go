package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/user/domain"
)

func NewPostgresDB(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to PostgreSQL Database")
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return nil, errors.Wrap(err, "failed to migrate User model in database")
	}
	return db, nil
}
