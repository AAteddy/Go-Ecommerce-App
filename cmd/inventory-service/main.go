package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	inventoryHttp "github.com/AAteddy/go-ecommerce-app/internal/inventory/delivery/http"
	inventoryDB "github.com/AAteddy/go-ecommerce-app/internal/inventory/infrastructure/db"
	"github.com/AAteddy/go-ecommerce-app/internal/inventory/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/inventory/usecase"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	"github.com/AAteddy/go-ecommerce-app/internal/user/infrastructure/redis"
	userRepo "github.com/AAteddy/go-ecommerce-app/internal/user/repository"
	UserUC "github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing Inventory Service")

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		dbConn, err := db.NewPostgresDB(cfg.DBURL)
		if err != nil {
			log.Fatal("Failed to connect to database", "error", err)
		}
		if err := inventoryDB.Migrate(dbConn); err != nil {
			log.Fatal("Failed to run migrations", "error", err)
		}
		log.Info("Inventory schema migrations completed")
	}

	dbConn, err := db.NewPostgresDB(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	tokenStore := redis.NewRedisTokenStore(cfg.RedisAddr)
	userRepo := userRepo.NewPostgresUserRepository(dbConn)
	userUseCase := UserUC.NewUserUseCase(userRepo, tokenStore, nil, cfg.JWTSecret, log)

	repo := repository.NewPostgresInventoryRepository(dbConn)
	usecase := usecase.NewInventoryUseCase(repo, log)

	route := chi.NewRouter()
	handler := inventoryHttp.NewInventoryHandler(usecase, userUseCase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting inventory service on port: " + cfg.InventoryHTTPPort)
	if err := http.ListenAndServe(cfg.InventoryHTTPPort, route); err != nil {
		log.Fatal("Failed to start inventory service", "error", err)
	}
}
