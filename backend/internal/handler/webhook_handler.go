package handler

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/igor-trentini/diffable/backend/internal/handler/dto"
	"github.com/igor-trentini/diffable/backend/internal/repository"
	"github.com/igor-trentini/diffable/backend/internal/service"
	"github.com/igor-trentini/diffable/backend/internal/webhook"
)

type WebhookHandler struct {
	analysisService   service.AnalysisService
	webhookRepository repository.WebhookRepository
}

func NewWebhookHandler(analysisSvc service.AnalysisService, webhookRepo repository.WebhookRepository) *WebhookHandler {
	return &WebhookHandler{
		analysisService:   analysisSvc,
		webhookRepository: webhookRepo,
	}
}

func (h *WebhookHandler) HandleBitbucket(w http.ResponseWriter, r *http.Request) {
	eventKey := r.Header.Get("X-Event-Key")
	if eventKey != "pullrequest:created" {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_event", "Only pullrequest:created events are supported")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_request", "Failed to read request body")
		return
	}

	payload, err := webhook.ParsePayload(body)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	// Log the webhook
	webhookLog := &repository.WebhookLog{
		EventKey: eventKey,
		Payload:  body,
		Status:   "received",
	}
	if err := h.webhookRepository.Create(r.Context(), webhookLog); err != nil {
		slog.Error("failed to log webhook", "error", err)
	}

	// Process in background
	params := webhook.ExtractAnalysisParams(payload)
	go h.processWebhook(webhookLog.ID, params)

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":     "accepted",
		"webhook_id": webhookLog.ID,
	})
}

func (h *WebhookHandler) processWebhook(webhookID string, params webhook.AnalysisParams) {
	ctx := context.Background()

	analysis, err := h.analysisService.AnalyzePR(ctx, &dto.AnalyzePRRequest{
		Workspace: params.Workspace,
		RepoSlug:  params.RepoSlug,
		PRID:      params.PRID,
	})

	if err != nil {
		slog.Error("webhook: failed to analyze PR",
			"webhook_id", webhookID,
			"workspace", params.Workspace,
			"repo", params.RepoSlug,
			"pr_id", params.PRID,
			"error", err,
		)
		if updateErr := h.webhookRepository.UpdateStatus(ctx, webhookID, "failed", nil, err.Error()); updateErr != nil {
			slog.Error("webhook: failed to update status", "error", updateErr)
		}
		return
	}

	slog.Info("webhook: analysis complete",
		"webhook_id", webhookID,
		"analysis_id", analysis.ID,
	)
	if updateErr := h.webhookRepository.UpdateStatus(ctx, webhookID, "completed", &analysis.ID, ""); updateErr != nil {
		slog.Error("webhook: failed to update status", "error", updateErr)
	}
}
