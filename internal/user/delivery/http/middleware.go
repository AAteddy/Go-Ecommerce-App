package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/dto"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func AuthMiddleware(usecase *usecase.UserUseCase) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logging.Init()
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				log.Warn("Missing or invalid Authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if blacklisted, err := usecase.TokenStore.IsTokenBlacklisted(r.Context(), tokenStr); err != nil || blacklisted {
				log.Error("Token blacklisted or error ", "error ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(usecase.JwtSecret), nil
			}, jwt.WithLeeway(5*time.Second))
			if err != nil || !token.Valid {
				log.Error("Invalid token ", "error ", err)
				response := dto.ApiResponse{
					StatusCode: http.StatusUnauthorized,
					Message:    "Unauthorized",
					Error:      "Invalid or missing token",
					Data:       nil,
				}

				writeResponse(w, response)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Error("Invalid claims ")
				response := dto.ApiResponse{
					StatusCode: http.StatusUnauthorized,
					Message:    "Unauthorized",
					Error:      "Invalid token claims",
					Data:       nil,
				}

				writeResponse(w, response)
				return
			}
			log.Info("Token claims", "claims", claims)

			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				log.Error("Invalid user_id claim ")
				response := dto.ApiResponse{
					StatusCode: http.StatusUnauthorized,
					Message:    "Unauthorized",
					Error:      "Invalid or missing user_id in context",
					Data:       nil,
				}

				writeResponse(w, response)
				return
			}

			role, ok := claims["role"].(string)
			if !ok || role == "" {
				log.Error("Invalid role claim ")
				response := dto.ApiResponse{
					StatusCode: http.StatusUnauthorized,
					Message:    "Unauthorized",
					Error:      "Invalid or missing role in context",
					Data:       nil,
				}

				writeResponse(w, response)
				return
			}

			// Passes authenticated data (from JWT) to downstream handlers without tight coupling.
			// WithValue creates a child context with key-value pairs (user_id, role from JWT claims).
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, roleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
