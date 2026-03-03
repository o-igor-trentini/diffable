# Fase 1: Fundação

**Objetivo:** Estabelecer a estrutura do projeto, Docker Compose, banco de dados com migrations, servidor Go com health checks e app React com Tailwind e navegação por tabs.

**Pré-requisitos:** Nenhum (primeira fase)

---

## Checklist Backend

### Estrutura e Configuração

- [ ] Criar `backend/go.mod` com módulo `github.com/igor-trentini/diffable/backend`
- [ ] Criar `backend/cmd/server/main.go` — entry point: carrega config, conecta DB, roda migrations, inicia servidor HTTP
- [ ] Criar `backend/internal/config/config.go` — carrega variáveis de ambiente com defaults:
  - `PORT` (default: 8080)
  - `DATABASE_URL`
  - `FRONTEND_URL` (default: http://localhost:3000)
  - `LOG_LEVEL` (default: info)
- [ ] Criar `backend/internal/config/config_test.go` — testa carregamento com defaults e overrides

### Servidor HTTP

- [ ] Criar `backend/internal/server/server.go` — configuração do chi router com middlewares:
  - `middleware.RequestID`
  - `middleware.RealIP`
  - `middleware.Logger`
  - `middleware.Recoverer`
  - CORS (permite frontend origin)
  - Timeout (30s)
- [ ] Criar `backend/internal/server/routes.go` — registra rotas de health + placeholders para `/api/v1/`

### Handlers

- [ ] Criar `backend/internal/handler/health_handler.go`:
  - `GET /healthz` → retorna `{"status":"ok"}`
  - `GET /readyz` → verifica conexão com DB, retorna `{"status":"ok"}` ou `{"status":"error"}`
- [ ] Criar `backend/internal/handler/health_handler_test.go` — testa ambos endpoints

### Domain

- [ ] Criar `backend/internal/domain/model.go`:
  - `AnalysisType` enum: `single_commit`, `commit_range`, `pull_request`
  - `AnalysisStatus` enum: `pending`, `processing`, `completed`, `failed`
  - Struct `Analysis` com todos campos da migration
  - Struct `Refinement` com todos campos da migration
- [ ] Criar `backend/internal/domain/errors.go`:
  - `ErrNotFound`
  - `ErrValidation`
  - `ErrExternalService`

### Middleware

- [ ] Criar `backend/internal/middleware/cors.go` — CORS permitindo `FRONTEND_URL`
- [ ] Criar `backend/internal/middleware/logging.go` — structured logging com `slog` (JSON)

### Migrations

- [ ] Criar `backend/migrations/000001_create_analyses_table.up.sql` — schema completo (ver PLANO.md)
- [ ] Criar `backend/migrations/000001_create_analyses_table.down.sql` — `DROP TABLE`, `DROP TYPE`

### Build e Deploy

- [ ] Criar `backend/Dockerfile` — multi-stage:
  - Stage 1: `golang:1.22-alpine` para build
  - Stage 2: `gcr.io/distroless/static-debian12` para runtime
- [ ] Criar `backend/Makefile` com targets: `build`, `run`, `test`, `lint`, `migrate-up`, `migrate-down`
- [ ] Criar `backend/.env.example` — todas variáveis documentadas

---

## Checklist Frontend

### Estrutura e Configuração

- [ ] Criar projeto com `npm create vite@latest frontend -- --template react-ts`
- [ ] Instalar dependências: `tailwindcss`, `postcss`, `autoprefixer`, `@tanstack/react-query`, `lucide-react`, `axios`
- [ ] Criar `frontend/tailwind.config.ts` — content paths, cores customizadas
- [ ] Criar `frontend/postcss.config.js` — tailwind + autoprefixer
- [ ] Criar `frontend/vite.config.ts` — proxy `/api` para `http://localhost:8080`
- [ ] Configurar `frontend/tsconfig.json` — strict mode, path alias `@/` = `src/`

### App Shell

- [ ] Criar `frontend/src/main.tsx` — ReactDOM.createRoot, QueryClientProvider
- [ ] Criar `frontend/src/App.tsx` — layout com TabNavigation, renderiza tab ativa
- [ ] Criar `frontend/src/assets/styles/index.css` — directives do Tailwind

### Componentes Shared

- [ ] Criar `frontend/src/features/shared/TabNavigation.tsx`:
  - 4 tabs: Commit (GitCommit), Range (GitBranch), PR (GitPullRequest), Refinar (RefreshCw)
  - Ícones Lucide para cada tab
  - Estado ativo com estilo visual
- [ ] Criar `frontend/src/features/shared/Button.tsx` — variantes: primary, secondary, ghost + loading state
- [ ] Criar `frontend/src/features/shared/TextArea.tsx` — textarea estilizado com label e erro
- [ ] Criar `frontend/src/features/shared/LoadingSpinner.tsx` — spinner animado
- [ ] Criar `frontend/src/features/shared/ErrorAlert.tsx` — alerta vermelho com ícone AlertCircle

### API Client

- [ ] Criar `frontend/src/lib/api/client.ts` — Axios instance com `baseURL: '/api/v1'`, interceptor de erro

### Build e Deploy

- [ ] Criar `frontend/Dockerfile` — multi-stage:
  - Stage 1: `node:20-alpine` para build
  - Stage 2: `nginx:1.27-alpine` para serve
- [ ] Criar `frontend/.env.example` — `VITE_API_URL`

### Testes

- [ ] Configurar Vitest: `vitest.config.ts`, `src/test/setup.ts`
- [ ] Instalar `@testing-library/react`, `@testing-library/jest-dom`, `jsdom`
- [ ] Criar teste: TabNavigation renderiza 4 tabs e alterna estado ativo

---

## Checklist Raiz

### Docker Compose

- [ ] Criar `docker-compose.yml`:
  - `db`: postgres:16-alpine, volume persistente, healthcheck (`pg_isready`)
  - `backend`: build ./backend, depends_on db (condition: healthy), env_file
  - `frontend`: build ./frontend, depends_on backend, porta 3000
- [ ] Criar `docker-compose.dev.yml`:
  - `backend`: monta volume, usa `air` para hot reload
  - `frontend`: monta volume, usa `vite dev`
- [ ] Criar `Makefile` com targets: `up`, `down`, `dev`, `test-all`, `logs`

### Documentação

- [ ] Criar `README.md` — overview do projeto, quick-start com Docker Compose
- [ ] Criar `docs/architecture.md` — decisões: chi, pgx, TanStack Query, clean arch
- [ ] Criar `docs/phases.md` — link para cada arquivo de fase

### Outros

- [ ] Criar `.gitignore` — Go binaries, node_modules, .env, dist, tmp
- [ ] Criar `.editorconfig` — 2-space para TS/JSON, tab para Go, LF

---

## Testes desta Fase

| Teste | Tipo | Validação |
|-------|------|-----------|
| `config_test.go` | Unit | Config carrega defaults e overrides de env |
| `health_handler_test.go` | Unit | `/healthz` → 200, `/readyz` → 200 com DB up |
| TabNavigation.test.tsx | Unit (Vitest) | Renderiza 4 tabs, click alterna ativo |
| Docker integration | Integration | `docker compose up` + `curl /healthz` → 200 |

---

## Critérios de Aceite

- [ ] `docker compose up` inicia os 3 serviços sem erro
- [ ] `curl http://localhost:8080/healthz` retorna `{"status":"ok"}`
- [ ] `curl http://localhost:8080/readyz` retorna `{"status":"ok"}`
- [ ] Frontend carrega em `http://localhost:3000` com 4 tabs visíveis
- [ ] Clicar em cada tab alterna o painel visível (placeholders por enquanto)
- [ ] `make test-all` passa todos os testes unitários
- [ ] Migration 000001 executa com sucesso, tabela `analyses` existe no PostgreSQL
- [ ] Todas env vars documentadas em `.env.example`
