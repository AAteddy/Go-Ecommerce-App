package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/user/domain"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	// Save user to the database
	// Integrates context with GORM to make DB queries cancelable, timeout-aware, and transaction-safe. 
	// Context propagates from use cases to the DB layer for end-to-end request management.
	// Context is used to set a timeout or cancel the operation.
	// If the context is canceled or times out, the DB operation will be aborted.
	// This is particularly useful in web applications where requests can be canceled by the user.
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to find user by email")
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to find user by id")
	}
	return &user, nil
}
