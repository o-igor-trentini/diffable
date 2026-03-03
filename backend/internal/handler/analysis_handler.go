package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/igor-trentini/diffable/backend/internal/handler/dto"
	"github.com/igor-trentini/diffable/backend/internal/service"
)

type AnalysisHandler struct {
	analysisService    service.AnalysisService
	refinementService  service.RefinementService
}

func NewAnalysisHandler(svc service.AnalysisService, refineSvc service.RefinementService) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService:   svc,
		refinementService: refineSvc,
	}
}

func (h *AnalysisHandler) AnalyzeCommit(w http.ResponseWriter, r *http.Request) {
	var req dto.AnalyzeCommitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	if err := req.Validate(); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	analysis, err := h.analysisService.AnalyzeCommit(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.AnalysisToResponse(analysis))
}

func (h *AnalysisHandler) AnalyzeRange(w http.ResponseWriter, r *http.Request) {
	var req dto.AnalyzeRangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	if err := req.Validate(); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	analysis, err := h.analysisService.AnalyzeRange(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.AnalysisToResponse(analysis))
}

func (h *AnalysisHandler) AnalyzePR(w http.ResponseWriter, r *http.Request) {
	var req dto.AnalyzePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	if err := req.Validate(); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	analysis, err := h.analysisService.AnalyzePR(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.AnalysisToResponse(analysis))
}

func (h *AnalysisHandler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", "id is required")
		return
	}

	analysis, err := h.analysisService.GetAnalysis(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.AnalysisToResponse(analysis))
}

func (h *AnalysisHandler) RefineDescription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", "id is required")
		return
	}

	var req dto.RefineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	if err := req.Validate(); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	refinement, err := h.refinementService.Refine(r.Context(), id, req.Instruction)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.RefinementToResponse(refinement))
}

func (h *AnalysisHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeErrorJSON(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, domain.ErrValidation):
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", err.Error())
	case errors.Is(err, domain.ErrExternalService):
		writeErrorJSON(w, http.StatusBadGateway, "external_service_error", err.Error())
	default:
		writeErrorJSON(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
	}
}

func writeErrorJSON(w http.ResponseWriter, status int, errType, message string) {
	writeJSON(w, status, dto.ErrorResponse{
		Error:   errType,
		Message: message,
	})
}
