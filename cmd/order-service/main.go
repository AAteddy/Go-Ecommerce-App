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
)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing Order Service")

	dbConn, err := db.Migrate(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	repo := repository.NewPostgresOrderRepository(dbConn)
	usecase := usecase.NewOrderUseCase(repo, log)

	route := chi.NewRouter()
	handler := orderHttp.NewOrderHandler(usecase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting order service on port: " + cfg.OrderHTTPPort)
	if err := http.ListenAndServe(cfg.OrderHTTPPort, route); err != nil {
		log.Fatal("Failed to start order service", "error", err)
	}

}
