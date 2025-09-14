package main

import (
	"context"
	"fmt"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/product/domain"
	"github.com/AAteddy/go-ecommerce-app/internal/product/repository"
)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing Product Service")

	dbConn, err := db.Migrate(cfg.DBURL) // migrates all tables
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	repo := repository.NewPostgresProductRepository(dbConn)
	product, err := domain.NewProduct("Laptop", "Gaming laptop", 999.99)
	if err != nil {
		log.Fatal("Failed to create product", "error", err)
	}

	if err := repo.Save(context.Background(), product); err != nil {
		log.Fatal("Failed to save product", "error", err)
	}

	log.Info("Product created", "id", product.ID)
	fmt.Printf("Product ID: %s\n", product.ID)

	// tokenStore := redis.NewRedisTokenStore(cfg.RedisAddr)
	// repo := repository.NewPostgresProductRepository(dbConn)
	// // usecase := usecase.NewUserUseCase(repo, tokenStore, nil, cfg.JWTSecret, log) // nil for Kafka (stubbed)

	// route := chi.NewRouter()
	// // handler := userhttp.NewUserHandler(usecase, log)
	// handler.RegisterRoutes(route)

	// log.Info("Starting user service on port: " + cfg.HTTPPort)
	// if err := http.ListenAndServe(cfg.HTTPPort, route); err != nil {
	// 	log.Fatal("Failed to start server", "error", err)
	// }
}
