package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	oai "github.com/sashabaranov/go-openai"

	"github.com/igor-trentini/diffable/backend/internal/bitbucket"
	"github.com/igor-trentini/diffable/backend/internal/cache"
	"github.com/igor-trentini/diffable/backend/internal/config"
	genoai "github.com/igor-trentini/diffable/backend/internal/openai"
	"github.com/igor-trentini/diffable/backend/internal/repository"
	"github.com/igor-trentini/diffable/backend/internal/server"
	"github.com/igor-trentini/diffable/backend/internal/service"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	slog.Info("starting diffable backend", "port", cfg.Port)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Bitbucket client
	bbClient := bitbucket.NewClient(bitbucket.Config{
		BaseURL:  cfg.BitbucketBaseURL,
		Email:    cfg.BitbucketEmail,
		APIToken: cfg.BitbucketAPIToken,
		Timeout:  cfg.BitbucketTimeout,
	})

	// In-memory cache
	diffCache := cache.NewInMemoryCache()

	// OpenAI generator
	oaiClient := oai.NewClient(cfg.OpenAIAPIKey)
	generator := genoai.NewGenerator(oaiClient, diffCache, genoai.GeneratorConfig{
		DefaultModel:   cfg.OpenAIDefaultModel,
		ComplexModel:   cfg.OpenAIComplexModel,
		MaxTokens:      cfg.OpenAIMaxTokens,
		Temperature:    float32(cfg.OpenAITemperature),
		TokenThreshold: cfg.OpenAITokenThreshold,
		CacheTTL:       cfg.CacheTTL,
	})

	// Repository
	analysisRepo := repository.NewPostgresAnalysisRepository(pool)

	// Service
	analysisSvc := service.NewAnalysisService(bbClient, generator, analysisRepo, diffCache)

	// Server
	srv := server.New(pool, cfg.FrontendURL, analysisSvc)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: srv,
	}

	go func() {
		slog.Info("server listening", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}

func setupLogger(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}

	slog.Info("migrations applied successfully")
	return nil
}
