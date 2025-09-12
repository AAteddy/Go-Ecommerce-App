package db

import (
	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
)

func NewPostgresDB(connString string) (*gorm.DB, error) {
	return db.Migrate(connString)
}
