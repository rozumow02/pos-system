package config

import (
	"os"
	"strconv"
)

type Config struct {
	BackendPort       string
	DatabaseURL       string
	LogLevel          string
	FrontendOrigin    string
	LowStockThreshold int
	MigrationsPath    string
	DBConnectRetries  int
}

func Load() Config {
	return Config{
		BackendPort:       getEnv("BACKEND_PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://posuser:pospass@localhost:5432/posdb?sslmode=disable"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		FrontendOrigin:    getEnv("FRONTEND_ORIGIN", "http://localhost"),
		LowStockThreshold: getEnvAsInt("LOW_STOCK_THRESHOLD", 5),
		MigrationsPath:    getEnv("MIGRATIONS_PATH", "./migrations"),
		DBConnectRetries:  getEnvAsInt("DB_CONNECT_RETRIES", 20),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
