package dto

import "github.com/igor-trentini/diffable/backend/internal/domain"

type AnalysisResponse struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Model       string `json:"model_used"`
	TokensUsed  int    `json:"tokens_used"`
	CreatedAt   string `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func AnalysisToResponse(a *domain.Analysis) *AnalysisResponse {
	tokensUsed := 0
	if a.TokensUsed != nil {
		tokensUsed = *a.TokensUsed
	}
	return &AnalysisResponse{
		ID:          a.ID,
		Type:        string(a.AnalysisType),
		Description: a.GeneratedDesc,
		Model:       a.ModelUsed,
		TokensUsed:  tokensUsed,
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
