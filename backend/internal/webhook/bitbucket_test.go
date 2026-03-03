package webhook

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePayload_Valid(t *testing.T) {
	payload := WebhookPayload{
		PullRequest: WebhookPullRequest{
			ID:    42,
			Title: "Test PR",
		},
		Repository: WebhookRepository{
			Slug:     "my-repo",
			FullName: "ws/my-repo",
			Owner:    WebhookOwner{Username: "ws"},
		},
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	result, err := ParsePayload(body)

	require.NoError(t, err)
	assert.Equal(t, 42, result.PullRequest.ID)
	assert.Equal(t, "my-repo", result.Repository.Slug)
}

func TestParsePayload_InvalidJSON(t *testing.T) {
	_, err := ParsePayload([]byte("not json"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid webhook payload")
}

func TestParsePayload_MissingPRID(t *testing.T) {
	payload := WebhookPayload{
		Repository: WebhookRepository{
			Slug:     "my-repo",
			FullName: "ws/my-repo",
		},
	}

	body, _ := json.Marshal(payload)
	_, err := ParsePayload(body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing pull request ID")
}

func TestParsePayload_MissingRepository(t *testing.T) {
	payload := WebhookPayload{
		PullRequest: WebhookPullRequest{ID: 1},
	}

	body, _ := json.Marshal(payload)
	_, err := ParsePayload(body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing repository information")
}

func TestExtractAnalysisParams(t *testing.T) {
	payload := &WebhookPayload{
		PullRequest: WebhookPullRequest{ID: 42},
		Repository: WebhookRepository{
			Slug:  "my-repo",
			Owner: WebhookOwner{Username: "my-workspace"},
		},
	}

	params := ExtractAnalysisParams(payload)

	assert.Equal(t, "my-workspace", params.Workspace)
	assert.Equal(t, "my-repo", params.RepoSlug)
	assert.Equal(t, 42, params.PRID)
}

func TestExtractAnalysisParams_FallsBackToProject(t *testing.T) {
	payload := &WebhookPayload{
		PullRequest: WebhookPullRequest{ID: 1},
		Repository: WebhookRepository{
			Slug:    "repo",
			Project: WebhookProject{Key: "PROJ"},
		},
	}

	params := ExtractAnalysisParams(payload)

	assert.Equal(t, "PROJ", params.Workspace)
}
