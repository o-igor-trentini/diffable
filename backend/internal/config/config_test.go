package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var requiredVars = map[string]string{
	"DATABASE_URL":        "postgres://test:test@localhost:5432/test",
	"BITBUCKET_EMAIL":     "dev@company.com",
	"BITBUCKET_API_TOKEN": "my-token",
	"OPENAI_API_KEY":      "sk-test-key",
}

// allEnvKeys lists every env var the config reads, for cleanup.
var allEnvKeys = []string{
	"DATABASE_URL", "BITBUCKET_EMAIL", "BITBUCKET_API_TOKEN", "OPENAI_API_KEY",
	"PORT", "FRONTEND_URL", "LOG_LEVEL", "SHUTDOWN_TIMEOUT",
	"BITBUCKET_BASE_URL", "BITBUCKET_TIMEOUT",
	"OPENAI_DEFAULT_MODEL", "OPENAI_COMPLEX_MODEL", "OPENAI_MAX_TOKENS",
	"OPENAI_TEMPERATURE", "OPENAI_TOKEN_THRESHOLD", "CACHE_TTL",
	"RATE_LIMIT_RPM", "DB_MAX_CONNS", "DB_MIN_CONNS",
}

func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range allEnvKeys {
		t.Setenv(k, "")
		os.Unsetenv(k)
	}
}

func setRequired(t *testing.T) {
	t.Helper()
	for k, v := range requiredVars {
		t.Setenv(k, v)
	}
}

func TestLoad_MissingRequiredVars(t *testing.T) {
	clearEnv(t)

	_, err := Load()
	require.Error(t, err)

	msg := err.Error()
	assert.Contains(t, msg, "DATABASE_URL")
	assert.Contains(t, msg, "BITBUCKET_EMAIL")
	assert.Contains(t, msg, "BITBUCKET_API_TOKEN")
	assert.Contains(t, msg, "OPENAI_API_KEY")
}

func TestLoad_OptionalDefaults(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "http://localhost:3000", cfg.FrontendURL)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 30*time.Second, cfg.ShutdownTimeout)
	assert.Equal(t, "https://api.bitbucket.org/2.0", cfg.BitbucketBaseURL)
	assert.Equal(t, 30*time.Second, cfg.BitbucketTimeout)
	assert.Equal(t, "gpt-4o-mini", cfg.OpenAIDefaultModel)
	assert.Equal(t, "gpt-4o", cfg.OpenAIComplexModel)
	assert.Equal(t, 1024, cfg.OpenAIMaxTokens)
	assert.InDelta(t, 0.3, cfg.OpenAITemperature, 0.001)
	assert.Equal(t, 4000, cfg.OpenAITokenThreshold)
	assert.Equal(t, 24*time.Hour, cfg.CacheTTL)
	assert.Equal(t, 60, cfg.RateLimitRPM)
	assert.Equal(t, 25, cfg.DBMaxConns)
	assert.Equal(t, 5, cfg.DBMinConns)
}

func TestLoad_AllOverrides(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	t.Setenv("PORT", "9090")
	t.Setenv("FRONTEND_URL", "http://example.com")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("SHUTDOWN_TIMEOUT", "10s")
	t.Setenv("BITBUCKET_BASE_URL", "https://custom.bitbucket.org/2.0")
	t.Setenv("BITBUCKET_TIMEOUT", "15s")
	t.Setenv("OPENAI_DEFAULT_MODEL", "gpt-4o-mini-custom")
	t.Setenv("OPENAI_COMPLEX_MODEL", "gpt-4o-custom")
	t.Setenv("OPENAI_MAX_TOKENS", "2048")
	t.Setenv("OPENAI_TEMPERATURE", "0.5")
	t.Setenv("OPENAI_TOKEN_THRESHOLD", "8000")
	t.Setenv("CACHE_TTL", "12h")
	t.Setenv("RATE_LIMIT_RPM", "120")
	t.Setenv("DB_MAX_CONNS", "50")
	t.Setenv("DB_MIN_CONNS", "10")

	cfg, err := Load()
	require.NoError(t, err)

	// Required vars
	assert.Equal(t, "postgres://test:test@localhost:5432/test", cfg.DatabaseURL)
	assert.Equal(t, "dev@company.com", cfg.BitbucketEmail)
	assert.Equal(t, "my-token", cfg.BitbucketAPIToken)
	assert.Equal(t, "sk-test-key", cfg.OpenAIAPIKey)

	// Overridden optional vars
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "http://example.com", cfg.FrontendURL)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
	assert.Equal(t, "https://custom.bitbucket.org/2.0", cfg.BitbucketBaseURL)
	assert.Equal(t, 15*time.Second, cfg.BitbucketTimeout)
	assert.Equal(t, "gpt-4o-mini-custom", cfg.OpenAIDefaultModel)
	assert.Equal(t, "gpt-4o-custom", cfg.OpenAIComplexModel)
	assert.Equal(t, 2048, cfg.OpenAIMaxTokens)
	assert.InDelta(t, 0.5, cfg.OpenAITemperature, 0.001)
	assert.Equal(t, 8000, cfg.OpenAITokenThreshold)
	assert.Equal(t, 12*time.Hour, cfg.CacheTTL)
	assert.Equal(t, 120, cfg.RateLimitRPM)
	assert.Equal(t, 50, cfg.DBMaxConns)
	assert.Equal(t, 10, cfg.DBMinConns)
}

func TestParseDuration_Invalid(t *testing.T) {
	d := parseDuration("invalid", 5*time.Second)
	assert.Equal(t, 5*time.Second, d)
}
