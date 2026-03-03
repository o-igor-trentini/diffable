package openai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSystemPromptForLevel_Technical(t *testing.T) {
	prompt := buildSystemPromptForLevel("technical")

	assert.Contains(t, prompt, "engenheiro de software")
	assert.Contains(t, prompt, "Resumo Técnico")
	assert.Contains(t, prompt, "Decisões Técnicas")
}

func TestBuildSystemPromptForLevel_Functional(t *testing.T) {
	prompt := buildSystemPromptForLevel("functional")

	assert.Contains(t, prompt, "analista sênior")
	assert.Contains(t, prompt, "Resumo:")
	assert.Contains(t, prompt, "Mudanças Realizadas")
	assert.Contains(t, prompt, "Impacto Funcional")
}

func TestBuildSystemPromptForLevel_Executive(t *testing.T) {
	prompt := buildSystemPromptForLevel("executive")

	assert.Contains(t, prompt, "analista de negócios")
	assert.Contains(t, prompt, "Resumo Executivo")
	assert.Contains(t, prompt, "2-3 frases")
}

func TestBuildSystemPromptForLevel_Default(t *testing.T) {
	prompt := buildSystemPromptForLevel("")
	defaultPrompt := buildSystemPromptForLevel("functional")

	assert.Equal(t, defaultPrompt, prompt)
}

func TestBuildSystemPromptForLevel_EachLevelIsDifferent(t *testing.T) {
	technical := buildSystemPromptForLevel("technical")
	functional := buildSystemPromptForLevel("functional")
	executive := buildSystemPromptForLevel("executive")

	assert.NotEqual(t, technical, functional)
	assert.NotEqual(t, functional, executive)
	assert.NotEqual(t, technical, executive)
}

func TestBuildSystemPrompt_ReturnsFunctional(t *testing.T) {
	prompt := buildSystemPrompt()
	functional := buildSystemPromptForLevel("functional")

	assert.Equal(t, functional, prompt)
}

func TestBuildSystemPromptForLevel_AllInPortuguese(t *testing.T) {
	for _, level := range []string{"technical", "functional", "executive"} {
		prompt := buildSystemPromptForLevel(level)
		assert.True(t, strings.Contains(prompt, "Português (BR)"), "Level %s should mention Portuguese", level)
	}
}
