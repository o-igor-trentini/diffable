package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testModelCfg = ModelConfig{
	DefaultModel:   "gpt-4o-mini",
	ComplexModel:   "gpt-4o",
	TokenThreshold: 4000,
}

func TestSelectModel_SmallDiff(t *testing.T) {
	model := SelectModel(testModelCfg, 500, "single_commit")
	assert.Equal(t, "gpt-4o-mini", model)
}

func TestSelectModel_LargeDiff(t *testing.T) {
	model := SelectModel(testModelCfg, 5000, "single_commit")
	assert.Equal(t, "gpt-4o", model)
}

func TestSelectModel_PRAlwaysComplex(t *testing.T) {
	model := SelectModel(testModelCfg, 100, "pull_request")
	assert.Equal(t, "gpt-4o", model)
}

func TestSelectModel_ExactThreshold(t *testing.T) {
	model := SelectModel(testModelCfg, 4000, "single_commit")
	assert.Equal(t, "gpt-4o-mini", model)
}

func TestSelectModel_CommitRange_Small(t *testing.T) {
	model := SelectModel(testModelCfg, 2000, "commit_range")
	assert.Equal(t, "gpt-4o-mini", model)
}

func TestSelectModel_CommitRange_Large(t *testing.T) {
	model := SelectModel(testModelCfg, 4001, "commit_range")
	assert.Equal(t, "gpt-4o", model)
}
