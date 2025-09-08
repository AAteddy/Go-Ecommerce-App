package usecase

import (
	"context"
	"time"

	"github.com/AAteddy/go-ecommerce-app/internal/user/domain"
)

// UserRepository defines the data access interfaces/methods for User entities.
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
}

// TokenStore defines the interface for token storage (Redis).
type TokenStore interface {
	SaveToken(ctx context.Context, token string, userID string, expiry time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	BlacklistToken(ctx context.Context, token string, expiry time.Duration) error
}

// EventProducer defines the interface for publishing events (Kafka).
type EventProducer interface {
	Produce(ctx context.Context, topic string, event []byte) error
}
