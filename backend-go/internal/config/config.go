package config

import "os"

type Config struct {
	AppEnv           string
	Port             string
	DatabaseURL      string
	JWTSecret        string
	PythonServiceURL string
}

func Load() Config {
	return Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://careerflow:careerflow@localhost:5432/careerflow?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "change-me"),
		PythonServiceURL: getEnv("PYTHON_SERVICE_URL", "http://localhost:8000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
