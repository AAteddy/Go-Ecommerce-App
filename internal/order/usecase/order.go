package usecase

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/middleware"
	"github.com/AAteddy/go-ecommerce-app/internal/product/repository"
	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo        OrderRepository
	productRepo repository.ProductRepository
	log         *logging.Logger
}

func NewOrderUseCase(repo OrderRepository, productRepo repository.ProductRepository, log *logging.Logger) *OrderUseCase {
	return &OrderUseCase{repo, productRepo, log}
}

type CreateOrderRequest struct {
	ProductIDs []string `json:"product_ids"`
	Quantity   int      `json:"quantity"`
	Total      float64  `json:"total"`
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, req CreateOrderRequest) (string, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		uc.log.Error("Invalid or missing user_id in context")
		return "", errors.ErrInvalidInput
	}

	productIDs := req.ProductIDs
	// Validate ProductIDs
	for _, pid := range productIDs {
		_, err := uuid.Parse(pid)
		if err != nil {
			uc.log.Error("Invalid product_id format", "product_id", pid, "error", err)
			return "", errors.ErrInvalidInput
		}
		// Check if product exists
		exists, err := uc.productRepo.Exists(ctx, pid)
		if err != nil {
			uc.log.Error("Failed to check if product exists", "product_id", pid, "error", err)
			return "", errors.Wrap(err, "Failed to validate product")
		}
		if !exists {
			uc.log.Error("Product does not exist", "product_id", pid)
			return "", errors.ErrNotFound
		}
	}

	uc.log.Info("Creating order", "user_id", userID, "product_ids", productIDs)

	order, err := domain.NewOrder(userID, productIDs, req.Total, req.Quantity)
	if err != nil {
		uc.log.Error("Failed to create order", "error", err)
		return "", err
	}

	if err := uc.repo.Save(ctx, order); err != nil {
		uc.log.Error("Failed to save order", "error", err)
		return "", errors.Wrap(err, "Failed to save order")
	}

	return order.ID.String(), nil
}

func (uc *OrderUseCase) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	uc.log.Info("Retrieving order by ID ", "id ", id)
	order, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to find order by ID", "error", err)
		return nil, errors.Wrap(err, "Failed to find order by ID")
	}

	return order, nil
}
