package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/AAteddy/go-ecommerce-app/internal/order/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/order/usecase"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/dto"
	customErrors "github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/middleware"
	UserUC "github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

type OrderHandler struct {
	uc     *usecase.OrderUseCase
	userUC *UserUC.UserUseCase
	log    *logging.Logger
}

func NewOrderHandler(uc *usecase.OrderUseCase, userUC *UserUC.UserUseCase, log *logging.Logger) *OrderHandler {
	return &OrderHandler{uc, userUC, log}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.With(middleware.AuthMiddleware(h.userUC, h.log)).Post("/create_order", h.CreateOrder)
	r.With(middleware.AuthMiddleware(h.userUC, h.log)).Get("/get_order/{order_id}", h.GetOrderByID)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid request body", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request body",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var orderID string
	g.Go(func() error {
		var err error
		orderID, err = h.uc.CreateOrder(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to create order", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to create order",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Order created successfully",
		Error:      "",
		Data:       orderID,
	}

	writeResponse(w, response)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "order_id")
	orderID := vars
	if err := uuid.Validate(orderID); err != nil {
		h.log.Error("Invalid order ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var order *domain.Order
	g.Go(func() error {
		var err error
		order, err = h.uc.GetOrderByID(ctx, orderID)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to get order by ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to get order by ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Order retrieved successfully",
		Error:      "",
		Data:       order,
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
