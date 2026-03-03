# DIFFABLE: Plano de Implementação

Plataforma fullstack (React + Go) para gerar descrições automáticas de commits e PRs do Bitbucket Cloud, usando OpenAI para geração de texto focado em QA/PO e JIRA.

## Problema

Preenchimento manual de cards JIRA a partir de PRs/commits é lento e gera descrições inconsistentes entre Dev, QA e PO.

## Solução

Plataforma que busca diffs do Bitbucket Cloud via API REST, envia para OpenAI e gera descrições prontas para JIRA em linguagem não-técnica.

## Decisões Técnicas

| Decisão | Escolha |
|---------|---------|
| Bitbucket | Cloud (api.bitbucket.org/2.0/) |
| Autenticação Bitbucket | API Token (email:token Basic Auth) |
| LLM | OpenAI (gpt-4o-mini default, gpt-4o para complexos) |
| Infra | Docker Compose (PostgreSQL + Go + React) |
| Testes | Obrigatórios em cada fase |

## Fases

| Fase | Descrição | Dependências | Arquivo |
|------|-----------|--------------|---------|
| 1 | Fundação — Scaffolding, Docker, DB, Health | — | [fase-01-fundacao.md](fase-01-fundacao.md) |
| 2 | Bitbucket — Cliente API completo | Fase 1 | [fase-02-bitbucket.md](fase-02-bitbucket.md) |
| 3 | OpenAI — Gerador de descrições | Fase 1 | [fase-03-openai.md](fase-03-openai.md) |
| 4 | Core — Endpoints + Frontend funcional | Fases 2, 3 | [fase-04-core.md](fase-04-core.md) |
| 5 | Refinamento e Histórico | Fase 4 | [fase-05-refinamento-historico.md](fase-05-refinamento-historico.md) |
| 6 | Polish e Produção | Fase 5 | [fase-06-polish-producao.md](fase-06-polish-producao.md) |
| 7 | Melhorias Futuras (opcional) | Fase 6 | [fase-07-melhorias-futuras.md](fase-07-melhorias-futuras.md) |

## Sequenciamento

```
Fase 1 (Fundação)
  ├── Fase 2 (Bitbucket) ──┐
  └── Fase 3 (OpenAI) ─────┤  (podem ser paralelas)
                            ▼
                    Fase 4 (Core Features)
                            ▼
                    Fase 5 (Refine + Histórico)
                            ▼
                    Fase 6 (Polish + Produção)
                            ▼
                    Fase 7 (Opcional)
```

> Fases 2 e 3 podem ser desenvolvidas em paralelo após Fase 1.

## Estrutura do Projeto

### Backend (Go)

```
backend/
  cmd/server/main.go
  internal/
    config/config.go
    server/{server,routes}.go
    handler/{analysis,history,health}_handler.go
    handler/dto/{request,response}.go
    service/{analysis,refinement,history}_service.go
    repository/analysis_repository.go
    domain/{model,errors}.go
    bitbucket/{client,types,pagination,ratelimit,errors}.go
    openai/{generator,prompts,tokenizer,model_selector,retry}.go
    cache/cache.go
    middleware/{cors,logging,ratelimit,requestid}.go
  migrations/
  Dockerfile, Makefile, .env.example
```

### Frontend (React)

```
frontend/src/
  main.tsx, App.tsx
  lib/api/{client,endpoints,types}.ts
  lib/hooks/{useAnalysis,useHistory,useClipboard}.ts
  features/
    commit/{CommitAnalysis,CommitForm}.tsx
    range/{RangeAnalysis,RangeForm}.tsx
    pull-request/{PrAnalysis,PrForm}.tsx
    refine/{RefineDescription,RefineForm}.tsx
    history/{HistoryPanel,HistoryItem}.tsx
    shared/{ResultDisplay,TabNavigation,Button,TextArea,
            CopyButton,LoadingSpinner,ErrorAlert,ErrorBoundary}.tsx
  test/{setup.ts, mocks/}
  Dockerfile, .env.example
```

### Raiz

```
docker-compose.yml, docker-compose.dev.yml, Makefile
docs/{architecture,api,phases,bitbucket-integration,openai-integration}.md
```

## Banco de Dados

### Migration 000001: analyses

```sql
CREATE TYPE analysis_type AS ENUM ('single_commit','commit_range','pull_request');
CREATE TYPE analysis_status AS ENUM ('pending','processing','completed','failed');

CREATE TABLE analyses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_type   analysis_type NOT NULL,
    status          analysis_status NOT NULL DEFAULT 'pending',
    workspace       VARCHAR(255),
    repo_slug       VARCHAR(255),
    commit_hash     VARCHAR(40),
    from_hash       VARCHAR(40),
    to_hash         VARCHAR(40),
    pr_id           INTEGER,
    raw_diff        TEXT,
    diff_hash       VARCHAR(64),
    generated_desc  TEXT,
    model_used      VARCHAR(50),
    tokens_used     INTEGER,
    error_message   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analyses_diff_hash ON analyses(diff_hash) WHERE diff_hash IS NOT NULL;
CREATE INDEX idx_analyses_created_at ON analyses(created_at DESC);
CREATE INDEX idx_analyses_type ON analyses(analysis_type);
```

### Migration 000002: refinements

```sql
CREATE TABLE refinements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_id     UUID REFERENCES analyses(id) ON DELETE CASCADE,
    instruction     TEXT NOT NULL,
    original_desc   TEXT NOT NULL,
    refined_desc    TEXT,
    model_used      VARCHAR(50),
    tokens_used     INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refinements_analysis_id ON refinements(analysis_id);
```

## API Endpoints

| Método | Path | Fase | Descrição |
|--------|------|------|-----------|
| GET | `/healthz` | 1 | Liveness probe |
| GET | `/readyz` | 1 | Readiness probe (DB check) |
| POST | `/api/v1/analyses/commit` | 4 | Análise de commit único |
| POST | `/api/v1/analyses/range` | 4 | Análise de range de commits |
| POST | `/api/v1/analyses/pr` | 4 | Análise de PR |
| GET | `/api/v1/analyses` | 5 | Listar histórico (paginado, filtrável) |
| GET | `/api/v1/analyses/{id}` | 4 | Buscar análise por ID |
| POST | `/api/v1/analyses/{id}/refine` | 5 | Refinar descrição |
| GET | `/api/v1/analyses/{id}/refinements` | 5 | Listar refinamentos |

## Variáveis de Ambiente

### Backend

```env
PORT=8080
DATABASE_URL=postgres://user:pass@db:5432/bbgendesc?sslmode=disable
FRONTEND_URL=http://localhost:3000
LOG_LEVEL=info

BITBUCKET_BASE_URL=https://api.bitbucket.org/2.0
BITBUCKET_EMAIL=user@example.com
BITBUCKET_API_TOKEN=your-token
BITBUCKET_TIMEOUT=30s

OPENAI_API_KEY=sk-...
OPENAI_DEFAULT_MODEL=gpt-4o-mini
OPENAI_COMPLEX_MODEL=gpt-4o
OPENAI_MAX_TOKENS=1024
OPENAI_TEMPERATURE=0.3
OPENAI_TOKEN_THRESHOLD=4000

CACHE_TTL=24h
RATE_LIMIT_RPM=60
DB_MAX_CONNS=25
DB_MIN_CONNS=5
SHUTDOWN_TIMEOUT=30s
```

### Frontend

```env
VITE_API_URL=/api/v1
```

## Dependências Principais

### Backend

| Pacote | Versão | Uso |
|--------|--------|-----|
| go-chi/chi/v5 | v5.2+ | HTTP router |
| jackc/pgx/v5 | v5.7+ | PostgreSQL driver |
| golang-migrate/migrate/v4 | v4.18+ | DB migrations |
| sashabaranov/go-openai | latest | OpenAI API client |
| pkoukk/tiktoken-go | latest | Token counting |
| stretchr/testify | v1.9+ | Test assertions |
| testcontainers/testcontainers-go | latest | Integration test DB |

### Frontend

| Pacote | Versão | Uso |
|--------|--------|-----|
| react | ^18.3 | UI framework |
| vite | ^5.4+ | Build tool |
| tailwindcss | ^3.4+ | CSS utility framework |
| @tanstack/react-query | ^5.60+ | Server state management |
| lucide-react | latest | Icon library |
| axios | ^1.7+ | HTTP client |
| vitest | ^2.1+ | Unit test runner |
| @testing-library/react | ^16+ | Component testing |
| playwright | ^1.48+ | E2E testing |
| msw | ^2.4+ | API mocking |

### Infra

| Componente | Versão |
|------------|--------|
| PostgreSQL | 16-alpine |
| Docker Compose | v2 |
| Node | 20-alpine |
| nginx | 1.27-alpine |

## Verificação End-to-End

Para validar a implementação completa:

1. `docker compose up` — todos serviços sobem
2. `curl localhost:8080/healthz` → `{"status":"ok"}`
3. Frontend em `localhost:3000` — 4 tabs navegáveis
4. Tab Commit: preencher workspace/repo/hash → Gerar → descrição aparece
5. Tab PR: preencher workspace/repo/PR ID → Gerar → descrição aparece
6. Copiar descrição → funciona
7. Refinar → escrever instrução → versão refinada aparece
8. Histórico → análises anteriores listadas
9. `make test-all` — todos testes passam
