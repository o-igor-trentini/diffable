# Fase 2: Integração Bitbucket Cloud

**Objetivo:** Construir um cliente HTTP completo para a API REST do Bitbucket Cloud v2.0 com autenticação, paginação e tratamento de rate limiting.

**Pré-requisitos:** Fase 1 concluída

**Nota:** Esta fase pode ser desenvolvida em paralelo com a Fase 3 (OpenAI).

---

## Referência da API Bitbucket Cloud

Base URL: `https://api.bitbucket.org/2.0/`

Auth: `Authorization: Basic base64(email:api_token)`

Rate limit: 1000 req/hora (autenticado). Headers: `X-RateLimit-Limit`, `X-RateLimit-Resource`, `X-RateLimit-NearLimit`.

Paginação: campo `next` na resposta. Seguir até não existir. `pagelen` max ~100.

### Endpoints Utilizados

| Endpoint | Método | Descrição |
|----------|--------|-----------|
| `/repositories/{workspace}/{repo}/commit/{hash}` | GET | Detalhes de um commit |
| `/repositories/{workspace}/{repo}/diff/{spec}` | GET | Diff raw (text/plain) |
| `/repositories/{workspace}/{repo}/diffstat/{spec}` | GET | Diffstat (JSON) |
| `/repositories/{workspace}/{repo}/commits?include=X&exclude=Y` | GET/POST | Range de commits |
| `/repositories/{workspace}/{repo}/pullrequests/{id}` | GET | Detalhes do PR |
| `/repositories/{workspace}/{repo}/pullrequests/{id}/diff` | GET | Diff do PR (text/plain) |
| `/repositories/{workspace}/{repo}/pullrequests/{id}/commits` | GET | Commits do PR |
| `/repositories/{workspace}` | GET | Listar repositórios |

---

## Checklist

### Interface e Client

- [ ] Criar `backend/internal/bitbucket/client.go`:
  - Interface `Client` com métodos:
    ```go
    type Client interface {
        GetCommit(ctx context.Context, workspace, repoSlug, hash string) (*Commit, error)
        GetCommitDiff(ctx context.Context, workspace, repoSlug, spec string) (string, error)
        GetDiffstat(ctx context.Context, workspace, repoSlug, spec string) (*DiffstatResponse, error)
        ListCommitsInRange(ctx context.Context, workspace, repoSlug, include, exclude string) ([]Commit, error)
        GetPullRequest(ctx context.Context, workspace, repoSlug string, prID int) (*PullRequest, error)
        GetPullRequestDiff(ctx context.Context, workspace, repoSlug string, prID int) (string, error)
        GetPullRequestCommits(ctx context.Context, workspace, repoSlug string, prID int) ([]Commit, error)
        ListRepositories(ctx context.Context, workspace string) ([]Repository, error)
    }
    ```
  - Struct `bitbucketClient` implementando a interface
  - Construtor `NewClient(config BitbucketConfig) Client`
  - Auth header: `Authorization: Basic base64(email:apiToken)` em toda request
  - HTTP client com timeout configurável
  - Logging estruturado (slog) para cada request

### Types

- [ ] Criar `backend/internal/bitbucket/types.go`:
  - `Commit` — hash, message, date, author, parents
  - `Author` — display_name, uuid
  - `PullRequest` — id, title, description, state, source/destination branches, author, created_on
  - `DiffstatResponse` — paginated, com `DiffstatEntry` (status, lines_added, lines_removed, paths)
  - `Repository` — slug, name, full_name, description
  - `PaginatedResponse[T]` — values, next, pagelen, size, page

### Paginação

- [ ] Criar `backend/internal/bitbucket/pagination.go`:
  - Função genérica `FetchAllPages[T](client, initialURL) ([]T, error)`
  - Segue campo `next` até nil
  - Respeita `pagelen` (default 100)
  - Context-aware (cancela se ctx cancelado)
- [ ] Criar `backend/internal/bitbucket/pagination_test.go`:
  - Mock server retorna 3 páginas → função coleta todos items
  - Edge case: primeira página vazia → retorna slice vazio
  - Edge case: `next` ausente na primeira página → retorna apenas essa página

### Rate Limiting

- [ ] Criar `backend/internal/bitbucket/ratelimit.go`:
  - Parse headers `X-RateLimit-Limit`, `X-RateLimit-NearLimit`
  - Struct `RateLimitInfo` com Limit, Remaining, NearLimit
  - Se `NearLimit == true`, emite `slog.Warn`
  - Se resposta 429, retorna `ErrRateLimited` com informação de retry
- [ ] Criar `backend/internal/bitbucket/ratelimit_test.go`:
  - Testa parse correto de headers
  - Testa warning quando NearLimit é true
  - Testa detecção de 429

### Errors

- [ ] Criar `backend/internal/bitbucket/errors.go`:
  - `ErrRateLimited` — com RetryAfter duration
  - `ErrNotFound` — recurso não encontrado (404)
  - `ErrUnauthorized` — credenciais inválidas (401)
  - `ErrBitbucketAPI` — erro genérico com status code e mensagem

### Configuração

- [ ] Atualizar `backend/internal/config/config.go` — adicionar:
  - `BitbucketBaseURL` (default: `https://api.bitbucket.org/2.0`)
  - `BitbucketEmail` (env: `BITBUCKET_EMAIL`)
  - `BitbucketAPIToken` (env: `BITBUCKET_API_TOKEN`)
  - `BitbucketTimeout` (default: 30s)
- [ ] Atualizar `backend/.env.example` com as novas variáveis

### Documentação

- [ ] Criar `docs/bitbucket-integration.md`:
  - Todos endpoints utilizados com exemplos de request/response
  - Formato de autenticação
  - Estratégia de rate limiting
  - Abordagem de paginação
  - Matriz de tratamento de erros
  - Como rodar testes de integração

---

## Testes desta Fase

| Teste | Tipo | Validação |
|-------|------|-----------|
| `client_test.go` — GetCommit | Unit | URL path correto, auth header, JSON parsing, 404 |
| `client_test.go` — GetCommitDiff | Unit | Retorna raw text diff, trata diffs grandes |
| `client_test.go` — GetDiffstat | Unit | JSON parsing correto |
| `client_test.go` — ListCommitsInRange | Unit | Query params `include`/`exclude`, paginação |
| `client_test.go` — GetPullRequest | Unit | URL com PR ID, JSON mapping |
| `client_test.go` — GetPullRequestDiff | Unit | Retorna raw text diff |
| `client_test.go` — GetPullRequestCommits | Unit | Paginação através dos commits |
| `client_test.go` — ListRepositories | Unit | Paginação, filtro |
| `client_test.go` — Auth header | Unit | Header formatado como `Basic base64(email:token)` |
| `client_test.go` — Error 401 | Unit | Retorna `ErrUnauthorized` |
| `client_test.go` — Error 404 | Unit | Retorna `ErrNotFound` |
| `client_test.go` — Error 429 | Unit | Retorna `ErrRateLimited` |
| `client_test.go` — Error 500 | Unit | Retorna `ErrBitbucketAPI` |
| `pagination_test.go` — multi-page | Unit | 3 páginas coletadas corretamente |
| `pagination_test.go` — empty | Unit | Página vazia retorna slice vazio |
| `ratelimit_test.go` — headers | Unit | Parse correto dos headers |
| `ratelimit_test.go` — threshold | Unit | Warning emitido quando NearLimit |
| Integration (opcional) | Integration | Chamada real a repo público, skip via env `BITBUCKET_INTEGRATION_TEST` |

> Todos testes unitários usam `httptest.NewServer` para mock do Bitbucket.

---

## Critérios de Aceite

- [ ] Todos métodos da interface `Client` possuem testes unitários passando
- [ ] Mock server prova construção correta de URL para cada endpoint
- [ ] Auth header formatado corretamente como `Basic base64(email:token)`
- [ ] Paginação segue todas páginas e coleta resultados completos
- [ ] Headers de rate limit são parseados e warning logado abaixo do threshold
- [ ] Respostas 429 resultam em `ErrRateLimited` com info de retry
- [ ] Tipos de erro são específicos (`ErrNotFound`, `ErrUnauthorized`, `ErrRateLimited`)
- [ ] `docs/bitbucket-integration.md` completo e revisado
- [ ] `go test ./internal/bitbucket/...` passa 100%
