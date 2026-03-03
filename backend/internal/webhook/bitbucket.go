package webhook

import (
	"encoding/json"
	"fmt"
)

type WebhookPayload struct {
	PullRequest WebhookPullRequest `json:"pullrequest"`
	Repository  WebhookRepository  `json:"repository"`
	Actor       WebhookActor       `json:"actor"`
}

type WebhookPullRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
}

type WebhookRepository struct {
	Slug     string          `json:"slug"`
	Name     string          `json:"name"`
	FullName string          `json:"full_name"`
	Project  WebhookProject  `json:"project"`
	Owner    WebhookOwner    `json:"owner"`
}

type WebhookProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type WebhookOwner struct {
	Username string `json:"username"`
}

type WebhookActor struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
}

type AnalysisParams struct {
	Workspace string
	RepoSlug  string
	PRID      int
}

func ParsePayload(body []byte) (*WebhookPayload, error) {
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("invalid webhook payload: %w", err)
	}

	if payload.PullRequest.ID == 0 {
		return nil, fmt.Errorf("missing pull request ID in webhook payload")
	}

	if payload.Repository.FullName == "" && payload.Repository.Slug == "" {
		return nil, fmt.Errorf("missing repository information in webhook payload")
	}

	return &payload, nil
}

func ExtractAnalysisParams(payload *WebhookPayload) AnalysisParams {
	workspace := payload.Repository.Owner.Username
	if workspace == "" {
		workspace = payload.Repository.Project.Key
	}

	return AnalysisParams{
		Workspace: workspace,
		RepoSlug:  payload.Repository.Slug,
		PRID:      payload.PullRequest.ID,
	}
}
