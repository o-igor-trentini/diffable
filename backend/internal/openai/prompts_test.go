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
	for _, level := range []string{"technical", "functional", "executive", "qa_detailed"} {
		prompt := buildSystemPromptForLevel(level)
		assert.True(t, strings.Contains(prompt, "Português (BR)"), "Level %s should mention Portuguese", level)
	}
}

func TestBuildSystemPromptForLevel_QADetailed(t *testing.T) {
	prompt := buildSystemPromptForLevel("qa_detailed")

	assert.Contains(t, prompt, "analista de qualidade")
	assert.Contains(t, prompt, "Contexto:")
	assert.Contains(t, prompt, "Mudanças no Banco de Dados:")
	assert.Contains(t, prompt, "Mudanças de API/Integrações:")
	assert.Contains(t, prompt, "Regras de Negócio:")
	assert.Contains(t, prompt, "Fluxos Afetados:")
	assert.Contains(t, prompt, "Cenários de Teste Sugeridos:")
	assert.Contains(t, prompt, "Observações:")
}

func TestBuildFewShotExamplesForLevel_QADetailed(t *testing.T) {
	qaExamples := buildFewShotExamplesForLevel("qa_detailed")
	functionalExamples := buildFewShotExamplesForLevel("functional")

	assert.NotEqual(t, qaExamples, functionalExamples, "qa_detailed should have different examples than functional")
	assert.Greater(t, len(qaExamples), 0, "qa_detailed should have at least one example")
}

func TestBuildUserPrompt_WithUserContext(t *testing.T) {
	prompt := buildUserPrompt("diff content", "pull_request", nil, "PR Title", "", "Este PR integra CNH da API Nexus")

	assert.Contains(t, prompt, "Contexto adicional fornecido pelo desenvolvedor:")
	assert.Contains(t, prompt, "Este PR integra CNH da API Nexus")
}

func TestBuildUserPrompt_WithoutUserContext(t *testing.T) {
	prompt := buildUserPrompt("diff content", "pull_request", nil, "PR Title", "", "")

	assert.NotContains(t, prompt, "Contexto adicional fornecido pelo desenvolvedor:")
}
