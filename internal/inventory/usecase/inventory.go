package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
)

type InventoryUseCase struct {
	repo InventoryRepository
	log  *logging.Logger
}

func NewInventoryUseCase(repo InventoryRepository, log *logging.Logger) *InventoryUseCase {
	return &InventoryUseCase{repo, log}
}

type CreateInventoryRequest struct {
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
}

func (uc *InventoryUseCase) CreateInventory(ctx context.Context, req CreateInventoryRequest) (string, error) {
	uc.log.Info("Creating new inventory for product ", "product_id ", req.ProductID)
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		uc.log.Error("Invalid product ID", "error", err)
		return "", errors.ErrInvalidInput
	}

	// var inventory *domain.Inventory
	inventory, err := domain.NewInventory(productID, req.Stock)
	if err != nil {
		uc.log.Error("Failed to create inventory", "error", err)
		return "", err
	}

	if err := uc.repo.Save(ctx, inventory); err != nil {
		uc.log.Error("Failed to save inventory", "error", err)
		return "", errors.Wrap(err, "Failed to save inventory")
	}

	return inventory.ID.String(), nil
}

func (uc *InventoryUseCase) GetInventoryByID(ctx context.Context, id string) (*domain.Inventory, error) {
	uc.log.Info("Retrieving inventory by ID ", "id ", id)
	inventory, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to find inventory by ID", "error", err)
		return nil, errors.Wrap(err, "Failed to find inventory by ID")
	}

	return inventory, nil
}

func (uc *InventoryUseCase) GetInventoryByProductID(ctx context.Context, productID string) (*domain.Inventory, error) {
	uc.log.Info("Retrieving inventory by Product ID ", "product_id ", productID)
	inventory, err := uc.repo.FindByProductID(ctx, productID)
	if err != nil {
		uc.log.Error("Failed to find inventory by Product ID", "error", err)
		return nil, errors.Wrap(err, "Failed to find inventory by Product ID")
	}

	return inventory, nil
}

// func (uc *InventoryUseCase) UpdateInventory(ctx context.Context, id string, req UpdateInventoryRequest) (*domain.Inventory, error) {
// 	uc.log.Info("Updating inventory by ID ", "id ", id)
// 	inventory, err := uc.repo.FindByID(ctx, id)
// 	if err != nil {
// 		uc.log.Error("Failed to find inventory by ID", "error", err)
// 		return nil, errors.Wrap(err, "Failed to find inventory by ID")
// 	}
// }
