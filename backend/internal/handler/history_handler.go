package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/igor-trentini/diffable/backend/internal/handler/dto"
	"github.com/igor-trentini/diffable/backend/internal/service"
)

type HistoryHandler struct {
	historyService service.HistoryService
}

func NewHistoryHandler(svc service.HistoryService) *HistoryHandler {
	return &HistoryHandler{historyService: svc}
}

func (h *HistoryHandler) ListAnalyses(w http.ResponseWriter, r *http.Request) {
	typeFilter := r.URL.Query().Get("type")
	page := parseIntParam(r.URL.Query().Get("page"), 1)
	pageSize := parseIntParam(r.URL.Query().Get("page_size"), 20)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	analyses, total, err := h.historyService.ListAnalyses(r.Context(), typeFilter, page, pageSize)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
		return
	}

	data := make([]dto.AnalysisResponse, 0, len(analyses))
	for i := range analyses {
		data = append(data, *dto.AnalysisToResponse(&analyses[i]))
	}

	writeJSON(w, http.StatusOK, dto.PaginatedAnalysesResponse{
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *HistoryHandler) GetRefinements(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", "id is required")
		return
	}

	refinements, err := h.historyService.GetRefinements(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeErrorJSON(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		writeErrorJSON(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
		return
	}

	data := make([]dto.RefinementResponse, 0, len(refinements))
	for i := range refinements {
		data = append(data, *dto.RefinementToResponse(&refinements[i]))
	}

	writeJSON(w, http.StatusOK, data)
}

func parseIntParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}
