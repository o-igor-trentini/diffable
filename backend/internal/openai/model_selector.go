package openai

type ModelConfig struct {
	DefaultModel   string
	ComplexModel   string
	TokenThreshold int
}

func SelectModel(cfg ModelConfig, tokenCount int, analysisType string) string {
	if analysisType == "pull_request" {
		return cfg.ComplexModel
	}
	if tokenCount > cfg.TokenThreshold {
		return cfg.ComplexModel
	}
	return cfg.DefaultModel
}
