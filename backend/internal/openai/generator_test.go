package openai

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/igor-trentini/diffable/backend/internal/cache"
	oai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockChatClient struct {
	calls    int
	response oai.ChatCompletionResponse
	err      error
	errOnce  bool // if true, only return error on first call
}

func (m *mockChatClient) CreateChatCompletion(_ context.Context, _ oai.ChatCompletionRequest) (oai.ChatCompletionResponse, error) {
	m.calls++
	if m.err != nil && (m.calls == 1 || !m.errOnce) {
		return oai.ChatCompletionResponse{}, m.err
	}
	return m.response, nil
}

var testGenConfig = GeneratorConfig{
	DefaultModel:   "gpt-4o-mini",
	ComplexModel:   "gpt-4o",
	MaxTokens:      1024,
	Temperature:    0.3,
	TokenThreshold: 4000,
	CacheTTL:       1 * time.Hour,
}

func TestGenerate_Success(t *testing.T) {
	mock := &mockChatClient{
		response: oai.ChatCompletionResponse{
			Choices: []oai.ChatCompletionChoice{
				{Message: oai.ChatCompletionMessage{Content: "**Resumo:** Mudança realizada"}},
			},
			Usage: oai.Usage{TotalTokens: 150},
		},
	}

	gen := NewGenerator(mock, cache.NewInMemoryCache(), testGenConfig)
	out, err := gen.Generate(context.Background(), GenerationInput{
		Diff:         "diff --git a/main.go b/main.go\n+func main() {}",
		AnalysisType: "single_commit",
	})

	require.NoError(t, err)
	assert.Equal(t, "**Resumo:** Mudança realizada", out.Description)
	assert.Equal(t, "gpt-4o-mini", out.Model)
	assert.Equal(t, 150, out.TokensUsed)
	assert.Equal(t, 1, mock.calls)
}

func TestGenerate_CacheHit(t *testing.T) {
	mock := &mockChatClient{
		response: oai.ChatCompletionResponse{
			Choices: []oai.ChatCompletionChoice{
				{Message: oai.ChatCompletionMessage{Content: "Cached description"}},
			},
			Usage: oai.Usage{TotalTokens: 100},
		},
	}

	c := cache.NewInMemoryCache()
	gen := NewGenerator(mock, c, testGenConfig)

	diff := "diff --git a/main.go b/main.go\n+func hello() {}"
	input := GenerationInput{Diff: diff, AnalysisType: "single_commit"}

	// First call — populates cache
	out1, err := gen.Generate(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, 1, mock.calls)

	// Second call — cache hit, no API call
	out2, err := gen.Generate(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, 1, mock.calls) // Still 1 — not called again
	assert.Equal(t, out1.Description, out2.Description)
	assert.Equal(t, "cache", out2.Model)
}

func TestGenerate_RetryOn429(t *testing.T) {
	mock := &mockChatClient{
		response: oai.ChatCompletionResponse{
			Choices: []oai.ChatCompletionChoice{
				{Message: oai.ChatCompletionMessage{Content: "Success after retry"}},
			},
			Usage: oai.Usage{TotalTokens: 120},
		},
		err:     &oai.APIError{HTTPStatusCode: http.StatusTooManyRequests, Message: "rate limited"},
		errOnce: true,
	}

	gen := NewGenerator(mock, cache.NewInMemoryCache(), testGenConfig)
	out, err := gen.Generate(context.Background(), GenerationInput{
		Diff:         "diff --git a/x.go b/x.go\n+change",
		AnalysisType: "single_commit",
	})

	require.NoError(t, err)
	assert.Equal(t, "Success after retry", out.Description)
	assert.Equal(t, 2, mock.calls)
}

func TestGenerate_PRUsesComplexModel(t *testing.T) {
	var capturedReq oai.ChatCompletionRequest
	mock := &mockChatClient{
		response: oai.ChatCompletionResponse{
			Choices: []oai.ChatCompletionChoice{
				{Message: oai.ChatCompletionMessage{Content: "PR description"}},
			},
			Usage: oai.Usage{TotalTokens: 200},
		},
	}

	// Wrap to capture request
	wrapper := &requestCapturingClient{inner: mock, captured: &capturedReq}
	gen := NewGenerator(wrapper, cache.NewInMemoryCache(), testGenConfig)
	out, err := gen.Generate(context.Background(), GenerationInput{
		Diff:         "diff --git a/main.go b/main.go\n+func main() {}",
		AnalysisType: "pull_request",
		PRTitle:      "Add feature",
	})

	require.NoError(t, err)
	assert.Equal(t, "gpt-4o", out.Model)
}

func TestRefine_Success(t *testing.T) {
	mock := &mockChatClient{
		response: oai.ChatCompletionResponse{
			Choices: []oai.ChatCompletionChoice{
				{Message: oai.ChatCompletionMessage{Content: "Refined description"}},
			},
			Usage: oai.Usage{TotalTokens: 180},
		},
	}

	gen := NewGenerator(mock, cache.NewInMemoryCache(), testGenConfig)
	out, err := gen.Refine(context.Background(), RefinementInput{
		OriginalDescription: "Original text",
		Instruction:         "Torne mais conciso",
	})

	require.NoError(t, err)
	assert.Equal(t, "Refined description", out.Description)
	assert.Equal(t, "gpt-4o-mini", out.Model)
	assert.Equal(t, 180, out.TokensUsed)
}

func TestGenerate_ContextCancelled(t *testing.T) {
	mock := &mockChatClient{
		err: &oai.APIError{HTTPStatusCode: http.StatusTooManyRequests},
	}

	gen := NewGenerator(mock, cache.NewInMemoryCache(), testGenConfig)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := gen.Generate(ctx, GenerationInput{
		Diff:         "diff --git a/x.go b/x.go\n+x",
		AnalysisType: "single_commit",
	})
	require.Error(t, err)
}

type requestCapturingClient struct {
	inner    *mockChatClient
	captured *oai.ChatCompletionRequest
}

func (r *requestCapturingClient) CreateChatCompletion(ctx context.Context, req oai.ChatCompletionRequest) (oai.ChatCompletionResponse, error) {
	*r.captured = req
	return r.inner.CreateChatCompletion(ctx, req)
}
