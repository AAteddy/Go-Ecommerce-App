package usecase

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/middleware"
)

type OrderUseCase struct {
	repo OrderRepository
	log  *logging.Logger
}

func NewOrderUseCase(repo OrderRepository, log *logging.Logger) *OrderUseCase {
	return &OrderUseCase{repo, log}
}

type CreateOrderRequest struct {
	ProductIDs []string `json:"product_ids"`
	Quantity   int      `json:"quantity"`
	Total      float64  `json:"total"`
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, req CreateOrderRequest) (string, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		uc.log.Error("Invalid or missing user_id in context")
		return "", errors.ErrInvalidInput
	}

	uc.log.Info("Creating new order for user ", "user_id ", userID)

	productIDs := req.ProductIDs

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
