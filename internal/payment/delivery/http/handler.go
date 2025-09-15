package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/AAteddy/go-ecommerce-app/internal/payment/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/payment/usecase"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/dto"
	customErrors "github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
)

type PaymentHandler struct {
	uc  *usecase.PaymentUseCase
	log *logging.Logger
}

func NewPaymentHandler(uc *usecase.PaymentUseCase, log *logging.Logger) *PaymentHandler {
	return &PaymentHandler{uc, log}
}

func (h *PaymentHandler) RegisterRoutes(r chi.Router) {
	r.Post("/create_payment", h.CreatePayment)
	r.Get("/payments/{id}", h.GetPaymentByID)
	r.Get("/payments/order/{order_id}", h.GetPaymentByOrderID)
	r.Post("/payments/{id}/update_status", h.UpdatePaymentStatus)
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreatePaymentRequest
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
	var payment *domain.Payment
	g.Go(func() error {
		var err error
		payment, err = h.uc.CreatePayment(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to create payment", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to create payment",
			Error:      err.Error(),
			Data:       nil,
		}

		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Payment created successfully",
		Error:      "",
		Data:       payment,
	}

	writeResponse(w, response)
}

func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		h.log.Error("Payment ID is required")
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Payment ID is required",
			Error:      "",
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	if err := uuid.Validate(paymentID); err != nil {
		h.log.Error("Invalid payment ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid payment ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var payment *domain.Payment
	g.Go(func() error {
		var err error
		payment, err = h.uc.GetPaymentID(ctx, paymentID)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to get payment", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to get payment",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Payment retrieved successfully",
		Error:      "",
		Data:       payment,
	}

	writeResponse(w, response)
}

func (h *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")
	if orderID == "" {
		h.log.Error("Order ID is required")
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Order ID is required",
			Error:      "",
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	// if err := uuid.Validate(orderID); err != nil {
	// 	h.log.Error("Invalid order ID", "error", err)
	// 	response := dto.ApiResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Message:    "Invalid order ID",
	// 		Error:      err.Error(),
	// 		Data:       nil,
	// 	}
	// 	writeResponse(w, response)
	// 	return
	// }

	g, ctx := errgroup.WithContext(r.Context())
	var payment *domain.Payment
	g.Go(func() error {
		var err error
		payment, err = h.uc.GetPaymentByOrderID(ctx, orderID)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to get payment by order ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to get payment by order ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Payment retrieved successfully",
		Error:      "",
		Data:       payment,
	}

	writeResponse(w, response)
}

func (h *PaymentHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		h.log.Error("Payment ID is required")
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Payment ID is required",
			Error:      "",
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	if err := uuid.Validate(paymentID); err != nil {
		h.log.Error("Invalid payment ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid payment ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode request body", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to decode request body",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var payment *domain.Payment
	g.Go(func() error {
		var err error
		payment, err = h.uc.UpdatePaymentStatus(ctx, paymentID, req.Status)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Failed to update payment status", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to update payment status",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Payment status updated successfully",
		Error:      "",
		Data:       payment,
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
