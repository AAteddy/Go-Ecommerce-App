package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	userhttp "github.com/AAteddy/go-ecommerce-app/internal/user/delivery/http"
	"github.com/AAteddy/go-ecommerce-app/internal/product/repository"

)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing User Service")

	dbConn, err := db.Migrate(cfg.DBURL) // migrates all tables
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	// tokenStore := redis.NewRedisTokenStore(cfg.RedisAddr)
	repo := repository.NewPostgresProductRepository(dbConn)
	// usecase := usecase.NewUserUseCase(repo, tokenStore, nil, cfg.JWTSecret, log) // nil for Kafka (stubbed)

	route := chi.NewRouter()
	// handler := userhttp.NewUserHandler(usecase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting user service on port: " + cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, route); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
