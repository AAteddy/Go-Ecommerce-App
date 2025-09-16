package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	inventoryHttp "github.com/AAteddy/go-ecommerce-app/internal/inventory/delivery/http"
	"github.com/AAteddy/go-ecommerce-app/internal/inventory/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/inventory/usecase"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
)

func main() {
	cfg := config.LoadConfig()
	log := logging.Init()

	log.Info("Initializing Inventory Service")

	dbConn, err := db.Migrate(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	repo := repository.NewPostgresInventoryRepository(dbConn)
	usecase := usecase.NewInventoryUseCase(repo, log)

	route := chi.NewRouter()
	handler := inventoryHttp.NewInventoryHandler(usecase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting inventory service on port: " + cfg.InventoryHTTPPort)
	if err := http.ListenAndServe(cfg.InventoryHTTPPort, route); err != nil {
		log.Fatal("Failed to start inventory service", "error", err)
	}
}
