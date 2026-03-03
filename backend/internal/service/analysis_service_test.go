package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/igor-trentini/diffable/backend/internal/bitbucket"
	"github.com/igor-trentini/diffable/backend/internal/cache"
	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/igor-trentini/diffable/backend/internal/handler/dto"
	"github.com/igor-trentini/diffable/backend/internal/openai"
	"github.com/igor-trentini/diffable/backend/internal/repository"
)

// --- Mocks ---

type mockBitbucketClient struct {
	commit     *bitbucket.Commit
	diff       string
	pr         *bitbucket.PullRequest
	commits    []bitbucket.Commit
	commitErr  error
	diffErr    error
	prErr      error
	commitsErr error
}

func (m *mockBitbucketClient) GetCommit(_ context.Context, _, _, _ string) (*bitbucket.Commit, error) {
	return m.commit, m.commitErr
}

func (m *mockBitbucketClient) GetCommitDiff(_ context.Context, _, _, _ string) (string, error) {
	return m.diff, m.diffErr
}

func (m *mockBitbucketClient) GetDiffstat(_ context.Context, _, _, _ string) (*bitbucket.PaginatedResponse[bitbucket.DiffstatEntry], error) {
	return nil, nil
}

func (m *mockBitbucketClient) ListCommitsInRange(_ context.Context, _, _, _, _ string) ([]bitbucket.Commit, error) {
	return m.commits, m.commitsErr
}

func (m *mockBitbucketClient) GetPullRequest(_ context.Context, _, _ string, _ int) (*bitbucket.PullRequest, error) {
	return m.pr, m.prErr
}

func (m *mockBitbucketClient) GetPullRequestDiff(_ context.Context, _, _ string, _ int) (string, error) {
	return m.diff, m.diffErr
}

func (m *mockBitbucketClient) GetPullRequestCommits(_ context.Context, _, _ string, _ int) ([]bitbucket.Commit, error) {
	return m.commits, m.commitsErr
}

func (m *mockBitbucketClient) ListRepositories(_ context.Context, _ string) ([]bitbucket.Repository, error) {
	return nil, nil
}

type mockGenerator struct {
	output *openai.GenerationOutput
	err    error
	called bool
}

func (m *mockGenerator) Generate(_ context.Context, _ openai.GenerationInput) (*openai.GenerationOutput, error) {
	m.called = true
	return m.output, m.err
}

func (m *mockGenerator) Refine(_ context.Context, _ openai.RefinementInput) (*openai.GenerationOutput, error) {
	return m.output, m.err
}

type mockRepository struct {
	analysis     *domain.Analysis
	analyses     []domain.Analysis
	total        int
	createErr    error
	getByIDErr   error
	getByHashErr error
	listErr      error
	created      *domain.Analysis
}

func (m *mockRepository) Create(_ context.Context, a *domain.Analysis) error {
	m.created = a
	if m.createErr == nil {
		a.ID = "generated-id"
		a.CreatedAt = time.Now()
		a.UpdatedAt = time.Now()
	}
	return m.createErr
}

func (m *mockRepository) GetByID(_ context.Context, _ string) (*domain.Analysis, error) {
	return m.analysis, m.getByIDErr
}

func (m *mockRepository) GetByDiffHash(_ context.Context, _ string) (*domain.Analysis, error) {
	return m.analysis, m.getByHashErr
}

func (m *mockRepository) List(_ context.Context, _ repository.AnalysisFilter, _, _ int) ([]domain.Analysis, int, error) {
	return m.analyses, m.total, m.listErr
}

type mockCache struct {
	data map[string]string
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string]string)}
}

func (m *mockCache) Get(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *mockCache) Set(key, value string, _ time.Duration) {
	m.data[key] = value
}

// --- Tests ---

func TestAnalyzeCommit_WithHash_CallsBitbucket(t *testing.T) {
	bb := &mockBitbucketClient{
		diff:   "diff --git a/main.go\n+hello",
		commit: &bitbucket.Commit{Hash: "abc123", Message: "fix bug"},
	}
	gen := &mockGenerator{
		output: &openai.GenerationOutput{
			Description: "Generated description",
			Model:       "gpt-4o-mini",
			TokensUsed:  100,
		},
	}
	repo := &mockRepository{getByHashErr: domain.ErrNotFound}
	c := newMockCache()

	svc := NewAnalysisService(bb, gen, repo, c)

	result, err := svc.AnalyzeCommit(context.Background(), &dto.AnalyzeCommitRequest{
		Workspace:  "ws",
		RepoSlug:   "repo",
		CommitHash: "abc123",
	})

	require.NoError(t, err)
	assert.Equal(t, "Generated description", result.GeneratedDesc)
	assert.True(t, gen.called)
	assert.NotNil(t, repo.created)
}

func TestAnalyzeCommit_WithRawDiff_SkipsBitbucket(t *testing.T) {
	bb := &mockBitbucketClient{
		diffErr: assert.AnError, // Would fail if called
	}
	gen := &mockGenerator{
		output: &openai.GenerationOutput{
			Description: "Generated from raw diff",
			Model:       "gpt-4o-mini",
			TokensUsed:  50,
		},
	}
	repo := &mockRepository{getByHashErr: domain.ErrNotFound}
	c := newMockCache()

	svc := NewAnalysisService(bb, gen, repo, c)

	result, err := svc.AnalyzeCommit(context.Background(), &dto.AnalyzeCommitRequest{
		RawDiff: "diff --git a/main.go\n+hello",
	})

	require.NoError(t, err)
	assert.Equal(t, "Generated from raw diff", result.GeneratedDesc)
	assert.True(t, gen.called)
}

func TestAnalyzeCommit_CacheHit_SkipsGenerator(t *testing.T) {
	bb := &mockBitbucketClient{}
	gen := &mockGenerator{}
	repo := &mockRepository{getByHashErr: domain.ErrNotFound}
	c := newMockCache()

	diff := "diff --git a/main.go\n+cached"
	diffHash := cache.DiffCacheKey(diff)
	c.Set(diffHash, "Cached description", time.Hour)

	svc := NewAnalysisService(bb, gen, repo, c)

	result, err := svc.AnalyzeCommit(context.Background(), &dto.AnalyzeCommitRequest{
		RawDiff: diff,
	})

	require.NoError(t, err)
	assert.Equal(t, "Cached description", result.GeneratedDesc)
	assert.False(t, gen.called)
}

func TestAnalyzeCommit_DBHit_SkipsGenerator(t *testing.T) {
	bb := &mockBitbucketClient{}
	gen := &mockGenerator{}
	tokens := 100
	repo := &mockRepository{
		analysis: &domain.Analysis{
			ID:            "existing-id",
			GeneratedDesc: "DB cached description",
			ModelUsed:     "gpt-4o-mini",
			TokensUsed:    &tokens,
		},
		getByHashErr: nil, // DB hit
	}
	c := newMockCache()

	svc := NewAnalysisService(bb, gen, repo, c)

	result, err := svc.AnalyzeCommit(context.Background(), &dto.AnalyzeCommitRequest{
		RawDiff: "diff --git a/main.go\n+db-cached",
	})

	require.NoError(t, err)
	assert.Equal(t, "existing-id", result.ID)
	assert.Equal(t, "DB cached description", result.GeneratedDesc)
	assert.False(t, gen.called)
}

func TestAnalyzeCommit_BitbucketNotFound_ReturnsNotFound(t *testing.T) {
	bb := &mockBitbucketClient{
		diffErr: &bitbucket.NotFoundError{Resource: "commit"},
	}
	gen := &mockGenerator{}
	repo := &mockRepository{getByHashErr: domain.ErrNotFound}
	c := newMockCache()

	svc := NewAnalysisService(bb, gen, repo, c)

	_, err := svc.AnalyzeCommit(context.Background(), &dto.AnalyzeCommitRequest{
		Workspace:  "ws",
		RepoSlug:   "repo",
		CommitHash: "nonexistent",
	})

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestAnalyzeCommit_GeneratorError_ReturnsExternalService(t *testing.T) {
	bb := &mockBitbucketClient{}
	gen := &mockGenerator{err: assert.AnError}
	repo := &mockRepository{getByHashErr: domain.ErrNotFound}
	c := newMockCache()

	svc := NewAnalysisService(bb, gen, repo, c)

	_, err := svc.AnalyzeCommit(context.Background(), &dto.AnalyzeCommitRequest{
		RawDiff: "diff --git a/main.go\n+error",
	})

	assert.ErrorIs(t, err, domain.ErrExternalService)
}

func TestAnalyzeRange_Success(t *testing.T) {
	bb := &mockBitbucketClient{
		commits: []bitbucket.Commit{
			{Hash: "abc123", Message: "commit 1"},
			{Hash: "def456", Message: "commit 2"},
		},
		diff: "diff --git a/main.go\n+range changes",
	}
	gen := &mockGenerator{
		output: &openai.GenerationOutput{
			Description: "Range description",
			Model:       "gpt-4o-mini",
			TokensUsed:  200,
		},
	}
	repo := &mockRepository{getByHashErr: domain.ErrNotFound}
	c := newMockCache()

	svc := NewAnalysisService(bb, gen, repo, c)

	result, err := svc.AnalyzeRange(context.Background(), &dto.AnalyzeRangeRequest{
		Workspace: "ws",
		RepoSlug:  "repo",
		FromHash:  "abc123",
		ToHash:    "def456",
	})

	require.NoError(t, err)
	assert.Equal(t, "Range description", result.GeneratedDesc)
	assert.Equal(t, domain.AnalysisTypeCommitRange, result.AnalysisType)
}

func TestAnalyzePR_WithPRID_Success(t *testing.T) {
	bb := &mockBitbucketClient{
		diff: "diff --git a/main.go\n+pr changes",
		pr: &bitbucket.PullRequest{
			ID:          42,
			Title:       "PR Title",
			Description: "PR Description",
		},
	}
	gen := &mockGenerator{
		output: &openai.GenerationOutput{
			Description: "PR description",
			Model:       "gpt-4o",
			TokensUsed:  300,
		},
	}
	repo := &mockRepository{getByHashErr: domain.ErrNotFound}
	c := newMockCache()

	svc := NewAnalysisService(bb, gen, repo, c)

	result, err := svc.AnalyzePR(context.Background(), &dto.AnalyzePRRequest{
		Workspace: "ws",
		RepoSlug:  "repo",
		PRID:      42,
	})

	require.NoError(t, err)
	assert.Equal(t, "PR description", result.GeneratedDesc)
	assert.Equal(t, domain.AnalysisTypePullRequest, result.AnalysisType)
}

func TestGetAnalysis_Found(t *testing.T) {
	tokens := 100
	repo := &mockRepository{
		analysis: &domain.Analysis{
			ID:            "test-id",
			GeneratedDesc: "Found",
			TokensUsed:    &tokens,
		},
	}

	svc := NewAnalysisService(nil, nil, repo, nil)

	result, err := svc.GetAnalysis(context.Background(), "test-id")

	require.NoError(t, err)
	assert.Equal(t, "test-id", result.ID)
}

func TestGetAnalysis_NotFound(t *testing.T) {
	repo := &mockRepository{getByIDErr: domain.ErrNotFound}

	svc := NewAnalysisService(nil, nil, repo, nil)

	_, err := svc.GetAnalysis(context.Background(), "nonexistent")

	assert.ErrorIs(t, err, domain.ErrNotFound)
}
