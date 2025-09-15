package usecase

import (
	"context"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/product/domain"
)

type ProductUseCase struct {
	repo ProductRepository
	log  *logging.Logger
}

func NewProductUseCase(repo ProductRepository, log *logging.Logger) *ProductUseCase {
	return &ProductUseCase{repo, log}
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, req CreateProductRequest) (*domain.Product, error) {
	uc.log.Info("Creating new product ", "name ", req.Name)
	product, err := domain.NewProduct(req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		uc.log.Error("Failed to create product", "error", err)
		return nil, err
	}

	if err := uc.repo.Save(ctx, product); err != nil {
		uc.log.Error("Failed to save product", "error", err)
		return nil, errors.Wrap(err, "Failed to save product")
	}

	return product, nil
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	uc.log.Info("Retrieving product by ID ", "id ", id)
	product, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to find product by ID", "error", err)
		return nil, errors.Wrap(err, "Failed to find product by ID")
	}

	return product, nil
}

func (uc *ProductUseCase) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	uc.log.Info("Listing all products")
	products, err := uc.repo.List(ctx)
	if err != nil {
		uc.log.Error("Failed to list products", "error", err)
		return nil, errors.Wrap(err, "Failed to list products")
	}

	return products, nil
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// UpdateProduct updates an existing product by id.
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*domain.Product, error) {
	uc.log.Info("Updating product ", "id ", id)
	// first retrieve the existing product
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to find product by ID", "error", err)
		return nil, errors.Wrap(err, "Failed to find product by ID")
	}

	// update the existing product with the new values
	existing.Name = req.Name
	existing.Description = req.Description
	existing.Price = req.Price
	existing.Stock = req.Stock

	// save the updated product
	updatedProduct, err := uc.repo.Update(ctx, id, existing)
	if err != nil {
		uc.log.Error("Failed to update product", "error", err)
		return nil, errors.Wrap(err, "Failed to update product")
	}

	return updatedProduct, nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	uc.log.Info("Deleting product ", "id ", id)
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.log.Error("Failed to delete product", "error", err)
		return errors.Wrap(err, "Failed to delete product")
	}

	return nil
}
