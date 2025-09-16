package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/AAteddy/go-ecommerce-app/internal/inventory/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/inventory/usecase"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/dto"
	customErrors "github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
)

type InventoryHandler struct {
	uc  *usecase.InventoryUseCase
	log *logging.Logger
}

func NewInventoryHandler(uc *usecase.InventoryUseCase, log *logging.Logger) *InventoryHandler {
	return &InventoryHandler{uc, log}
}

func (h *InventoryHandler) RegisterRoutes(r chi.Router) {
	r.Post("/create_inventory", h.CreateInventory)
	r.Get("/get_inventory/{inventory_id}", h.GetInventoryByID)
	r.Get("/get_inventory/product/{product_id}", h.GetInventoryByProductID)
	// r.Get("/list_inventories", h.ListInventories)
	// r.Put("/update_inventory/{inventory_id}", h.UpdateInventory)
	// r.Delete("/delete_inventory/{inventory_id}", h.DeleteInventory)
}

func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateInventoryRequest
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
	var inventoryID string
	g.Go(func() error {
		var err error
		inventoryID, err = h.uc.Create(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Error creating inventory", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to create inventory",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusCreated,
		Message:    "Inventory created successfully",
		Error:      "",
		Data:       map[string]string{"inventory_id": inventoryID},
	}

	writeResponse(w, response)
}

func (h *InventoryHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	inventoryID := chi.URLParam(r, "inventory_id")
	if err := uuid.Validate(inventoryID); err != nil {
		h.log.Error("Invalid inventory ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid inventory ID",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var inventory *domain.Inventory
	g.Go(func() error {
		var err error
		inventory, err = h.uc.GetInventoryByID(ctx, inventoryID)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Error retrieving inventory by ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to retrieve inventory",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Inventory retrieved successfully",
		Error:      "",
		Data:       inventory,
	}

	writeResponse(w, response)
}

func (h *InventoryHandler) GetInventoryByProductID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
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
	var inventory *domain.Inventory
	g.Go(func() error {
		var err error
		inventory, err = h.uc.GetInventoryByProductID(ctx, productID)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Error retrieving inventory by product ID", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to retrieve inventory",
			Error:      err.Error(),
			Data:       nil,
		}
		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Inventory retrieved successfully",
		Error:      "",
		Data:       inventory,
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
