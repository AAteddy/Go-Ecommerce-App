package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	paymentHttp "github.com/AAteddy/go-ecommerce-app/internal/payment/delivery/http"
	"github.com/AAteddy/go-ecommerce-app/internal/payment/repository"
	"github.com/AAteddy/go-ecommerce-app/internal/payment/usecase"
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

	log.Info("Initializing Payment Service")

	dbConn, err := db.Migrate(cfg.DBURL)
	if err != nil {
		log.Error("Failed to connect and migrate database", "error", err)
	}

	tokenStore := redis.NewRedisTokenStore(cfg.RedisAddr)
	userRepo := userRepo.NewPostgresUserRepository(dbConn)
	userUseCase := UserUC.NewUserUseCase(userRepo, tokenStore, nil, cfg.JWTSecret, log)

	repo := repository.NewPostgresPaymentRepository(dbConn)
	usecase := usecase.NewPaymentUseCase(repo, log)

	route := chi.NewRouter()
	handler := paymentHttp.NewPaymentHandler(usecase, userUseCase, log)
	handler.RegisterRoutes(route)

	log.Info("Starting payment service on port: " + cfg.PaymentHTTPPort)
	if err := http.ListenAndServe(cfg.PaymentHTTPPort, route); err != nil {
		log.Fatal("Failed to start payment service", "error", err)
	}

}
