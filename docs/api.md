# API Reference

Base URL: `http://localhost:8080`

Todos os endpoints retornam JSON. Erros seguem o formato padrão `ErrorResponse`.

---

## Health Checks

### GET /healthz

Liveness probe.

```bash
curl http://localhost:8080/healthz
```

**Response 200:**
```json
{"status": "ok"}
```

### GET /readyz

Readiness probe (verifica conectividade com o banco).

```bash
curl http://localhost:8080/readyz
```

**Response 200:**
```json
{"status": "ok"}
```

**Response 503:**
```json
{"status": "unavailable"}
```

---

## Análises

### POST /api/v1/analyses/commit

Analisa um commit único e gera descrição.

**Request Body:**

Opção 1 — Via Bitbucket:
```json
{
  "workspace": "meu-workspace",
  "repo_slug": "meu-repo",
  "commit_hash": "abc1234def5678"
}
```

Opção 2 — Diff direto:
```json
{
  "raw_diff": "diff --git a/file.go ..."
}
```

```bash
curl -X POST http://localhost:8080/api/v1/analyses/commit \
  -H "Content-Type: application/json" \
  -d '{"workspace":"my-ws","repo_slug":"my-repo","commit_hash":"abc1234"}'
```

**Response 200:**
```json
{
  "id": "uuid",
  "type": "single_commit",
  "description": "Texto gerado...",
  "model_used": "gpt-4o-mini",
  "tokens_used": 150,
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Erros:** 400 (validação), 502 (Bitbucket/OpenAI indisponível), 500 (erro interno)

---

### POST /api/v1/analyses/range

Analisa um range de commits.

**Request Body:**
```json
{
  "workspace": "meu-workspace",
  "repo_slug": "meu-repo",
  "from_hash": "abc1234",
  "to_hash": "def5678"
}
```

```bash
curl -X POST http://localhost:8080/api/v1/analyses/range \
  -H "Content-Type: application/json" \
  -d '{"workspace":"my-ws","repo_slug":"my-repo","from_hash":"abc1234","to_hash":"def5678"}'
```

**Response 200:** Mesmo formato de `AnalysisResponse` com `type: "commit_range"`.

**Erros:** 400, 502, 500

---

### POST /api/v1/analyses/pr

Analisa um Pull Request.

**Request Body:**

Opção 1 — Via Bitbucket:
```json
{
  "workspace": "meu-workspace",
  "repo_slug": "meu-repo",
  "pr_id": 42
}
```

Opção 2 — Diff direto:
```json
{
  "raw_diff": "diff --git ...",
  "pr_title": "feat: nova funcionalidade",
  "pr_description": "Descrição opcional do PR"
}
```

```bash
curl -X POST http://localhost:8080/api/v1/analyses/pr \
  -H "Content-Type: application/json" \
  -d '{"workspace":"my-ws","repo_slug":"my-repo","pr_id":42}'
```

**Response 200:** Mesmo formato de `AnalysisResponse` com `type: "pull_request"`.

**Erros:** 400, 502, 500

---

### GET /api/v1/analyses/{id}

Busca uma análise por ID.

```bash
curl http://localhost:8080/api/v1/analyses/uuid-aqui
```

**Response 200:** `AnalysisResponse`

**Erros:** 400 (id inválido), 404 (não encontrado), 500

---

### GET /api/v1/analyses

Lista análises com paginação e filtro por tipo.

**Query Parameters:**

| Param     | Tipo   | Default | Descrição                                       |
|-----------|--------|---------|------------------------------------------------|
| type      | string | —       | Filtro: `single_commit`, `commit_range`, `pull_request` |
| page      | int    | 1       | Página atual (>= 1)                             |
| page_size | int    | 20      | Itens por página (1-100)                         |

```bash
curl "http://localhost:8080/api/v1/analyses?type=pull_request&page=1&page_size=10"
```

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "type": "pull_request",
      "description": "...",
      "model_used": "gpt-4o-mini",
      "tokens_used": 150,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 42,
  "page": 1,
  "page_size": 10
}
```

**Erros:** 400 (parâmetros inválidos), 500

---

## Refinamentos

### POST /api/v1/analyses/{id}/refine

Refina a descrição de uma análise existente.

**Request Body:**
```json
{
  "instruction": "simplifique e foque nos impactos para o QA"
}
```

```bash
curl -X POST http://localhost:8080/api/v1/analyses/uuid-aqui/refine \
  -H "Content-Type: application/json" \
  -d '{"instruction":"simplifique e foque nos impactos para o QA"}'
```

**Response 200:**
```json
{
  "id": "uuid",
  "analysis_id": "uuid-da-analise",
  "instruction": "simplifique...",
  "refined_description": "Texto refinado...",
  "model_used": "gpt-4o-mini",
  "tokens_used": 100,
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Erros:** 400 (instrução vazia), 404 (análise não encontrada), 502, 500

---

### GET /api/v1/analyses/{id}/refinements

Lista todos os refinamentos de uma análise.

```bash
curl http://localhost:8080/api/v1/analyses/uuid-aqui/refinements
```

**Response 200:**
```json
[
  {
    "id": "uuid",
    "analysis_id": "uuid",
    "instruction": "...",
    "refined_description": "...",
    "model_used": "gpt-4o-mini",
    "tokens_used": 100,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

**Erros:** 400, 500

---

## Formato de Erro

Todos os erros retornam:

```json
{
  "error": "error_type",
  "message": "Mensagem amigável"
}
```

| HTTP Status | error_type             | Quando                                        |
|-------------|------------------------|-----------------------------------------------|
| 400         | validation_error       | Campos obrigatórios ausentes ou inválidos      |
| 400         | invalid_request        | JSON malformado                                |
| 404         | not_found              | Recurso não encontrado                         |
| 422         | token_limit_exceeded   | Diff excede limite de tokens do modelo         |
| 429         | rate_limited           | Limite de requisições excedido                 |
| 502         | external_service_error | Bitbucket ou OpenAI indisponível               |
| 504         | timeout                | Timeout na requisição                          |
| 500         | internal_error         | Erro interno do servidor                       |

---

## Rate Limiting

- Limite padrão: **60 requisições por minuto** por IP
- Configurável via variável `RATE_LIMIT_RPM`
- Quando excedido, retorna **429** com header `Retry-After: 60`

---

## Headers

| Header         | Direção   | Descrição                          |
|---------------|-----------|------------------------------------|
| X-Request-ID  | Req/Resp  | ID único para correlação de logs   |
| Retry-After   | Response  | Segundos para retry (em 429)       |
| Content-Type  | Req/Resp  | application/json                   |

---

## Autenticação Bitbucket

O backend conecta ao Bitbucket Cloud usando Basic Auth (email:token). Configure:

```env
BITBUCKET_EMAIL=seu-email@exemplo.com
BITBUCKET_API_TOKEN=seu-app-password
```

O token é um [App Password](https://bitbucket.org/account/settings/app-passwords/) com permissões de leitura em repositórios.
