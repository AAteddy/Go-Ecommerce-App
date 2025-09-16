package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	userhttp "github.com/AAteddy/go-ecommerce-app/internal/user/delivery/http"

	"github.com/AAteddy/go-ecommerce-app/internal/user/infrastructure/redis"
	"github.com/AAteddy/go-ecommerce-app/internal/user/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing User Service")

	dbConn, err := db.Migrate(cfg.DBURL) // migrates all tables
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	tokenStore := redis.NewRedisTokenStore(cfg.RedisAddr)
	repo := repository.NewPostgresUserRepository(dbConn)
	usecase := usecase.NewUserUseCase(repo, tokenStore, nil, cfg.JWTSecret, log) // nil for Kafka (stubbed)

	route := chi.NewRouter()
	handler := userhttp.NewUserHandler(usecase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting user service on port: " + cfg.UserHTTPPort)
	if err := http.ListenAndServe(cfg.UserHTTPPort, route); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
