package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func AuthMiddleware(usecase *usecase.UserUseCase, log *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				log.Error("Missing Authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
				tokenStr = tokenStr[7:]
			}

			if blacklisted, err := usecase.TokenStore.IsTokenBlacklisted(r.Context(), tokenStr); err != nil || blacklisted {
				log.Error("Token blacklisted or error ", "error ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// if usecase.TokenStore.IsTokenBlacklisted(r.Context(), tokenStr) {
			// 	log.Error("Token is blacklisted")
			// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// 	return
			// }

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(usecase.JwtSecret), nil
			}, jwt.WithLeeway(5*time.Second))
			if err != nil {
				log.Error("Invalid JWT", "error", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, ok := claims["sub"].(string)
				if !ok {
					log.Error("Invalid user_id in token")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				role, ok := claims["role"].(string)
				if !ok {
					log.Error("Invalid role in token")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				log.Info("Token claims", "claims", claims)
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				ctx = context.WithValue(ctx, RoleKey, role)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				log.Error("Invalid token claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}
