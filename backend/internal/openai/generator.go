package openai

import (
	"context"
	"log/slog"
	"time"

	"github.com/igor-trentini/diffable/backend/internal/cache"
	oai "github.com/sashabaranov/go-openai"
)

type DescriptionGenerator interface {
	Generate(ctx context.Context, input GenerationInput) (*GenerationOutput, error)
	Refine(ctx context.Context, input RefinementInput) (*GenerationOutput, error)
}

type GenerationInput struct {
	Diff           string
	AnalysisType   string
	CommitMessages []string
	PRTitle        string
	PRDescription  string
	Level          string

	MaxTokensOverride   *int
	TemperatureOverride *float64
	ModelOverride       *string
}

func (i GenerationInput) HasOverrides() bool {
	if i.MaxTokensOverride != nil {
		return true
	}
	if i.TemperatureOverride != nil {
		return true
	}
	if i.ModelOverride != nil && *i.ModelOverride != "auto" {
		return true
	}
	return false
}

type GenerationOutput struct {
	Description string
	Model       string
	TokensUsed  int
}

type RefinementInput struct {
	OriginalDescription string
	Instruction         string
}

type ChatClient interface {
	CreateChatCompletion(ctx context.Context, req oai.ChatCompletionRequest) (oai.ChatCompletionResponse, error)
}

type GeneratorConfig struct {
	DefaultModel   string
	ComplexModel   string
	MaxTokens      int
	Temperature    float32
	TokenThreshold int
	CacheTTL       time.Duration
}

type openaiGenerator struct {
	client    ChatClient
	cache     cache.Cache
	config    GeneratorConfig
	modelCfg  ModelConfig
}

func NewGenerator(client ChatClient, c cache.Cache, cfg GeneratorConfig) DescriptionGenerator {
	return &openaiGenerator{
		client: client,
		cache:  c,
		config: cfg,
		modelCfg: ModelConfig{
			DefaultModel:   cfg.DefaultModel,
			ComplexModel:   cfg.ComplexModel,
			TokenThreshold: cfg.TokenThreshold,
		},
	}
}

func (g *openaiGenerator) resolveModel(input GenerationInput, tokenCount int) string {
	if input.ModelOverride != nil && *input.ModelOverride != "auto" {
		return *input.ModelOverride
	}
	return SelectModel(g.modelCfg, tokenCount, input.AnalysisType)
}

func (g *openaiGenerator) Generate(ctx context.Context, input GenerationInput) (*GenerationOutput, error) {
	level := input.Level
	if level == "" {
		level = "functional"
	}

	useCache := !input.HasOverrides()

	cacheKey := cache.DiffCacheKey(input.Diff + ":" + level)
	if useCache {
		if cached, ok := g.cache.Get(cacheKey); ok {
			slog.Debug("openai: cache hit", "key", cacheKey[:12])
			return &GenerationOutput{Description: cached, Model: "cache"}, nil
		}
	}

	processed := PreprocessDiff(input.Diff)
	tokenCount := CountTokens(processed, g.config.DefaultModel)
	model := g.resolveModel(input, tokenCount)

	maxTokens := g.config.MaxTokens
	if input.MaxTokensOverride != nil {
		maxTokens = *input.MaxTokensOverride
	}

	temperature := g.config.Temperature
	if input.TemperatureOverride != nil {
		temperature = float32(*input.TemperatureOverride)
	}

	messages := g.buildMessages(processed, input)

	var result oai.ChatCompletionResponse
	err := withRetry(ctx, func() error {
		var callErr error
		result, callErr = g.client.CreateChatCompletion(ctx, oai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			MaxTokens:   maxTokens,
			Temperature: temperature,
		})
		return callErr
	})
	if err != nil {
		return nil, err
	}

	description := result.Choices[0].Message.Content
	tokensUsed := result.Usage.TotalTokens

	slog.Info("openai: generation complete",
		"model", model,
		"tokens", tokensUsed,
		"analysis_type", input.AnalysisType,
	)

	if useCache {
		g.cache.Set(cacheKey, description, g.config.CacheTTL)
	}

	return &GenerationOutput{
		Description: description,
		Model:       model,
		TokensUsed:  tokensUsed,
	}, nil
}

func (g *openaiGenerator) Refine(ctx context.Context, input RefinementInput) (*GenerationOutput, error) {
	model := g.config.DefaultModel
	prompt := buildRefinePrompt(input.OriginalDescription, input.Instruction)

	messages := []oai.ChatCompletionMessage{
		{Role: oai.ChatMessageRoleSystem, Content: buildSystemPrompt()},
		{Role: oai.ChatMessageRoleUser, Content: prompt},
	}

	var result oai.ChatCompletionResponse
	err := withRetry(ctx, func() error {
		var callErr error
		result, callErr = g.client.CreateChatCompletion(ctx, oai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			MaxTokens:   g.config.MaxTokens,
			Temperature: g.config.Temperature,
		})
		return callErr
	})
	if err != nil {
		return nil, err
	}

	return &GenerationOutput{
		Description: result.Choices[0].Message.Content,
		Model:       model,
		TokensUsed:  result.Usage.TotalTokens,
	}, nil
}

func (g *openaiGenerator) buildMessages(diff string, input GenerationInput) []oai.ChatCompletionMessage {
	level := input.Level
	if level == "" {
		level = "functional"
	}
	messages := []oai.ChatCompletionMessage{
		{Role: oai.ChatMessageRoleSystem, Content: buildSystemPromptForLevel(level)},
	}

	for _, ex := range buildFewShotExamples() {
		messages = append(messages, oai.ChatCompletionMessage{
			Role:    ex.Role,
			Content: ex.Content,
		})
	}

	userPrompt := buildUserPrompt(diff, input.AnalysisType, input.CommitMessages, input.PRTitle, input.PRDescription)
	messages = append(messages, oai.ChatCompletionMessage{
		Role:    oai.ChatMessageRoleUser,
		Content: userPrompt,
	})

	return messages
}
