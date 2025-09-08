package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	customErrors "github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
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
	r.With(AuthMiddleware(h.usecase)).Get("/profile", h.GetProfile)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req usecase.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var token string
	g.Go(func() error {
		var err error
		token, err = h.usecase.Register(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Registration failed", "error", err)
		http.Error(w, err.Error(), httpStatusFromError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req usecase.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	var token string
	g.Go(func() error {
		var err error
		token, err = h.usecase.Login(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error("Login failed", "error", err)
		http.Error(w, err.Error(), httpStatusFromError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		h.log.Error("Invalid or missing user_id in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := r.Context().Value(roleKey).(string)
	if !ok || role == "" {
		h.log.Error("Invalid or missing role in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user details from db
	user, err := h.usecase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to fetch user profile", "error", err)
		http.Error(w, err.Error(), httpStatusFromError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile retrieved successfully",
		"user_id": userID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    role,
	})

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
