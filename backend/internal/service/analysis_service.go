package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/igor-trentini/diffable/backend/internal/bitbucket"
	"github.com/igor-trentini/diffable/backend/internal/cache"
	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/igor-trentini/diffable/backend/internal/handler/dto"
	"github.com/igor-trentini/diffable/backend/internal/openai"
	"github.com/igor-trentini/diffable/backend/internal/repository"
)

type AnalysisService interface {
	AnalyzeCommit(ctx context.Context, req *dto.AnalyzeCommitRequest) (*domain.Analysis, error)
	AnalyzeRange(ctx context.Context, req *dto.AnalyzeRangeRequest) (*domain.Analysis, error)
	AnalyzePR(ctx context.Context, req *dto.AnalyzePRRequest) (*domain.Analysis, error)
	GetAnalysis(ctx context.Context, id string) (*domain.Analysis, error)
}

type analysisService struct {
	bbClient   bitbucket.Client
	generator  openai.DescriptionGenerator
	repository repository.AnalysisRepository
	cache      cache.Cache
}

func NewAnalysisService(
	bbClient bitbucket.Client,
	generator openai.DescriptionGenerator,
	repo repository.AnalysisRepository,
	c cache.Cache,
) AnalysisService {
	return &analysisService{
		bbClient:   bbClient,
		generator:  generator,
		repository: repo,
		cache:      c,
	}
}

func defaultLevel(level string) string {
	if level == "" {
		return "functional"
	}
	return level
}

func applyOverrides(input *openai.GenerationInput, overrides *dto.GenerationOverrides) {
	if overrides == nil {
		return
	}
	input.MaxTokensOverride = overrides.MaxTokens
	input.TemperatureOverride = overrides.Temperature
	input.ModelOverride = overrides.Model
}

func hasOverrides(overrides *dto.GenerationOverrides) bool {
	if overrides == nil {
		return false
	}
	if overrides.MaxTokens != nil {
		return true
	}
	if overrides.Temperature != nil {
		return true
	}
	if overrides.Model != nil && *overrides.Model != "auto" {
		return true
	}
	return false
}

func (s *analysisService) AnalyzeCommit(ctx context.Context, req *dto.AnalyzeCommitRequest) (*domain.Analysis, error) {
	diff := req.RawDiff
	var commitMessages []string
	level := defaultLevel(req.Level)

	if diff == "" {
		var err error
		diff, err = s.bbClient.GetCommitDiff(ctx, req.Workspace, req.RepoSlug, req.CommitHash)
		if err != nil {
			return nil, s.mapExternalError(err)
		}

		commit, err := s.bbClient.GetCommit(ctx, req.Workspace, req.RepoSlug, req.CommitHash)
		if err != nil {
			slog.Warn("failed to get commit message, proceeding without it", "error", err)
		} else if commit.Message != "" {
			commitMessages = []string{commit.Message}
		}
	}

	if diff == "" {
		return nil, fmt.Errorf("%w: empty diff", domain.ErrValidation)
	}

	diffHash := cache.DiffCacheKey(diff + ":" + level)

	if !hasOverrides(req.Overrides) {
		if cached, ok := s.cache.Get(diffHash); ok {
			slog.Debug("service: cache hit", "hash", diffHash[:12])
			analysis := &domain.Analysis{
				AnalysisType:  domain.AnalysisTypeSingleCommit,
				Status:        domain.AnalysisStatusCompleted,
				Workspace:     req.Workspace,
				RepoSlug:      req.RepoSlug,
				CommitHash:    req.CommitHash,
				Level:         level,
				DiffHash:      diffHash,
				GeneratedDesc: cached,
				ModelUsed:     "cache",
			}
			if err := s.repository.Create(ctx, analysis); err != nil {
				slog.Warn("failed to save cached analysis", "error", err)
			}
			return analysis, nil
		}

		existing, err := s.repository.GetByDiffHash(ctx, diffHash)
		if err == nil {
			slog.Debug("service: db hit", "hash", diffHash[:12])
			return existing, nil
		}
	}

	genInput := openai.GenerationInput{
		Diff:           diff,
		AnalysisType:   string(domain.AnalysisTypeSingleCommit),
		CommitMessages: commitMessages,
		Level:          level,
	}
	applyOverrides(&genInput, req.Overrides)

	output, err := s.generator.Generate(ctx, genInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrExternalService, err.Error())
	}

	tokensUsed := output.TokensUsed
	analysis := &domain.Analysis{
		AnalysisType:  domain.AnalysisTypeSingleCommit,
		Status:        domain.AnalysisStatusCompleted,
		Workspace:     req.Workspace,
		RepoSlug:      req.RepoSlug,
		CommitHash:    req.CommitHash,
		Level:         level,
		RawDiff:       diff,
		DiffHash:      diffHash,
		GeneratedDesc: output.Description,
		ModelUsed:     output.Model,
		TokensUsed:    &tokensUsed,
	}

	if err := s.repository.Create(ctx, analysis); err != nil {
		slog.Error("failed to save analysis", "error", err)
	}

	return analysis, nil
}

func (s *analysisService) AnalyzeRange(ctx context.Context, req *dto.AnalyzeRangeRequest) (*domain.Analysis, error) {
	level := defaultLevel(req.Level)

	commits, err := s.bbClient.ListCommitsInRange(ctx, req.Workspace, req.RepoSlug, req.ToHash, req.FromHash)
	if err != nil {
		return nil, s.mapExternalError(err)
	}

	var commitMessages []string
	for _, c := range commits {
		if c.Message != "" {
			commitMessages = append(commitMessages, c.Message)
		}
	}

	diff, err := s.bbClient.GetCommitDiff(ctx, req.Workspace, req.RepoSlug, req.FromHash+".."+req.ToHash)
	if err != nil {
		// Fallback: concatenate individual commit diffs
		var diffs []string
		for _, c := range commits {
			d, dErr := s.bbClient.GetCommitDiff(ctx, req.Workspace, req.RepoSlug, c.Hash)
			if dErr != nil {
				continue
			}
			diffs = append(diffs, d)
		}
		diff = strings.Join(diffs, "\n")
	}

	if diff == "" {
		return nil, fmt.Errorf("%w: empty diff for range", domain.ErrValidation)
	}

	diffHash := cache.DiffCacheKey(diff + ":" + level)

	if !hasOverrides(req.Overrides) {
		if cached, ok := s.cache.Get(diffHash); ok {
			analysis := &domain.Analysis{
				AnalysisType:  domain.AnalysisTypeCommitRange,
				Status:        domain.AnalysisStatusCompleted,
				Workspace:     req.Workspace,
				RepoSlug:      req.RepoSlug,
				FromHash:      req.FromHash,
				ToHash:        req.ToHash,
				Level:         level,
				DiffHash:      diffHash,
				GeneratedDesc: cached,
				ModelUsed:     "cache",
			}
			if err := s.repository.Create(ctx, analysis); err != nil {
				slog.Warn("failed to save cached analysis", "error", err)
			}
			return analysis, nil
		}

		existing, err := s.repository.GetByDiffHash(ctx, diffHash)
		if err == nil {
			return existing, nil
		}
	}

	genInput := openai.GenerationInput{
		Diff:           diff,
		AnalysisType:   string(domain.AnalysisTypeCommitRange),
		CommitMessages: commitMessages,
		Level:          level,
	}
	applyOverrides(&genInput, req.Overrides)

	output, err := s.generator.Generate(ctx, genInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrExternalService, err.Error())
	}

	tokensUsed := output.TokensUsed
	analysis := &domain.Analysis{
		AnalysisType:  domain.AnalysisTypeCommitRange,
		Status:        domain.AnalysisStatusCompleted,
		Workspace:     req.Workspace,
		RepoSlug:      req.RepoSlug,
		FromHash:      req.FromHash,
		ToHash:        req.ToHash,
		Level:         level,
		RawDiff:       diff,
		DiffHash:      diffHash,
		GeneratedDesc: output.Description,
		ModelUsed:     output.Model,
		TokensUsed:    &tokensUsed,
	}

	if err := s.repository.Create(ctx, analysis); err != nil {
		slog.Error("failed to save analysis", "error", err)
	}

	return analysis, nil
}

func (s *analysisService) AnalyzePR(ctx context.Context, req *dto.AnalyzePRRequest) (*domain.Analysis, error) {
	diff := req.RawDiff
	prTitle := req.PRTitle
	prDesc := req.PRDescription
	level := defaultLevel(req.Level)

	if diff == "" {
		var err error
		diff, err = s.bbClient.GetPullRequestDiff(ctx, req.Workspace, req.RepoSlug, req.PRID)
		if err != nil {
			return nil, s.mapExternalError(err)
		}
	}

	if (prTitle == "" || prDesc == "") && req.PRID > 0 {
		pr, err := s.bbClient.GetPullRequest(ctx, req.Workspace, req.RepoSlug, req.PRID)
		if err != nil {
			slog.Warn("failed to get PR details, proceeding without them", "error", err)
		} else {
			if prTitle == "" {
				prTitle = pr.Title
			}
			if prDesc == "" {
				prDesc = pr.Description
			}
		}
	}

	if diff == "" {
		return nil, fmt.Errorf("%w: empty diff", domain.ErrValidation)
	}

	diffHash := cache.DiffCacheKey(diff + ":" + level)

	if !hasOverrides(req.Overrides) {
		if cached, ok := s.cache.Get(diffHash); ok {
			prID := req.PRID
			analysis := &domain.Analysis{
				AnalysisType:  domain.AnalysisTypePullRequest,
				Status:        domain.AnalysisStatusCompleted,
				Workspace:     req.Workspace,
				RepoSlug:      req.RepoSlug,
				PrID:          &prID,
				Level:         level,
				DiffHash:      diffHash,
				GeneratedDesc: cached,
				ModelUsed:     "cache",
			}
			if err := s.repository.Create(ctx, analysis); err != nil {
				slog.Warn("failed to save cached analysis", "error", err)
			}
			return analysis, nil
		}

		existing, err := s.repository.GetByDiffHash(ctx, diffHash)
		if err == nil {
			return existing, nil
		}
	}

	genInput := openai.GenerationInput{
		Diff:          diff,
		AnalysisType:  string(domain.AnalysisTypePullRequest),
		PRTitle:       prTitle,
		PRDescription: prDesc,
		Level:         level,
	}
	applyOverrides(&genInput, req.Overrides)

	output, err := s.generator.Generate(ctx, genInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrExternalService, err.Error())
	}

	tokensUsed := output.TokensUsed
	prID := req.PRID
	analysis := &domain.Analysis{
		AnalysisType:  domain.AnalysisTypePullRequest,
		Status:        domain.AnalysisStatusCompleted,
		Workspace:     req.Workspace,
		RepoSlug:      req.RepoSlug,
		PrID:          &prID,
		Level:         level,
		RawDiff:       diff,
		DiffHash:      diffHash,
		GeneratedDesc: output.Description,
		ModelUsed:     output.Model,
		TokensUsed:    &tokensUsed,
	}

	if err := s.repository.Create(ctx, analysis); err != nil {
		slog.Error("failed to save analysis", "error", err)
	}

	return analysis, nil
}

func (s *analysisService) GetAnalysis(ctx context.Context, id string) (*domain.Analysis, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *analysisService) mapExternalError(err error) error {
	var notFound *bitbucket.NotFoundError
	if errors.As(err, &notFound) {
		return fmt.Errorf("%w: %s", domain.ErrNotFound, err.Error())
	}

	var unauthorized *bitbucket.UnauthorizedError
	if errors.As(err, &unauthorized) {
		return fmt.Errorf("%w: %s", domain.ErrExternalService, err.Error())
	}

	return fmt.Errorf("%w: %s", domain.ErrExternalService, err.Error())
}
