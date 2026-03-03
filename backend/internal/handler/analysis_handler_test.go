package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/igor-trentini/diffable/backend/internal/handler/dto"
)

type mockAnalysisService struct {
	analysis *domain.Analysis
	err      error
}

func (m *mockAnalysisService) AnalyzeCommit(_ context.Context, _ *dto.AnalyzeCommitRequest) (*domain.Analysis, error) {
	return m.analysis, m.err
}

func (m *mockAnalysisService) AnalyzeRange(_ context.Context, _ *dto.AnalyzeRangeRequest) (*domain.Analysis, error) {
	return m.analysis, m.err
}

func (m *mockAnalysisService) AnalyzePR(_ context.Context, _ *dto.AnalyzePRRequest) (*domain.Analysis, error) {
	return m.analysis, m.err
}

func (m *mockAnalysisService) GetAnalysis(_ context.Context, _ string) (*domain.Analysis, error) {
	return m.analysis, m.err
}

func newTestAnalysis() *domain.Analysis {
	tokens := 100
	return &domain.Analysis{
		ID:            "test-id-123",
		AnalysisType:  domain.AnalysisTypeSingleCommit,
		Status:        domain.AnalysisStatusCompleted,
		GeneratedDesc: "Test description",
		ModelUsed:     "gpt-4o-mini",
		TokensUsed:    &tokens,
		CreatedAt:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestAnalyzeCommit_ValidJSON_Returns200(t *testing.T) {
	svc := &mockAnalysisService{analysis: newTestAnalysis()}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{RawDiff: "diff --git a/main.go"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.AnalysisResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "test-id-123", resp.ID)
	assert.Equal(t, "Test description", resp.Description)
}

func TestAnalyzeCommit_MissingFields_Returns400(t *testing.T) {
	svc := &mockAnalysisService{}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyzeCommit_InvalidJSON_Returns400(t *testing.T) {
	svc := &mockAnalysisService{}
	h := NewAnalysisHandler(svc, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyzeCommit_ServiceError_Returns500(t *testing.T) {
	svc := &mockAnalysisService{err: assert.AnError}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{RawDiff: "diff --git a/main.go"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAnalyzeRange_ValidJSON_Returns200(t *testing.T) {
	svc := &mockAnalysisService{analysis: newTestAnalysis()}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeRangeRequest{
		Workspace: "ws",
		RepoSlug:  "repo",
		FromHash:  "abc123",
		ToHash:    "def456",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/range", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeRange(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyzeRange_MissingHashes_Returns400(t *testing.T) {
	svc := &mockAnalysisService{}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeRangeRequest{Workspace: "ws"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/range", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeRange(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyzePR_ValidJSON_Returns200(t *testing.T) {
	svc := &mockAnalysisService{analysis: newTestAnalysis()}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzePRRequest{
		Workspace: "ws",
		RepoSlug:  "repo",
		PRID:      42,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/pr", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzePR(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyzePR_MissingFields_Returns400(t *testing.T) {
	svc := &mockAnalysisService{}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzePRRequest{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/pr", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzePR(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyzePR_NotFound_Returns404(t *testing.T) {
	svc := &mockAnalysisService{err: domain.ErrNotFound}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzePRRequest{
		Workspace: "ws",
		RepoSlug:  "repo",
		PRID:      999,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/pr", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzePR(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAnalysis_ExistingID_Returns200(t *testing.T) {
	svc := &mockAnalysisService{analysis: newTestAnalysis()}
	h := NewAnalysisHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analyses/test-id-123", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-id-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetAnalysis(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.AnalysisResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "test-id-123", resp.ID)
}

func TestGetAnalysis_NonExistingID_Returns404(t *testing.T) {
	svc := &mockAnalysisService{err: domain.ErrNotFound}
	h := NewAnalysisHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analyses/non-existent", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "non-existent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetAnalysis(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAnalyzeCommit_ExternalServiceError_Returns502(t *testing.T) {
	svc := &mockAnalysisService{err: domain.ErrExternalService}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{RawDiff: "diff --git a/main.go"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
}

func TestHandleServiceError_RateLimited_Returns429(t *testing.T) {
	svc := &mockAnalysisService{err: domain.ErrRateLimited}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{RawDiff: "diff --git a/main.go"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Equal(t, "60", w.Header().Get("Retry-After"))
}

func TestHandleServiceError_Timeout_Returns504(t *testing.T) {
	svc := &mockAnalysisService{err: domain.ErrTimeout}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{RawDiff: "diff --git a/main.go"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
}

func TestHandleServiceError_TokenLimitExceeded_Returns422(t *testing.T) {
	svc := &mockAnalysisService{err: domain.ErrTokenLimitExceeded}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{RawDiff: "diff --git a/main.go"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestHandleServiceError_ValidationError_Returns400(t *testing.T) {
	svc := &mockAnalysisService{err: domain.ErrValidation}
	h := NewAnalysisHandler(svc, nil)

	body, _ := json.Marshal(dto.AnalyzeCommitRequest{RawDiff: "diff --git a/main.go"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyses/commit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AnalyzeCommit(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
