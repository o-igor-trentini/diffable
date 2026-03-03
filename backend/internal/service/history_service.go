package service

import (
	"context"

	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/igor-trentini/diffable/backend/internal/repository"
)

type HistoryService interface {
	ListAnalyses(ctx context.Context, typeFilter string, page, pageSize int) ([]domain.Analysis, int, error)
	GetRefinements(ctx context.Context, analysisID string) ([]domain.Refinement, error)
}

type historyService struct {
	repository repository.AnalysisRepository
}

func NewHistoryService(repo repository.AnalysisRepository) HistoryService {
	return &historyService{repository: repo}
}

func (s *historyService) ListAnalyses(ctx context.Context, typeFilter string, page, pageSize int) ([]domain.Analysis, int, error) {
	offset := (page - 1) * pageSize
	filter := repository.AnalysisFilter{Type: typeFilter}
	return s.repository.List(ctx, filter, offset, pageSize)
}

func (s *historyService) GetRefinements(ctx context.Context, analysisID string) ([]domain.Refinement, error) {
	return s.repository.ListRefinements(ctx, analysisID)
}
