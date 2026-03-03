package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	oai "github.com/sashabaranov/go-openai"

	"github.com/igor-trentini/diffable/backend/internal/bitbucket"
	"github.com/igor-trentini/diffable/backend/internal/cache"
	"github.com/igor-trentini/diffable/backend/internal/config"
	"github.com/igor-trentini/diffable/backend/internal/handler"
	genoai "github.com/igor-trentini/diffable/backend/internal/openai"
	"github.com/igor-trentini/diffable/backend/internal/repository"
	"github.com/igor-trentini/diffable/backend/internal/server"
	"github.com/igor-trentini/diffable/backend/internal/service"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	slog.Info("starting diffable backend", "port", cfg.Port)

	// Database connection pool with tuning
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to parse database URL", "error", err)
		os.Exit(1)
	}

	poolConfig.MaxConns = int32(cfg.DBMaxConns)
	poolConfig.MinConns = int32(cfg.DBMinConns)
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database",
		"max_conns", cfg.DBMaxConns,
		"min_conns", cfg.DBMinConns,
	)

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

	// Services
	analysisSvc := service.NewAnalysisService(bbClient, generator, analysisRepo, diffCache)
	refinementSvc := service.NewRefinementService(generator, analysisRepo)
	historySvc := service.NewHistoryService(analysisRepo)

	// Handlers
	bbHandler := handler.NewBitbucketHandler(bbClient, diffCache)
	webhookRepo := repository.NewPostgresWebhookRepository(pool)
	whHandler := handler.NewWebhookHandler(analysisSvc, webhookRepo)

	// Server
	srv := server.New(pool, cfg.FrontendURL, cfg.RateLimitRPM, analysisSvc, refinementSvc, historySvc, bbHandler, whHandler)

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
	sig := <-quit

	slog.Info("received shutdown signal", "signal", sig.String())
	slog.Info("shutting down server, draining active connections...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	pool.Close()
	slog.Info("database connections closed")
	slog.Info("server stopped gracefully")
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
