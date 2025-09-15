package main

import (
	"net/http"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/config"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/db"
	"github.com/AAteddy/go-ecommerce-app/internal/pkg/logging"
	productHttp "github.com/AAteddy/go-ecommerce-app/internal/product/delivery/http"
	"github.com/AAteddy/go-ecommerce-app/internal/product/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/product/usecase"
	"github.com/go-chi/chi/v5"
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
	usecase := usecase.NewProductUseCase(repo, log)

	route := chi.NewRouter()
	handler := productHttp.NewProductHandler(usecase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting product service on port: " + cfg.ProductHTTPPort)
	if err := http.ListenAndServe(cfg.ProductHTTPPort, route); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
