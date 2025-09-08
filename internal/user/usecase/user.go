package usecase

import (
	"context"
	"time"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/user/domain"
	"github.com/golang-jwt/jwt/v5"
)

// UserUseCase defines the interface for user-related use cases.
type UserUseCase struct {
	repo          UserRepository
	TokenStore    TokenStore
	eventProducer EventProducer
	JwtSecret     string
	log           *logging.Logger
}

func NewUserUseCase(repo UserRepository, TokenStore TokenStore, eventProducer EventProducer, JwtSecret string, log *logging.Logger) *UserUseCase {
	return &UserUseCase{repo, TokenStore, eventProducer, JwtSecret, log}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

func (uc *UserUseCase) Register(ctx context.Context, req RegisterRequest) (string, error) {
	uc.log.Info("Registering new user ", "email ", req.Email)
	user, err := domain.NewUser(req.Email, req.Password, req.Role, req.Name)
	if err != nil {
		uc.log.Error("Failed to create user", "error", err)
		return "", err
	}

	if err := uc.repo.Save(ctx, user); err != nil {
		uc.log.Error("Failed to save user", "error", err)
		return "", errors.Wrap(err, "Failed to save user")
	}

	token, err := uc.generateJWT(user.ID.String(), user.Role)
	if err != nil {
		uc.log.Error("Failed to generate JWT", "error", err)
		return "", errors.Wrap(err, "failed to generate JWT")
	}

	if err := uc.TokenStore.SaveToken(ctx, token, user.ID.String(), time.Hour*24); err != nil {
		uc.log.Error("Failed to save token", "error", err)
		return "", errors.Wrap(err, "failed to save token")
	}

	// Stub for Kafka event (detailed in Part 6)
	// u.eventProducer.Produce(ctx, "user-registered", []byte(fmt.Sprintf(`{"user_id": "%s"}`, user.ID)))

	return token, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uc *UserUseCase) Login(ctx context.Context, req LoginRequest) (string, error) {
	uc.log.Info("User login attempt", "email", req.Email)
	user, err := uc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		uc.log.Error("Failed to find user by email", "error", err)
		return "", errors.ErrInvalidInput
	}

	if !user.VerifyPassword(req.Password) {
		uc.log.Warn("Invalid password attempt", "email", req.Email)
		return "", errors.ErrInvalidInput
	}

	token, err := uc.generateJWT(user.ID.String(), user.Role)
	if err != nil {
		uc.log.Error("Failed to generate JWT", "error", err)
		return "", errors.Wrap(err, "failed to generate JWT")
	}

	if err := uc.TokenStore.SaveToken(ctx, token, user.ID.String(), time.Hour*24); err != nil {
		uc.log.Error("Failed to save token", "error", err)
		return "", errors.Wrap(err, "failed to save token")
	}

	return token, nil
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	uc.log.Info("Fetching user by ID", "userID", userID)
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		uc.log.Error("Failed to find user by ID", "error", err)
		return nil, errors.Wrap(err, "failed to find user by ID")
	}
	return user, nil
}

func (uc *UserUseCase) generateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"iat":  time.Now().Unix(),
	}

	uc.log.Info("Generating JWT", "user_id", userID, "exp", claims["exp"])
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(uc.JwtSecret))
	if err != nil {
		uc.log.Error("Failed to sign JWT", "error", err)
		return "", err
	}

	return signedToken, nil
}
