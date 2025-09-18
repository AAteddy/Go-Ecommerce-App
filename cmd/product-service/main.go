package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	productHttp "github.com/AAteddy/go-ecommerce-app/internal/product/delivery/http"
	productDB "github.com/AAteddy/go-ecommerce-app/internal/product/infrastructure/db"
	"github.com/AAteddy/go-ecommerce-app/internal/product/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/product/usecase"
	"github.com/AAteddy/go-ecommerce-app/internal/user/infrastructure/redis"
	userRepo "github.com/AAteddy/go-ecommerce-app/internal/user/repository"
	UserUC "github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing Product Service")

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		dbConn, err := db.NewPostgresDB(cfg.DBURL)
		if err != nil {
			log.Fatal("Failed to connect to database", "error", err)
		}
		if err := productDB.Migrate(dbConn); err != nil {
			log.Fatal("Failed to run migrations", "error", err)
		}
		log.Info("Product schema migrations completed")
	}

	dbConn, err := db.NewPostgresDB(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	tokenStore := redis.NewRedisTokenStore(cfg.RedisAddr)
	userRepo := userRepo.NewPostgresUserRepository(dbConn)
	userUseCase := UserUC.NewUserUseCase(userRepo, tokenStore, nil, cfg.JWTSecret, log)

	repo := repository.NewPostgresProductRepository(dbConn)
	usecase := usecase.NewProductUseCase(repo, log)

	route := chi.NewRouter()
	handler := productHttp.NewProductHandler(usecase, userUseCase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting product service on port: " + cfg.ProductHTTPPort)
	if err := http.ListenAndServe(cfg.ProductHTTPPort, route); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
