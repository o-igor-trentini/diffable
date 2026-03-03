package dto

import (
	"testing"

	"github.com/igor-trentini/diffable/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func intPtr(v int) *int          { return &v }
func float64Ptr(v float64) *float64 { return &v }
func strPtr(v string) *string    { return &v }

func TestGenerationOverrides_Validate(t *testing.T) {
	tests := []struct {
		name      string
		overrides *GenerationOverrides
		wantErr   bool
	}{
		{
			name:      "nil overrides passes",
			overrides: nil,
			wantErr:   false,
		},
		{
			name:      "empty struct passes",
			overrides: &GenerationOverrides{},
			wantErr:   false,
		},
		{
			name:      "valid max_tokens",
			overrides: &GenerationOverrides{MaxTokens: intPtr(1024)},
			wantErr:   false,
		},
		{
			name:      "valid temperature",
			overrides: &GenerationOverrides{Temperature: float64Ptr(0.5)},
			wantErr:   false,
		},
		{
			name:      "valid model auto",
			overrides: &GenerationOverrides{Model: strPtr("auto")},
			wantErr:   false,
		},
		{
			name:      "valid model gpt-4o-mini",
			overrides: &GenerationOverrides{Model: strPtr("gpt-4o-mini")},
			wantErr:   false,
		},
		{
			name:      "valid model gpt-4o",
			overrides: &GenerationOverrides{Model: strPtr("gpt-4o")},
			wantErr:   false,
		},
		{
			name:      "all valid fields",
			overrides: &GenerationOverrides{MaxTokens: intPtr(2048), Temperature: float64Ptr(0.7), Model: strPtr("gpt-4o")},
			wantErr:   false,
		},
		{
			name:      "max_tokens boundary min valid",
			overrides: &GenerationOverrides{MaxTokens: intPtr(64)},
			wantErr:   false,
		},
		{
			name:      "max_tokens boundary max valid",
			overrides: &GenerationOverrides{MaxTokens: intPtr(4096)},
			wantErr:   false,
		},
		{
			name:      "temperature boundary min valid",
			overrides: &GenerationOverrides{Temperature: float64Ptr(0.0)},
			wantErr:   false,
		},
		{
			name:      "temperature boundary max valid",
			overrides: &GenerationOverrides{Temperature: float64Ptr(2.0)},
			wantErr:   false,
		},
		// Invalid cases
		{
			name:      "max_tokens too low",
			overrides: &GenerationOverrides{MaxTokens: intPtr(63)},
			wantErr:   true,
		},
		{
			name:      "max_tokens too high",
			overrides: &GenerationOverrides{MaxTokens: intPtr(4097)},
			wantErr:   true,
		},
		{
			name:      "max_tokens zero",
			overrides: &GenerationOverrides{MaxTokens: intPtr(0)},
			wantErr:   true,
		},
		{
			name:      "max_tokens negative",
			overrides: &GenerationOverrides{MaxTokens: intPtr(-1)},
			wantErr:   true,
		},
		{
			name:      "temperature negative",
			overrides: &GenerationOverrides{Temperature: float64Ptr(-0.1)},
			wantErr:   true,
		},
		{
			name:      "temperature too high",
			overrides: &GenerationOverrides{Temperature: float64Ptr(2.1)},
			wantErr:   true,
		},
		{
			name:      "invalid model",
			overrides: &GenerationOverrides{Model: strPtr("gpt-5")},
			wantErr:   true,
		},
		{
			name:      "empty model string is invalid",
			overrides: &GenerationOverrides{Model: strPtr("")},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.overrides.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, domain.ErrValidation)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAnalyzeCommitRequest_Validate_WithOverrides(t *testing.T) {
	t.Run("valid request with valid overrides", func(t *testing.T) {
		req := &AnalyzeCommitRequest{
			RawDiff:   "diff --git a/main.go",
			Overrides: &GenerationOverrides{MaxTokens: intPtr(2048)},
		}
		assert.NoError(t, req.Validate())
	})

	t.Run("valid request with invalid overrides", func(t *testing.T) {
		req := &AnalyzeCommitRequest{
			RawDiff:   "diff --git a/main.go",
			Overrides: &GenerationOverrides{MaxTokens: intPtr(0)},
		}
		err := req.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrValidation)
	})

	t.Run("valid request without overrides", func(t *testing.T) {
		req := &AnalyzeCommitRequest{
			RawDiff: "diff --git a/main.go",
		}
		assert.NoError(t, req.Validate())
	})
}

func TestAnalyzeRangeRequest_Validate_WithOverrides(t *testing.T) {
	t.Run("valid request with invalid overrides", func(t *testing.T) {
		req := &AnalyzeRangeRequest{
			Workspace: "ws",
			RepoSlug:  "repo",
			FromHash:  "abc",
			ToHash:    "def",
			Overrides: &GenerationOverrides{Model: strPtr("invalid")},
		}
		err := req.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}

func TestAnalyzePRRequest_Validate_WithOverrides(t *testing.T) {
	t.Run("valid request with invalid overrides", func(t *testing.T) {
		req := &AnalyzePRRequest{
			RawDiff:   "diff --git a/main.go",
			PRTitle:   "title",
			Overrides: &GenerationOverrides{Temperature: float64Ptr(3.0)},
		}
		err := req.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}
