package usecase

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/google/uuid"
)

type PaymentUseCase struct {
	repo PaymentRepository
	log  *logging.Logger
}

func NewPaymentUseCase(repo PaymentRepository, log *logging.Logger) *PaymentUseCase {
	return &PaymentUseCase{repo, log}
}

type CreatePaymentRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func (uc *PaymentUseCase) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*domain.Payment, error) {
	uc.log.Info("Creating new payment for order ", "order_id ", req.OrderID)
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		uc.log.Error("Invalid order ID", "error", err)
		return nil, errors.ErrInvalidInput
	}

	payment, err := domain.NewPayment(orderID, req.Amount)
	if err != nil {
		uc.log.Error("Failed to create payment", "error", err)
		return nil, err
	}

	if err := uc.repo.Save(ctx, payment); err != nil {
		uc.log.Error("Failed to save payment", "error", err)
		return nil, errors.Wrap(err, "Failed to save payment")
	}

	return payment, nil
}

func (uc *PaymentUseCase) GetPaymentID(ctx context.Context, id string) (*domain.Payment, error) {
	uc.log.Info("Retrieving payment by ID ", "id", id)
	payment, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to find payment by ID", "error", err)
		return nil, errors.Wrap(err, "Failed to find payment by ID")
	}

	return payment, nil
}

func (uc *PaymentUseCase) GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	uc.log.Info("Retrieving payment by order ID ", "order_id", orderID)
	payment, err := uc.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		uc.log.Error("Failed to find payment by order ID", "error", err)
		return nil, errors.Wrap(err, "Failed to find payment by order ID")
	}

	return payment, nil
}

func (uc *PaymentUseCase) UpdatePaymentStatus(ctx context.Context, id string, status string) (*domain.Payment, error) {
	uc.log.Info("Updating payment status ", "id", id, "status", status)
	payment, err := uc.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		uc.log.Error("Failed to update payment status", "error", err)
		return nil, errors.Wrap(err, "Failed to update payment status")
	}

	return payment, nil
}
