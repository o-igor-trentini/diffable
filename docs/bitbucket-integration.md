# Integração Bitbucket Cloud

## Visão Geral

Cliente HTTP para a API REST do Bitbucket Cloud v2.0 com autenticação Basic Auth, paginação automática e tratamento de rate limiting.

## Autenticação

Todas as requisições utilizam Basic Auth:

```
Authorization: Basic base64(email:api_token)
```

Configuração via variáveis de ambiente:

| Variável | Descrição | Default |
|----------|-----------|---------|
| `BITBUCKET_BASE_URL` | URL base da API | `https://api.bitbucket.org/2.0` |
| `BITBUCKET_EMAIL` | Email da conta Bitbucket | — |
| `BITBUCKET_API_TOKEN` | App password / API token | — |
| `BITBUCKET_TIMEOUT` | Timeout HTTP | `30s` |

## Endpoints Utilizados

### Commits

**GET** `/repositories/{workspace}/{repo}/commit/{hash}`

Retorna detalhes de um commit específico.

```json
{
  "hash": "abc123...",
  "message": "fix: corrige validação",
  "date": "2025-01-01T00:00:00+00:00",
  "author": {
    "raw": "User <user@example.com>",
    "user": { "display_name": "User", "uuid": "{...}" }
  },
  "parents": [{ "hash": "def456..." }]
}
```

### Diffs

**GET** `/repositories/{workspace}/{repo}/diff/{spec}`

Retorna o diff raw em formato text/plain.

**GET** `/repositories/{workspace}/{repo}/diffstat/{spec}`

Retorna estatísticas do diff em JSON (arquivos alterados, linhas adicionadas/removidas).

### Range de Commits

**GET** `/repositories/{workspace}/{repo}/commits?include=X&exclude=Y`

Lista commits entre dois pontos. Suporta paginação automática.

### Pull Requests

**GET** `/repositories/{workspace}/{repo}/pullrequests/{id}`

Detalhes do PR (título, descrição, estado, branches).

**GET** `/repositories/{workspace}/{repo}/pullrequests/{id}/diff`

Diff completo do PR em text/plain.

**GET** `/repositories/{workspace}/{repo}/pullrequests/{id}/commits`

Lista de commits do PR com paginação automática.

### Repositórios

**GET** `/repositories/{workspace}`

Lista repositórios do workspace com paginação automática.

## Paginação

A API do Bitbucket utiliza paginação baseada em cursor. O campo `next` na resposta indica a URL da próxima página.

**Estratégia:**
- Segue automaticamente o campo `next` até que não exista mais
- `pagelen` configurado como 100 (máximo)
- Respeita cancelamento de context (context-aware)
- Acumula todos os `values` de todas as páginas

## Rate Limiting

Limite: **1000 requisições/hora** (autenticado).

**Headers monitorados:**
- `X-RateLimit-Limit`: limite total
- `X-RateLimit-NearLimit`: `true` quando próximo do limite

**Comportamento:**
- Quando `NearLimit == true`: emite warning via `slog.Warn`
- Quando recebe HTTP 429: retorna `RateLimitedError` com duração de retry (header `Retry-After` ou 60s default)

## Tratamento de Erros

| Status HTTP | Erro | Descrição |
|-------------|------|-----------|
| 401 | `UnauthorizedError` | Credenciais inválidas |
| 404 | `NotFoundError` | Recurso não encontrado |
| 429 | `RateLimitedError` | Rate limit excedido |
| 4xx/5xx | `APIError` | Erro genérico com status code e mensagem |

Todos os erros implementam a interface `error` e podem ser verificados com `errors.As()`.

## Testes

### Executar testes unitários

```bash
cd backend
go test ./internal/bitbucket/... -v
```

### Testes de integração (opcional)

Para rodar testes contra a API real do Bitbucket (requer credenciais):

```bash
BITBUCKET_INTEGRATION_TEST=true \
BITBUCKET_EMAIL=user@example.com \
BITBUCKET_API_TOKEN=your-token \
go test ./internal/bitbucket/... -v -run Integration
```

## Estrutura de Arquivos

```
backend/internal/bitbucket/
├── client.go           # Interface Client e implementação HTTP
├── client_test.go      # Testes do client com httptest
├── types.go            # Tipos de dados (Commit, PR, Repository, etc.)
├── pagination.go       # Função genérica fetchAllPages
├── pagination_test.go  # Testes de paginação
├── ratelimit.go        # Parse de headers e detecção de rate limit
├── ratelimit_test.go   # Testes de rate limit
└── errors.go           # Tipos de erro específicos
```
