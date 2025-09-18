package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	orderHttp "github.com/AAteddy/go-ecommerce-app/internal/order/delivery/http"
	"github.com/AAteddy/go-ecommerce-app/internal/order/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/order/usecase"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	productRepo "github.com/AAteddy/go-ecommerce-app/internal/product/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/user/infrastructure/redis"
	userRepo "github.com/AAteddy/go-ecommerce-app/internal/user/repository"
	UserUC "github.com/AAteddy/go-ecommerce-app/internal/user/usecase"
)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing Order Service")

	// if os.Getenv("RUN_MIGRATIONS") == "true" {
	// 	dbConn, err := db.NewPostgresDB(cfg.DBURL)
	// 	if err != nil {
	// 		log.Fatal("Failed to connect to database", "error", err)
	// 	}
	// 	if err := orderDB.Migrate(dbConn); err != nil {
	// 		log.Fatal("Failed to run migrations", "error", err)
	// 	}
	// 	log.Info("Order schema migrations completed")
	// }

	dbConn, err := db.NewPostgresDB(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	tokenStore := redis.NewRedisTokenStore(cfg.RedisAddr)
	userRepo := userRepo.NewPostgresUserRepository(dbConn)
	userUseCase := UserUC.NewUserUseCase(userRepo, tokenStore, nil, cfg.JWTSecret, log) // nil for repo, tokenStore, kafka (stubbed)

	productRepo := productRepo.NewPostgresProductRepository(dbConn)
	repo := repository.NewPostgresOrderRepository(dbConn)
	usecase := usecase.NewOrderUseCase(repo, productRepo, log)

	route := chi.NewRouter()
	handler := orderHttp.NewOrderHandler(usecase, userUseCase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting order service on port: " + cfg.OrderHTTPPort)
	if err := http.ListenAndServe(cfg.OrderHTTPPort, route); err != nil {
		log.Fatal("Failed to start order service", "error", err)
	}

}
