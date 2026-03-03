//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/igor-trentini/diffable/backend/internal/domain"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()
	dsn := "postgres://postgres:postgres@localhost:5432/diffable_test?sslmode=disable"

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	// Clean up before test
	_, err = pool.Exec(ctx, "DELETE FROM analyses")
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestPostgresRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewPostgresAnalysisRepository(pool)
	ctx := context.Background()

	tokens := 100
	analysis := &domain.Analysis{
		AnalysisType:  domain.AnalysisTypeSingleCommit,
		Status:        domain.AnalysisStatusCompleted,
		Workspace:     "test-ws",
		RepoSlug:      "test-repo",
		CommitHash:    "abc123",
		DiffHash:      fmt.Sprintf("hash-%d", time.Now().UnixNano()),
		GeneratedDesc: "Test description",
		ModelUsed:     "gpt-4o-mini",
		TokensUsed:    &tokens,
	}

	err := repo.Create(ctx, analysis)
	require.NoError(t, err)
	assert.NotEmpty(t, analysis.ID)
	assert.False(t, analysis.CreatedAt.IsZero())
}

func TestPostgresRepository_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewPostgresAnalysisRepository(pool)
	ctx := context.Background()

	tokens := 100
	analysis := &domain.Analysis{
		AnalysisType:  domain.AnalysisTypeSingleCommit,
		Status:        domain.AnalysisStatusCompleted,
		Workspace:     "test-ws",
		RepoSlug:      "test-repo",
		DiffHash:      fmt.Sprintf("hash-%d", time.Now().UnixNano()),
		GeneratedDesc: "Test description",
		ModelUsed:     "gpt-4o-mini",
		TokensUsed:    &tokens,
	}

	require.NoError(t, repo.Create(ctx, analysis))

	found, err := repo.GetByID(ctx, analysis.ID)
	require.NoError(t, err)
	assert.Equal(t, analysis.ID, found.ID)
	assert.Equal(t, "Test description", found.GeneratedDesc)
}

func TestPostgresRepository_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewPostgresAnalysisRepository(pool)

	_, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPostgresRepository_GetByDiffHash(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewPostgresAnalysisRepository(pool)
	ctx := context.Background()

	hash := fmt.Sprintf("hash-%d", time.Now().UnixNano())
	tokens := 100
	analysis := &domain.Analysis{
		AnalysisType:  domain.AnalysisTypeSingleCommit,
		Status:        domain.AnalysisStatusCompleted,
		DiffHash:      hash,
		GeneratedDesc: "Hash found",
		ModelUsed:     "gpt-4o-mini",
		TokensUsed:    &tokens,
	}

	require.NoError(t, repo.Create(ctx, analysis))

	found, err := repo.GetByDiffHash(ctx, hash)
	require.NoError(t, err)
	assert.Equal(t, "Hash found", found.GeneratedDesc)
}

func TestPostgresRepository_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewPostgresAnalysisRepository(pool)
	ctx := context.Background()

	tokens := 50
	for i := 0; i < 3; i++ {
		a := &domain.Analysis{
			AnalysisType:  domain.AnalysisTypeSingleCommit,
			Status:        domain.AnalysisStatusCompleted,
			DiffHash:      fmt.Sprintf("list-hash-%d-%d", time.Now().UnixNano(), i),
			GeneratedDesc: fmt.Sprintf("Description %d", i),
			ModelUsed:     "gpt-4o-mini",
			TokensUsed:    &tokens,
		}
		require.NoError(t, repo.Create(ctx, a))
	}

	analyses, total, err := repo.List(ctx, AnalysisFilter{}, 0, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, analyses, 2)
}
