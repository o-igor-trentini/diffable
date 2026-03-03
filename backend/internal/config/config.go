package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	FrontendURL     string
	LogLevel        string
	ShutdownTimeout time.Duration

	BitbucketBaseURL  string
	BitbucketEmail    string
	BitbucketAPIToken string
	BitbucketTimeout  time.Duration

	OpenAIAPIKey         string
	OpenAIDefaultModel   string
	OpenAIComplexModel   string
	OpenAIMaxTokens      int
	OpenAITemperature    float64
	OpenAITokenThreshold int
	CacheTTL             time.Duration

	RateLimitRPM int
	DBMaxConns   int
	DBMinConns   int
}

func Load() (*Config, error) {
	// Load .env file if present (silently ignored when absent, e.g. inside Docker)
	_ = godotenv.Load()

	var missing []string
	requireEnv := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	cfg := &Config{
		// Required (secrets)
		DatabaseURL:     requireEnv("DATABASE_URL"),
		BitbucketEmail:  requireEnv("BITBUCKET_EMAIL"),
		BitbucketAPIToken: requireEnv("BITBUCKET_API_TOKEN"),
		OpenAIAPIKey:    requireEnv("OPENAI_API_KEY"),

		// Optional (safe defaults)
		Port:            getEnv("PORT", "8080"),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:3000"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout: parseDuration(getEnv("SHUTDOWN_TIMEOUT", "30s"), 30*time.Second),

		BitbucketBaseURL: getEnv("BITBUCKET_BASE_URL", "https://api.bitbucket.org/2.0"),
		BitbucketTimeout: parseDuration(getEnv("BITBUCKET_TIMEOUT", "30s"), 30*time.Second),

		OpenAIDefaultModel:   getEnv("OPENAI_DEFAULT_MODEL", "gpt-4o-mini"),
		OpenAIComplexModel:   getEnv("OPENAI_COMPLEX_MODEL", "gpt-4o"),
		OpenAIMaxTokens:      parseInt(getEnv("OPENAI_MAX_TOKENS", "1024"), 1024),
		OpenAITemperature:    parseFloat(getEnv("OPENAI_TEMPERATURE", "0.3"), 0.3),
		OpenAITokenThreshold: parseInt(getEnv("OPENAI_TOKEN_THRESHOLD", "4000"), 4000),
		CacheTTL:             parseDuration(getEnv("CACHE_TTL", "24h"), 24*time.Hour),

		RateLimitRPM: parseInt(getEnv("RATE_LIMIT_RPM", "60"), 60),
		DBMaxConns:   parseInt(getEnv("DB_MAX_CONNS", "25"), 25),
		DBMinConns:   parseInt(getEnv("DB_MIN_CONNS", "5"), 5),
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
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

func parseInt(s string, fallback int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}

func parseFloat(s string, fallback float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fallback
	}
	return v
}
