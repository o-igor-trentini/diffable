package dto

import "github.com/igor-trentini/diffable/backend/internal/domain"

type AnalysisResponse struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Level       string `json:"level"`
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

type RefinementResponse struct {
	ID          string `json:"id"`
	AnalysisID  string `json:"analysis_id"`
	Instruction string `json:"instruction"`
	RefinedDesc string `json:"refined_description"`
	Model       string `json:"model_used"`
	TokensUsed  int    `json:"tokens_used"`
	CreatedAt   string `json:"created_at"`
}

type PaginatedAnalysesResponse struct {
	Data     []AnalysisResponse `json:"data"`
	Total    int                `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

func AnalysisToResponse(a *domain.Analysis) *AnalysisResponse {
	tokensUsed := 0
	if a.TokensUsed != nil {
		tokensUsed = *a.TokensUsed
	}
	level := a.Level
	if level == "" {
		level = "functional"
	}
	return &AnalysisResponse{
		ID:          a.ID,
		Type:        string(a.AnalysisType),
		Level:       level,
		Description: a.GeneratedDesc,
		Model:       a.ModelUsed,
		TokensUsed:  tokensUsed,
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func RefinementToResponse(r *domain.Refinement) *RefinementResponse {
	tokensUsed := 0
	if r.TokensUsed != nil {
		tokensUsed = *r.TokensUsed
	}
	return &RefinementResponse{
		ID:          r.ID,
		AnalysisID:  r.AnalysisID,
		Instruction: r.Instruction,
		RefinedDesc: r.RefinedDesc,
		Model:       r.ModelUsed,
		TokensUsed:  tokensUsed,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
