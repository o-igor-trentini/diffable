# Fase 7: Melhorias Futuras (Opcional)

**Objetivo:** Features avançadas do backlog que melhoram significativamente a usabilidade. Cada sub-fase é independente e pode ser implementada individualmente.

**Pré-requisitos:** Fase 6 concluída

---

## 7a. Export Markdown

**Objetivo:** Permitir download da descrição gerada como arquivo `.md` formatado.

### Checklist

- [ ] Criar botão "Exportar Markdown" no `ResultDisplay.tsx` (ícone FileDown do Lucide)
- [ ] Gerar conteúdo Markdown com metadados:
  ```markdown
  # Análise: [tipo]
  **Data:** [created_at]
  **Modelo:** [model_used]
  **Repositório:** [workspace/repo]

  ---

  [descrição gerada]
  ```
- [ ] Download via Blob URL no browser
- [ ] Teste: botão gera e faz download do arquivo .md com conteúdo correto

### Critérios de Aceite

- [ ] Botão "Exportar Markdown" visível no ResultDisplay
- [ ] Click gera download de arquivo `.md` com formatação correta
- [ ] Metadados (tipo, data, modelo, repo) incluídos no cabeçalho

---

## 7b. Geração em Múltiplos Níveis

**Objetivo:** Permitir escolha do nível de descrição: Técnico, Funcional ou Executivo.

### Checklist

- [ ] Adicionar dropdown "Nível" nos formulários de cada tab:
  - Técnico: inclui referências a código, padrões, impacto técnico
  - Funcional: default atual (QA/PO, linguagem não-técnica)
  - Executivo: alto nível, impacto de negócio, 2-3 frases máximo
- [ ] Criar prompts de sistema específicos por nível em `prompts.go`:
  - `buildSystemPromptForLevel(level string) string`
  - Nível "technical": pode mencionar arquivos, padrões, trade-offs
  - Nível "functional": default atual
  - Nível "executive": "Resuma em 2-3 frases o impacto de negócio"
- [ ] Adicionar campo `level` no `GenerationInput`:
  ```go
  type GenerationInput struct {
      // ... campos existentes
      Level string // "technical", "functional", "executive"
  }
  ```
- [ ] Adicionar campo `level` nos DTOs de request e na tabela `analyses`
- [ ] Frontend: dropdown com 3 opções, default "Funcional"
- [ ] Teste: cada nível gera prompt de sistema diferente
- [ ] Teste: dropdown renderiza, seleção muda valor enviado

### Critérios de Aceite

- [ ] Dropdown "Nível" visível em todos formulários
- [ ] Cada nível gera descrição com tom e profundidade diferentes
- [ ] Nível salvo na análise para referência futura

---

## 7c. Autocomplete de Workspace/Repositório

**Objetivo:** Facilitar preenchimento dos campos workspace e repositório com autocomplete baseado na API do Bitbucket.

### Checklist

**Backend:**
- [ ] Criar endpoint `GET /api/v1/bitbucket/repositories?workspace={workspace}&q={query}`:
  - Proxy para `BitbucketClient.ListRepositories`
  - Filtro por nome (query string `q`)
  - Cache da lista por 5 minutos
- [ ] Teste: endpoint retorna lista filtrada de repositórios

**Frontend:**
- [ ] Criar componente `AutocompleteInput.tsx`:
  - Input com dropdown de sugestões
  - Debounce de 300ms no input
  - Loading state enquanto busca
  - Seleção por click ou Enter
  - Limpa sugestões ao selecionar
- [ ] Substituir inputs de texto de workspace/repo por `AutocompleteInput`
- [ ] Hook `useRepositories(workspace, query)` — `useQuery` com `enabled: query.length >= 2`
- [ ] Teste: debounce funciona, sugestões aparecem, seleção preenche campo

### Critérios de Aceite

- [ ] Campo workspace e repo_slug oferecem autocomplete
- [ ] Busca acontece após 2+ caracteres com debounce de 300ms
- [ ] Lista de repositórios cacheada por 5 minutos
- [ ] Seleção preenche campo corretamente

---

## 7d. Webhook para PRs

**Objetivo:** Receber webhooks do Bitbucket quando um PR é criado e auto-gerar descrição.

### Checklist

**Backend:**
- [ ] Criar endpoint `POST /api/v1/webhooks/bitbucket`:
  - Parse payload do webhook de PR do Bitbucket Cloud
  - Valida header `X-Event-Key` = `pullrequest:created`
  - Valida assinatura (se configurada)
  - Extrai workspace, repo, PR ID do payload
  - Chama `analysisService.AnalyzePR` em background (goroutine)
  - Retorna 202 Accepted imediatamente
- [ ] Criar `backend/internal/webhook/bitbucket.go`:
  - Struct para payload do webhook
  - Parser e validador
  - Background processing com error logging
- [ ] Criar migration para tabela de webhook logs (opcional):
  ```sql
  CREATE TABLE webhook_logs (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      event_key VARCHAR(100),
      payload JSONB,
      status VARCHAR(20),
      analysis_id UUID REFERENCES analyses(id),
      error_message TEXT,
      created_at TIMESTAMPTZ DEFAULT NOW()
  );
  ```
- [ ] Documentar setup do webhook no Bitbucket:
  - Settings → Webhooks → Add webhook
  - URL: `https://your-domain/api/v1/webhooks/bitbucket`
  - Triggers: Pull Request → Created
- [ ] Teste: payload válido → 202, análise criada em background
- [ ] Teste: event key inválido → 400
- [ ] Teste: payload malformado → 400

### Critérios de Aceite

- [ ] Webhook recebe evento `pullrequest:created` e retorna 202
- [ ] Análise é gerada em background automaticamente
- [ ] Resultado pode ser consultado via `GET /api/v1/analyses`
- [ ] Documentação explica como configurar webhook no Bitbucket

---

## 7e. Integração CI/CD

**Objetivo:** Permitir geração de descrições via pipeline CI (GitHub Actions, Bitbucket Pipelines).

### Checklist

- [ ] Criar CLI tool: `bb-gen-desc-cli`:
  - Flag `--workspace`, `--repo`, `--pr-id` ou `--commit-hash`
  - Flag `--api-url` para apontar para a plataforma
  - Flag `--output` para formato (text, json, markdown)
  - Faz POST para API e imprime resultado
- [ ] Criar exemplo de Bitbucket Pipeline step:
  ```yaml
  - step:
      name: Generate PR Description
      script:
        - bb-gen-desc-cli --api-url $BB_GEN_DESC_URL --workspace $BITBUCKET_WORKSPACE --repo $BITBUCKET_REPO_SLUG --pr-id $BITBUCKET_PR_ID
  ```
- [ ] Documentar uso em CI/CD

### Critérios de Aceite

- [ ] CLI funciona standalone e gera descrição via API
- [ ] Exemplo de pipeline documentado e funcional

---

## Priorização Sugerida

| Sub-fase | Impacto | Esforço | Prioridade |
|----------|---------|---------|------------|
| 7b. Múltiplos Níveis | Alto | Baixo | 1 |
| 7a. Export Markdown | Médio | Baixo | 2 |
| 7c. Autocomplete | Médio | Médio | 3 |
| 7d. Webhook PRs | Alto | Alto | 4 |
| 7e. CLI/CI | Médio | Alto | 5 |
