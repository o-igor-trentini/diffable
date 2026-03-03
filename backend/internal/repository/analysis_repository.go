package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/igor-trentini/diffable/backend/internal/domain"
)

type AnalysisFilter struct {
	Type string
}

type AnalysisRepository interface {
	Create(ctx context.Context, analysis *domain.Analysis) error
	GetByID(ctx context.Context, id string) (*domain.Analysis, error)
	GetByDiffHash(ctx context.Context, hash string) (*domain.Analysis, error)
	List(ctx context.Context, filter AnalysisFilter, offset, limit int) ([]domain.Analysis, int, error)
}

type postgresAnalysisRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAnalysisRepository(pool *pgxpool.Pool) AnalysisRepository {
	return &postgresAnalysisRepository{pool: pool}
}

func (r *postgresAnalysisRepository) Create(ctx context.Context, analysis *domain.Analysis) error {
	query := `
		INSERT INTO analyses (
			analysis_type, status, workspace, repo_slug,
			commit_hash, from_hash, to_hash, pr_id,
			raw_diff, diff_hash, generated_desc,
			model_used, tokens_used, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		analysis.AnalysisType,
		analysis.Status,
		nullableString(analysis.Workspace),
		nullableString(analysis.RepoSlug),
		nullableString(analysis.CommitHash),
		nullableString(analysis.FromHash),
		nullableString(analysis.ToHash),
		analysis.PrID,
		nullableString(analysis.RawDiff),
		nullableString(analysis.DiffHash),
		nullableString(analysis.GeneratedDesc),
		nullableString(analysis.ModelUsed),
		analysis.TokensUsed,
		nullableString(analysis.ErrorMessage),
	).Scan(&analysis.ID, &analysis.CreatedAt, &analysis.UpdatedAt)
}

func (r *postgresAnalysisRepository) GetByID(ctx context.Context, id string) (*domain.Analysis, error) {
	query := `
		SELECT id, analysis_type, status, workspace, repo_slug,
			commit_hash, from_hash, to_hash, pr_id,
			raw_diff, diff_hash, generated_desc,
			model_used, tokens_used, error_message,
			created_at, updated_at
		FROM analyses WHERE id = $1`

	a := &domain.Analysis{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.AnalysisType, &a.Status,
		&a.Workspace, &a.RepoSlug,
		&a.CommitHash, &a.FromHash, &a.ToHash, &a.PrID,
		&a.RawDiff, &a.DiffHash, &a.GeneratedDesc,
		&a.ModelUsed, &a.TokensUsed, &a.ErrorMessage,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting analysis by id: %w", err)
	}
	return a, nil
}

func (r *postgresAnalysisRepository) GetByDiffHash(ctx context.Context, hash string) (*domain.Analysis, error) {
	query := `
		SELECT id, analysis_type, status, workspace, repo_slug,
			commit_hash, from_hash, to_hash, pr_id,
			raw_diff, diff_hash, generated_desc,
			model_used, tokens_used, error_message,
			created_at, updated_at
		FROM analyses
		WHERE diff_hash = $1 AND status = 'completed'
		ORDER BY created_at DESC
		LIMIT 1`

	a := &domain.Analysis{}
	err := r.pool.QueryRow(ctx, query, hash).Scan(
		&a.ID, &a.AnalysisType, &a.Status,
		&a.Workspace, &a.RepoSlug,
		&a.CommitHash, &a.FromHash, &a.ToHash, &a.PrID,
		&a.RawDiff, &a.DiffHash, &a.GeneratedDesc,
		&a.ModelUsed, &a.TokensUsed, &a.ErrorMessage,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting analysis by diff hash: %w", err)
	}
	return a, nil
}

func (r *postgresAnalysisRepository) List(ctx context.Context, filter AnalysisFilter, offset, limit int) ([]domain.Analysis, int, error) {
	countQuery := "SELECT COUNT(*) FROM analyses"
	listQuery := `
		SELECT id, analysis_type, status, workspace, repo_slug,
			commit_hash, from_hash, to_hash, pr_id,
			raw_diff, diff_hash, generated_desc,
			model_used, tokens_used, error_message,
			created_at, updated_at
		FROM analyses`

	var args []any
	where := ""
	if filter.Type != "" {
		where = " WHERE analysis_type = $1"
		args = append(args, filter.Type)
	}

	countQuery += where
	listQuery += where + " ORDER BY created_at DESC"

	if limit > 0 {
		if len(args) > 0 {
			listQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		} else {
			listQuery += " LIMIT $1 OFFSET $2"
		}
		args = append(args, limit, offset)
	}

	var total int
	countArgs := args
	if limit > 0 {
		countArgs = args[:len(args)-2]
	}
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting analyses: %w", err)
	}

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing analyses: %w", err)
	}
	defer rows.Close()

	var analyses []domain.Analysis
	for rows.Next() {
		var a domain.Analysis
		if err := rows.Scan(
			&a.ID, &a.AnalysisType, &a.Status,
			&a.Workspace, &a.RepoSlug,
			&a.CommitHash, &a.FromHash, &a.ToHash, &a.PrID,
			&a.RawDiff, &a.DiffHash, &a.GeneratedDesc,
			&a.ModelUsed, &a.TokensUsed, &a.ErrorMessage,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning analysis: %w", err)
		}
		analyses = append(analyses, a)
	}

	return analyses, total, nil
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
