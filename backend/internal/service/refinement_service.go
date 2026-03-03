package service

import (
	"context"
	"fmt"

	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/igor-trentini/diffable/backend/internal/openai"
	"github.com/igor-trentini/diffable/backend/internal/repository"
)

type RefinementService interface {
	Refine(ctx context.Context, analysisID string, instruction string) (*domain.Refinement, error)
}

type refinementService struct {
	generator  openai.DescriptionGenerator
	repository repository.AnalysisRepository
}

func NewRefinementService(generator openai.DescriptionGenerator, repo repository.AnalysisRepository) RefinementService {
	return &refinementService{
		generator:  generator,
		repository: repo,
	}
}

func (s *refinementService) Refine(ctx context.Context, analysisID string, instruction string) (*domain.Refinement, error) {
	if instruction == "" {
		return nil, fmt.Errorf("%w: instruction is required", domain.ErrValidation)
	}

	analysis, err := s.repository.GetByID(ctx, analysisID)
	if err != nil {
		return nil, err
	}

	output, err := s.generator.Refine(ctx, openai.RefinementInput{
		OriginalDescription: analysis.GeneratedDesc,
		Instruction:         instruction,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrExternalService, err.Error())
	}

	tokensUsed := output.TokensUsed
	refinement := &domain.Refinement{
		AnalysisID:  analysisID,
		Instruction: instruction,
		OriginalDesc: analysis.GeneratedDesc,
		RefinedDesc:  output.Description,
		ModelUsed:    output.Model,
		TokensUsed:   &tokensUsed,
	}

	if err := s.repository.CreateRefinement(ctx, refinement); err != nil {
		return nil, fmt.Errorf("saving refinement: %w", err)
	}

	return refinement, nil
}
