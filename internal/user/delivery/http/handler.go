package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/dto"
	customErrors "github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/middleware"
	"github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
	"golang.org/x/sync/errgroup"
)

type UserHandler struct {
	usecase *usecase.UserUseCase
	log     *logging.Logger
}

func NewUserHandler(usecase *usecase.UserUseCase, log *logging.Logger) *UserHandler {
	return &UserHandler{usecase, log}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.With(middleware.AuthMiddleware(h.usecase, h.log)).Get("/profile", h.GetProfile)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req usecase.RegisterRequest
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

	// Using errgroup to manage concurrency and error handling.
	// This allows for future scalability if registration process involves multiple steps.
	// Context is passed to allow cancellation and timeouts.
	// In this simple case, it's a single goroutine, but structured for extensibility.
	g, ctx := errgroup.WithContext(r.Context())
	var token string
	g.Go(func() error {
		var err error
		token, err = h.usecase.Register(ctx, req)
		return err
	})

	// Wait for all goroutines to finish and check for errors.
	// If any goroutine returns an error, it will be captured here.
	// If no error, continue processing.
	if err := g.Wait(); err != nil {
		h.log.Error("Registration failed", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "User failed to register",
			Error:      err.Error(),
			Data:       nil,
		}

		writeResponse(w, response)
		return
	}

	tokenData := map[string]string{
		"token": token,
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully",
		Error:      "",
		Data:       tokenData,
	}

	writeResponse(w, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req usecase.LoginRequest
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

	// Using errgroup to manage concurrency and error handling.
	// This allows for future scalability if login process involves multiple steps.
	// Context is passed to allow cancellation and timeouts.
	// In this simple case, it's a single goroutine, but structured for extensibility.
	g, ctx := errgroup.WithContext(r.Context())
	var token string
	g.Go(func() error {
		var err error
		token, err = h.usecase.Login(ctx, req)
		return err
	})

	// Wait for all goroutines to finish and check for errors.
	// If any goroutine returns an error, it will be captured here.
	// If no error, continue processing.
	if err := g.Wait(); err != nil {
		h.log.Error("Login failed", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "User failed to login",
			Error:      err.Error(),
			Data:       nil,
		}

		writeResponse(w, response)
		return
	}

	tokenData := map[string]string{
		"token": token,
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "User logged in successfully",
		Error:      "",
		Data:       tokenData,
	}

	writeResponse(w, response)

}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// retrieves JWT data from middleware, enabling authorization without globals.
	// Extract user_id and role from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.log.Error("Invalid or missing user_id in context")
		response := dto.ApiResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			Error:      "Invalid or missing user_id in context",
			Data:       nil,
		}

		writeResponse(w, response)
		return
	}

	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role == "" {
		h.log.Error("Invalid or missing role in context")
		response := dto.ApiResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			Error:      "Invalid or missing role in context",
			Data:       nil,
		}

		writeResponse(w, response)
		return
	}

	// Fetch user details from db
	user, err := h.usecase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to fetch user profile", "error", err)
		response := dto.ApiResponse{
			StatusCode: httpStatusFromError(err),
			Message:    "Failed to fetch user profile",
			Error:      err.Error(),
			Data:       nil,
		}

		writeResponse(w, response)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    "Profile retrieved successfully",
		Error:      "",
		Data:       user,
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
