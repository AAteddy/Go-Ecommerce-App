package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/dto"
	customErrors "github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/middleware"
	"github.com/AAteddy/go-ecommerce-app/internal/product/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/product/usecase"
	UserUC "github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

type ProductHandler struct {
	uc     *usecase.ProductUseCase
	userUC *UserUC.UserUseCase
	log    *logging.Logger
}

func NewProductHandler(uc *usecase.ProductUseCase, userUC *UserUC.UserUseCase, log *logging.Logger) *ProductHandler {
	return &ProductHandler{uc, userUC, log}
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	r.With(middleware.AuthMiddleware(h.userUC, h.log)).Post("/create_product", h.CreateProduct)
	r.With(middleware.AuthMiddleware(h.userUC, h.log)).Get("/products/{id}", h.GetProductByID)
	r.With(middleware.AuthMiddleware(h.userUC, h.log)).Get("/products", h.ListProducts)
	r.With(middleware.AuthMiddleware(h.userUC, h.log)).Post("/update_product/{id}", h.UpdateProduct)
	r.With(middleware.AuthMiddleware(h.userUC, h.log)).Post("/delete_product/{id}", h.DeleteProduct)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid request body", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var product *domain.Product
	g.Go(func() error {
		var err error
		product, err = h.uc.CreateProduct(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to create product", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to create product",
			Error:      err.Error(),
			Data:       nil,
		}

		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusCreated,
		Message:    "Product created successfully",
		Error:      "",
		Data:       product,
	}

	writeResponse(w, response)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	productID := vars
	if err := uuid.Validate(productID); err != nil {
		h.log.Error("Invalid product ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid product ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var product *domain.Product
	g.Go(func() error {
		var err error
		product, err = h.uc.GetProductByID(ctx, productID)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to get product", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to get product",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Product retrieved successfully",
		Error:      "",
		Data:       product,
	}

	writeResponse(w, response)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	g, ctx := errgroup.WithContext(r.Context())
	var products []*domain.Product
	g.Go(func() error {
		var err error
		products, err = h.uc.ListProducts(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to list products", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to list products",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Products retrieved successfully",
		Error:      "",
		Data:       products,
	}

	writeResponse(w, response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	productID := vars
	if err := uuid.Validate(productID); err != nil {
		h.log.Error("Invalid product ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid product ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	var req usecase.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid request body", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var product *domain.Product
	g.Go(func() error {
		var err error
		product, err = h.uc.UpdateProduct(ctx, productID, req)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to update product", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to update product",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Product updated successfully",
		Error:      "",
		Data:       product,
	}

	writeResponse(w, response)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	productID := vars
	if err := uuid.Validate(productID); err != nil {
		h.log.Error("Invalid product ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid product ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	g.Go(func() error {
		return h.uc.DeleteProduct(ctx, productID)
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to delete product", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to delete product",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Product deleted successfully",
		Error:      "",
		Data:       nil,
	}

	writeResponse(w, response)
}

func httpStatusFromError(err error) int {
	if errors.Is(err, customErrors.ErrInvalidInput) {
		return http.StatusBadRequest
	}
	if errors.Is(err, customErrors.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func writeResponse(w http.ResponseWriter, response dto.ApiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
