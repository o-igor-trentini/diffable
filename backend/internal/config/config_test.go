package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("FRONTEND_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("SHUTDOWN_TIMEOUT")
	os.Unsetenv("BITBUCKET_BASE_URL")
	os.Unsetenv("BITBUCKET_EMAIL")
	os.Unsetenv("BITBUCKET_API_TOKEN")
	os.Unsetenv("BITBUCKET_TIMEOUT")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("OPENAI_DEFAULT_MODEL")
	os.Unsetenv("OPENAI_COMPLEX_MODEL")
	os.Unsetenv("OPENAI_MAX_TOKENS")
	os.Unsetenv("OPENAI_TEMPERATURE")
	os.Unsetenv("OPENAI_TOKEN_THRESHOLD")
	os.Unsetenv("CACHE_TTL")

	cfg := Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "", cfg.DatabaseURL)
	assert.Equal(t, "http://localhost:3000", cfg.FrontendURL)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 30*time.Second, cfg.ShutdownTimeout)
	assert.Equal(t, "https://api.bitbucket.org/2.0", cfg.BitbucketBaseURL)
	assert.Equal(t, "", cfg.BitbucketEmail)
	assert.Equal(t, "", cfg.BitbucketAPIToken)
	assert.Equal(t, 30*time.Second, cfg.BitbucketTimeout)
	assert.Equal(t, "", cfg.OpenAIAPIKey)
	assert.Equal(t, "gpt-4o-mini", cfg.OpenAIDefaultModel)
	assert.Equal(t, "gpt-4o", cfg.OpenAIComplexModel)
	assert.Equal(t, 1024, cfg.OpenAIMaxTokens)
	assert.InDelta(t, 0.3, cfg.OpenAITemperature, 0.001)
	assert.Equal(t, 4000, cfg.OpenAITokenThreshold)
	assert.Equal(t, 24*time.Hour, cfg.CacheTTL)
}

func TestLoad_Overrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	t.Setenv("FRONTEND_URL", "http://example.com")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("SHUTDOWN_TIMEOUT", "10s")
	t.Setenv("BITBUCKET_BASE_URL", "https://custom.bitbucket.org/2.0")
	t.Setenv("BITBUCKET_EMAIL", "dev@company.com")
	t.Setenv("BITBUCKET_API_TOKEN", "my-token")
	t.Setenv("BITBUCKET_TIMEOUT", "15s")
	t.Setenv("OPENAI_API_KEY", "sk-test-key")
	t.Setenv("OPENAI_DEFAULT_MODEL", "gpt-4o-mini-custom")
	t.Setenv("OPENAI_COMPLEX_MODEL", "gpt-4o-custom")
	t.Setenv("OPENAI_MAX_TOKENS", "2048")
	t.Setenv("OPENAI_TEMPERATURE", "0.5")
	t.Setenv("OPENAI_TOKEN_THRESHOLD", "8000")
	t.Setenv("CACHE_TTL", "12h")

	cfg := Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "postgres://test:test@localhost:5432/test", cfg.DatabaseURL)
	assert.Equal(t, "http://example.com", cfg.FrontendURL)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
	assert.Equal(t, "https://custom.bitbucket.org/2.0", cfg.BitbucketBaseURL)
	assert.Equal(t, "dev@company.com", cfg.BitbucketEmail)
	assert.Equal(t, "my-token", cfg.BitbucketAPIToken)
	assert.Equal(t, 15*time.Second, cfg.BitbucketTimeout)
	assert.Equal(t, "sk-test-key", cfg.OpenAIAPIKey)
	assert.Equal(t, "gpt-4o-mini-custom", cfg.OpenAIDefaultModel)
	assert.Equal(t, "gpt-4o-custom", cfg.OpenAIComplexModel)
	assert.Equal(t, 2048, cfg.OpenAIMaxTokens)
	assert.InDelta(t, 0.5, cfg.OpenAITemperature, 0.001)
	assert.Equal(t, 8000, cfg.OpenAITokenThreshold)
	assert.Equal(t, 12*time.Hour, cfg.CacheTTL)
}

func TestParseDuration_Invalid(t *testing.T) {
	d := parseDuration("invalid", 5*time.Second)
	assert.Equal(t, 5*time.Second, d)
}
