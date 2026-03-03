package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	DatabaseURL     string
	FrontendURL     string
	LogLevel        string
	ShutdownTimeout time.Duration
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:3000"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout: parseDuration(getEnv("SHUTDOWN_TIMEOUT", "30s"), 30*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return fallback
	}
	return d
}
