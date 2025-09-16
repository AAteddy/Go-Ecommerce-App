package config

import (
	"os"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/joho/godotenv"
)

// Config holds the environment variable's configuration values for the application.
type Config struct {
	UserHTTPPort      string
	ProductHTTPPort   string
	OrderHTTPPort     string
	PaymentHTTPPort   string
	InventoryHTTPPort string
	DBURL             string
	KafkaBroker       string
	RedisAddr         string
	JWTSecret         string
}

// LoadConfig loads the environment variables from the .env file.
func LoadConfig() Config {
	// Load .env file if it exists
	// Only load .env if not running in Docker
	if _, exists := os.LookupEnv("DOCKER_ENV"); !exists {
		if err := godotenv.Load(); err != nil {
			// Log warning but continue, as env vars might be set externally
			// Logging not used here to avoid dependency; will log in main
		}
	}
	cfg := Config{
		UserHTTPPort:      getEnv("UserHTTP_PORT", ":8080"),
		ProductHTTPPort:   getEnv("ProductHTTP_PORT", ":8081"),
		OrderHTTPPort:     getEnv("OrderHTTP_PORT", ":8082"),
		PaymentHTTPPort:   getEnv("PaymentHTTP_PORT", ":8084"),
		InventoryHTTPPort: getEnv("InventoryHTTP_PORT", ":8083"),
		DBURL:             getEnv("DB_URL", "postgres://user:password@postgres:5432/ecommerce_db?sslmode=disable"),
		KafkaBroker:       getEnv("KAFKA_BROKER", "kafka:9092"),
		RedisAddr:         getEnv("REDIS_ADDR", "redis:6379"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
	}

	// Validate required fields
	if cfg.JWTSecret == "" {
		panic(errors.New("JWT_SECRET is required"))
	}
	if cfg.DBURL == "" {
		panic(errors.New("DB_URL is required"))
	}
	if cfg.KafkaBroker == "" {
		panic(errors.New("KAFKA_BROKER is required"))
	}
	if cfg.RedisAddr == "" {
		panic(errors.New("REDIS_ADDR is required"))
	}

	return cfg
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
