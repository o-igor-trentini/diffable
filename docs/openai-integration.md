# Integração OpenAI

## Visão Geral

Serviço de geração de descrições para cards JIRA a partir de diffs de código, usando a API da OpenAI com seleção automática de modelo, preprocessing de diff, cache por SHA-256, e retry com exponential backoff.

## Design do Prompt

### System Prompt

O modelo atua como **analista sênior de software** escrevendo para JIRA:
- Linguagem não-técnica em Português (BR)
- Acessível para QA e PO
- Foco no impacto funcional, não em detalhes de implementação
- Máximo ~200 palavras
- Formato: **Resumo** + **Mudanças Realizadas** (bullets) + **Impacto Funcional**

### Few-Shot Examples

2 pares de exemplo diff→descrição são incluídos para guiar o formato e tom:
1. Login com bloqueio de conta (commit simples)
2. Adição de PIX no checkout (PR)

### User Prompt

Varia conforme o tipo de análise:
- `single_commit`: "Descreva as mudanças deste commit:"
- `commit_range`: "Descreva as mudanças consolidadas destes commits:"
- `pull_request`: "Gere uma descrição para card JIRA baseado neste PR:"

Contexto adicional (título do PR, descrição, mensagens de commit) é incluído quando disponível.

## Seleção de Modelo

| Condição | Modelo |
|----------|--------|
| `analysisType == "pull_request"` | `gpt-4o` |
| `tokenCount > 4000` | `gpt-4o` |
| Caso contrário | `gpt-4o-mini` |

O threshold é configurável via `OPENAI_TOKEN_THRESHOLD`.

## Gerenciamento de Tokens

### Preprocessing de Diff

Antes da contagem de tokens, o diff passa por preprocessing que remove:
- Arquivos de teste (`*_test.go`, `*.spec.ts`, `*.test.tsx`, `__snapshots__`)
- Binários (`Binary files ... differ`)
- Lock files (`package-lock.json`, `go.sum`, `yarn.lock`, `pnpm-lock.yaml`)
- Arquivos gerados (`vendor/`, `dist/`, `build/`, `node_modules/`)
- Linhas de contexto (mantém apenas `+`, `-`, `@@`, headers)

### Contagem de Tokens

Usa `pkoukk/tiktoken-go` para contagem precisa compatível com os modelos OpenAI.

### Chunking

Para diffs muito grandes, `ChunkDiff` divide por boundary de arquivo (`diff --git`), mantendo cada chunk dentro do limite de tokens.

## Cache

- **Chave:** SHA-256 hex do diff raw
- **Implementação:** `sync.Map` in-memory com TTL
- **TTL padrão:** 24h (configurável via `CACHE_TTL`)
- Diffs idênticos retornam resultado cacheado sem chamar a API

## Retry

- **Estratégia:** Exponential backoff com jitter
- **Delays:** 1s, 2s, 4s (com +/- 20% de jitter)
- **Max tentativas:** 3
- **Retryable:** 429 (rate limit), 500, 503
- **Não retryable:** 400, 401, 404
- Context-aware (cancela se ctx cancelado)

## Configuração

| Variável | Default | Descrição |
|----------|---------|-----------|
| `OPENAI_API_KEY` | — | Chave da API OpenAI |
| `OPENAI_DEFAULT_MODEL` | `gpt-4o-mini` | Modelo para diffs simples |
| `OPENAI_COMPLEX_MODEL` | `gpt-4o` | Modelo para diffs complexos/PRs |
| `OPENAI_MAX_TOKENS` | `1024` | Max tokens no output |
| `OPENAI_TEMPERATURE` | `0.3` | Temperatura (consistência) |
| `OPENAI_TOKEN_THRESHOLD` | `4000` | Threshold para modelo complexo |
| `CACHE_TTL` | `24h` | TTL do cache de descrições |

## Estimativa de Custo

| Tipo | Modelo | Input estimado | Output estimado | Custo aprox. |
|------|--------|---------------|-----------------|--------------|
| Commit simples | gpt-4o-mini | ~500 tokens | ~300 tokens | ~$0.0003 |
| Commit range | gpt-4o-mini | ~2000 tokens | ~400 tokens | ~$0.001 |
| Pull Request | gpt-4o | ~3000 tokens | ~500 tokens | ~$0.02 |
| Refinamento | gpt-4o-mini | ~500 tokens | ~300 tokens | ~$0.0003 |

## Exemplos

### Input (Commit)

```json
{
  "diff": "diff --git a/src/auth/login.ts ...",
  "analysis_type": "single_commit",
  "commit_messages": ["feat: add account lockout"]
}
```

### Output

```json
{
  "description": "**Resumo:** Adicionado mecanismo de bloqueio...",
  "model": "gpt-4o-mini",
  "tokens_used": 150
}
```

## Testes

```bash
go test ./internal/openai/... -v
go test ./internal/cache/... -v
```

## Estrutura de Arquivos

```
backend/internal/openai/
├── generator.go            # Interface DescriptionGenerator e implementação
├── generator_test.go       # Testes com mock client
├── prompts.go              # System prompt, few-shot, user prompts
├── tokenizer.go            # CountTokens, PreprocessDiff, ChunkDiff
├── tokenizer_test.go       # Testes de tokenização e preprocessing
├── model_selector.go       # SelectModel por tokens e tipo
├── model_selector_test.go  # Testes de seleção de modelo
├── retry.go                # withRetry com exponential backoff
└── retry_test.go           # Testes de retry

backend/internal/cache/
├── cache.go                # Interface Cache, InMemoryCache, DiffCacheKey
└── cache_test.go           # Testes de cache com TTL e concorrência
```
