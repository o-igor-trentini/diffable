# Fase 3: Integração OpenAI

**Objetivo:** Construir o serviço de geração de descrições usando OpenAI, com interface abstrata, tokenização, preprocessing de diff, seleção de modelo, retry com backoff e cache.

**Pré-requisitos:** Fase 1 concluída

**Nota:** Esta fase pode ser desenvolvida em paralelo com a Fase 2 (Bitbucket).

---

## Decisões Técnicas

| Decisão | Escolha | Motivo |
|---------|---------|--------|
| Client Go | `github.com/sashabaranov/go-openai` | 9k+ stars, mantido ativamente |
| Modelo default | `gpt-4o-mini` | Custo-benefício para diffs pequenos |
| Modelo complexo | `gpt-4o` | Melhor qualidade para diffs grandes/PRs |
| Threshold | 4000 tokens | Acima disso, usa modelo complexo |
| Temperatura | 0.3 | Output consistente e conciso |
| Max tokens output | 1024 | Suficiente para descrições JIRA |
| Streaming | Não | Background processing, não precisa de UX real-time |
| Retry | Exponential backoff + jitter, max 3 | Padrão para APIs externas |
| Cache | SHA-256 do diff como chave | Evita reprocessamento de diffs idênticos |
| Token counting | `pkoukk/tiktoken-go` | Contagem precisa compatível com OpenAI |

---

## Checklist

### Interface e Gerador

- [ ] Criar `backend/internal/openai/generator.go`:
  - Interface `DescriptionGenerator`:
    ```go
    type DescriptionGenerator interface {
        Generate(ctx context.Context, input GenerationInput) (*GenerationOutput, error)
        Refine(ctx context.Context, input RefinementInput) (*GenerationOutput, error)
    }

    type GenerationInput struct {
        Diff           string
        AnalysisType   string   // "single_commit", "commit_range", "pull_request"
        CommitMessages []string // opcional, para contexto
        PRTitle        string   // opcional, para análise de PR
        PRDescription  string   // opcional, para análise de PR
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
    ```
  - Struct `OpenAIGenerator` implementando a interface
  - Usa `sashabaranov/go-openai` para chamadas
  - Temperatura 0.3, MaxTokens 1024
  - Logging de tokens usados e custo estimado
- [ ] Criar `backend/internal/openai/generator_test.go`:
  - Geração com sucesso (mock client)
  - Cache hit retorna resultado cacheado (generator NÃO chamado)
  - Retry após 429 (sucesso na segunda tentativa)
  - Refinamento com instrução válida
  - Context cancelado aborta operação

### Prompts

- [ ] Criar `backend/internal/openai/prompts.go`:
  - `buildSystemPrompt()` — system role:
    - Atua como analista sênior escrevendo para JIRA
    - Linguagem não-técnica, acessível para QA/PO
    - Output em Português (BR)
    - Formato: Summary, Mudanças Realizadas (bullets), Impacto Funcional
    - Foco no impacto funcional, não detalhes de implementação
    - Máximo ~200 palavras
  - `buildFewShotExamples()` — 2-3 pares diff→descrição de exemplo
  - `buildUserPrompt(diff, analysisType)` — prompt do usuário por tipo:
    - `single_commit`: "Descreva as mudanças deste commit:"
    - `commit_range`: "Descreva as mudanças consolidadas destes commits:"
    - `pull_request`: "Gere uma descrição para card JIRA baseado neste PR:"
  - `buildRefinePrompt(original, instruction)` — "Refine esta descrição: [original]. Instrução: [instrução]"

### Tokenização e Preprocessing

- [ ] Criar `backend/internal/openai/tokenizer.go`:
  - `CountTokens(text, model string) int` — usa tiktoken-go
  - `PreprocessDiff(rawDiff string) string`:
    1. Remove diffs de arquivos de teste (`*_test.go`, `*.spec.ts`, `*.test.tsx`, `__snapshots__`)
    2. Remove diffs de binários (`Binary files ... differ`)
    3. Remove diffs de lock files (`package-lock.json`, `go.sum`, `yarn.lock`, `pnpm-lock.yaml`)
    4. Remove diffs de arquivos gerados (vendor, dist, build)
    5. Reduz linhas de contexto (mantém apenas +/- e headers)
  - `ChunkDiff(diff string, maxTokens int) []string` — split por boundary de arquivo (`diff --git`), cada chunk dentro do limite
- [ ] Criar `backend/internal/openai/tokenizer_test.go`:
  - CountTokens retorna valor esperado para strings conhecidas
  - PreprocessDiff remove arquivos de teste
  - PreprocessDiff remove binários
  - PreprocessDiff remove lock files
  - ChunkDiff respeita limite de tokens
  - ChunkDiff com diff pequeno retorna 1 chunk
  - Diff vazio tratado corretamente

### Seleção de Modelo

- [ ] Criar `backend/internal/openai/model_selector.go`:
  - `SelectModel(tokenCount int, analysisType string) string`:
    - Se `tokenCount > threshold` (4000) → `gpt-4o`
    - Se `analysisType == "pull_request"` → `gpt-4o`
    - Caso contrário → `gpt-4o-mini`
  - Threshold configurável via config
- [ ] Criar `backend/internal/openai/model_selector_test.go`:
  - Diff pequeno → gpt-4o-mini
  - Diff grande (>4000 tokens) → gpt-4o
  - PR independente do tamanho → gpt-4o
  - Threshold boundary (exatamente 4000) → gpt-4o-mini

### Retry

- [ ] Criar `backend/internal/openai/retry.go`:
  - `withRetry(ctx, maxAttempts int, fn func() error) error`:
    - Exponential backoff: 1s, 2s, 4s
    - Jitter: +/- 20% do delay calculado
    - Retry em: 429 (rate limit), 500, 503
    - Não retry em: 400, 401, 404
    - Respeita context cancellation
    - Max 3 tentativas
- [ ] Criar `backend/internal/openai/retry_test.go`:
  - Sucesso na primeira tentativa → sem delay
  - Sucesso na segunda tentativa → delay ~1s
  - Falha após 3 tentativas → retorna último erro
  - Erro não-retryável (400) → falha imediata
  - Context cancelado → retorna ctx.Err()

### Cache

- [ ] Criar `backend/internal/cache/cache.go`:
  - Interface:
    ```go
    type Cache interface {
        Get(key string) (string, bool)
        Set(key, value string, ttl time.Duration)
    }
    ```
  - `InMemoryCache` com `sync.Map` e TTL
  - `DiffCacheKey(diff string) string` — retorna SHA-256 hex do diff
- [ ] Criar `backend/internal/cache/cache_test.go`:
  - Set/Get funciona
  - TTL expira corretamente
  - Cache miss retorna false
  - Acesso concorrente seguro

### Configuração

- [ ] Atualizar `backend/internal/config/config.go` — adicionar:
  - `OpenAIAPIKey` (env: `OPENAI_API_KEY`)
  - `OpenAIDefaultModel` (default: `gpt-4o-mini`)
  - `OpenAIComplexModel` (default: `gpt-4o`)
  - `OpenAIMaxTokens` (default: 1024)
  - `OpenAITemperature` (default: 0.3)
  - `OpenAITokenThreshold` (default: 4000)
  - `CacheTTL` (default: 24h)
- [ ] Atualizar `backend/.env.example` com as novas variáveis

### Documentação

- [ ] Criar `docs/openai-integration.md`:
  - Design do prompt (rationale, exemplos)
  - Lógica de seleção de modelo
  - Estratégia de gerenciamento de tokens
  - Abordagem de cache
  - Política de retry
  - Estimativa de custo por tipo de análise
  - Exemplos de input/output

---

## Testes desta Fase

| Teste | Tipo | Validação |
|-------|------|-----------|
| `generator_test.go` — Generate success | Unit | Mock client retorna completion, output parseado |
| `generator_test.go` — Cache hit | Unit | Segunda chamada com mesmo diff retorna cache |
| `generator_test.go` — Retry on 429 | Unit | 429 na primeira, sucesso na segunda |
| `generator_test.go` — Refine | Unit | Prompt de refinamento construído, output retornado |
| `generator_test.go` — Context cancel | Unit | Operação abortada quando ctx cancelado |
| `tokenizer_test.go` — CountTokens | Unit | String conhecida retorna contagem esperada |
| `tokenizer_test.go` — PreprocessDiff | Unit | Arquivos de teste removidos, binários removidos |
| `tokenizer_test.go` — ChunkDiff | Unit | Diff grande dividido, cada chunk no limite |
| `model_selector_test.go` — thresholds | Unit | Cada boundary retorna modelo correto |
| `cache_test.go` — CRUD + TTL | Unit | Set, Get, miss, expiração |
| `retry_test.go` — backoff | Unit | Progressão de delay correta, jitter, max tentativas |

---

## Critérios de Aceite

- [ ] Interface `DescriptionGenerator` é limpa e agnóstica de implementação
- [ ] Implementação OpenAI gera descrições a partir de diff text
- [ ] Token counting coincide com valores esperados para inputs conhecidos
- [ ] Preprocessor de diff remove ruído (testes, binários, lock files)
- [ ] Model selector escolhe `gpt-4o-mini` para diffs pequenos, `gpt-4o` para complexos
- [ ] Cache retorna resultado armazenado para diff idêntico (chave SHA-256)
- [ ] Retry faz backoff exponencial com jitter, para após 3 tentativas
- [ ] Todos testes unitários passam com client OpenAI mockado
- [ ] `docs/openai-integration.md` completo
- [ ] `go test ./internal/openai/... ./internal/cache/...` passa 100%
