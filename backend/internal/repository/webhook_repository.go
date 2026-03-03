package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WebhookLog struct {
	ID           string
	EventKey     string
	Payload      []byte
	Status       string
	AnalysisID   *string
	ErrorMessage string
	CreatedAt    time.Time
}

type WebhookRepository interface {
	Create(ctx context.Context, log *WebhookLog) error
	UpdateStatus(ctx context.Context, id, status string, analysisID *string, errorMsg string) error
}

type postgresWebhookRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresWebhookRepository(pool *pgxpool.Pool) WebhookRepository {
	return &postgresWebhookRepository{pool: pool}
}

func (r *postgresWebhookRepository) Create(ctx context.Context, log *WebhookLog) error {
	query := `
		INSERT INTO webhook_logs (event_key, payload, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query,
		log.EventKey,
		log.Payload,
		log.Status,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *postgresWebhookRepository) UpdateStatus(ctx context.Context, id, status string, analysisID *string, errorMsg string) error {
	query := `
		UPDATE webhook_logs
		SET status = $2, analysis_id = $3, error_message = $4
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id, status, analysisID, nullableString(errorMsg))
	if err != nil {
		return fmt.Errorf("updating webhook log status: %w", err)
	}
	return nil
}
