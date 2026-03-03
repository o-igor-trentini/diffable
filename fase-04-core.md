# Fase 4: Funcionalidades Core

**Objetivo:** Conectar Bitbucket + OpenAI em um service de análise, criar endpoints REST para os 3 tipos de análise e construir o frontend completo com formulários e exibição de resultado.

**Pré-requisitos:** Fases 2 e 3 concluídas

---

## Fluxo de Orquestração

```
Request HTTP
    ↓
Handler (valida input, extrai DTO)
    ↓
AnalysisService.Analyze*(ctx, req)
    ↓
┌─ Se hash fornecido: BitbucketClient.GetCommitDiff() / GetPRDiff()
│  Se raw_diff fornecido: usa direto
    ↓
Calcula diff_hash = SHA-256(diff)
    ↓
┌─ Checa cache (in-memory) → hit? retorna
│  Checa DB (GetByDiffHash) → hit? retorna
    ↓
Tokenizer.PreprocessDiff(diff)
    ↓
ModelSelector.SelectModel(tokenCount, type)
    ↓
Generator.Generate(diff, type, metadata)
    ↓
Repository.Create(analysis)
    ↓
Cache.Set(diffHash, description)
    ↓
Retorna AnalysisResponse
```

---

## Checklist Backend

### DTOs

- [ ] Criar `backend/internal/handler/dto/request.go`:
  - `AnalyzeCommitRequest`:
    ```go
    type AnalyzeCommitRequest struct {
        Workspace  string `json:"workspace"`
        RepoSlug   string `json:"repo_slug"`
        CommitHash string `json:"commit_hash"`
        RawDiff    string `json:"raw_diff"`
    }
    ```
    Validação: `commit_hash` (+ workspace + repo_slug) OU `raw_diff` obrigatório
  - `AnalyzeRangeRequest`:
    ```go
    type AnalyzeRangeRequest struct {
        Workspace string `json:"workspace"`
        RepoSlug  string `json:"repo_slug"`
        FromHash  string `json:"from_hash"`
        ToHash    string `json:"to_hash"`
    }
    ```
    Validação: todos campos obrigatórios
  - `AnalyzePRRequest`:
    ```go
    type AnalyzePRRequest struct {
        Workspace string `json:"workspace"`
        RepoSlug  string `json:"repo_slug"`
        PRID      int    `json:"pr_id"`
        RawDiff   string `json:"raw_diff"`
        PRTitle   string `json:"pr_title"`
        PRDesc    string `json:"pr_description"`
    }
    ```
    Validação: `pr_id` (+ workspace + repo_slug) OU (`raw_diff` + `pr_title`) obrigatório

- [ ] Criar `backend/internal/handler/dto/response.go`:
  - `AnalysisResponse`:
    ```go
    type AnalysisResponse struct {
        ID          string `json:"id"`
        Type        string `json:"type"`
        Description string `json:"description"`
        Model       string `json:"model_used"`
        TokensUsed  int    `json:"tokens_used"`
        CreatedAt   string `json:"created_at"`
    }
    ```
  - `ErrorResponse`:
    ```go
    type ErrorResponse struct {
        Error   string `json:"error"`
        Message string `json:"message"`
        Details string `json:"details,omitempty"`
    }
    ```

### Handler

- [ ] Criar `backend/internal/handler/analysis_handler.go`:
  - Struct `AnalysisHandler` com `analysisService`
  - `AnalyzeCommit(w, r)`:
    1. Decode JSON body para `AnalyzeCommitRequest`
    2. Valida input
    3. Chama `service.AnalyzeCommit(ctx, req)`
    4. Retorna 200 + `AnalysisResponse` ou erro mapeado
  - `AnalyzeRange(w, r)`:
    1. Decode JSON body para `AnalyzeRangeRequest`
    2. Valida input (todos campos obrigatórios)
    3. Chama `service.AnalyzeRange(ctx, req)`
    4. Retorna 200 + `AnalysisResponse`
  - `AnalyzePR(w, r)`:
    1. Decode JSON body para `AnalyzePRRequest`
    2. Valida input
    3. Chama `service.AnalyzePR(ctx, req)`
    4. Retorna 200 + `AnalysisResponse`
  - `GetAnalysis(w, r)`:
    1. Extrai `{id}` do path
    2. Chama `service.GetAnalysis(ctx, id)`
    3. Retorna 200 + `AnalysisResponse` ou 404
  - Mapeamento de erros: `ErrNotFound` → 404, `ErrValidation` → 400, `ErrExternalService` → 502
- [ ] Criar `backend/internal/handler/analysis_handler_test.go`:
  - POST commit: JSON válido → 200, campos faltando → 400, serviço erro → 500
  - POST range: JSON válido → 200, hashes inválidos → 400
  - POST PR: JSON válido → 200, PR not found → 404
  - GET analysis: ID existente → 200, ID inexistente → 404

### Service

- [ ] Criar `backend/internal/service/analysis_service.go`:
  - Struct `AnalysisService` com: `bbClient`, `generator`, `repository`, `cache`, `tokenizer`, `modelSelector`
  - `AnalyzeCommit(ctx, req)`:
    1. Se `raw_diff` vazio: `bbClient.GetCommitDiff(workspace, repo, hash)`
    2. Se `raw_diff` vazio: `bbClient.GetCommit(workspace, repo, hash)` para mensagem
    3. `diffHash = SHA-256(diff)`
    4. Checa cache → hit? retorna
    5. Checa DB (`GetByDiffHash`) → hit? retorna
    6. `processedDiff = tokenizer.PreprocessDiff(diff)`
    7. `tokenCount = tokenizer.CountTokens(processedDiff)`
    8. `model = modelSelector.SelectModel(tokenCount, "single_commit")`
    9. `output = generator.Generate(ctx, {Diff, "single_commit", commitMessages})`
    10. `repository.Create(analysis)`
    11. `cache.Set(diffHash, description, ttl)`
    12. Retorna resultado
  - `AnalyzeRange(ctx, req)`:
    1. `bbClient.ListCommitsInRange(workspace, repo, toHash, fromHash)`
    2. Para cada commit ou para o range: `bbClient.GetCommitDiff(workspace, repo, "fromHash..toHash")`
    3. Segue mesmo fluxo de cache → preprocess → generate → save
  - `AnalyzePR(ctx, req)`:
    1. Se `raw_diff` vazio: `bbClient.GetPullRequestDiff(workspace, repo, prID)`
    2. Se `pr_title` vazio: `bbClient.GetPullRequest(workspace, repo, prID)` para título/desc
    3. Segue mesmo fluxo com metadata do PR no input do generator
- [ ] Criar `backend/internal/service/analysis_service_test.go`:
  - Commit com hash: Bitbucket chamado → generator chamado → DB salvo
  - Commit com raw_diff: Bitbucket NÃO chamado → generator chamado
  - Range: commits listados → diff gerado → descrição consolidada
  - PR: diff do PR buscado → metadata do PR usada → descrição gerada
  - Cache hit: generator NÃO chamado, resultado cacheado retornado
  - DB hit (mesmo diff_hash): generator NÃO chamado
  - Bitbucket erro 404: retorna ErrNotFound
  - Generator erro: retorna ErrExternalService

### Repository

- [ ] Criar `backend/internal/repository/analysis_repository.go`:
  - Interface `AnalysisRepository`:
    ```go
    type AnalysisRepository interface {
        Create(ctx context.Context, analysis *domain.Analysis) error
        GetByID(ctx context.Context, id string) (*domain.Analysis, error)
        GetByDiffHash(ctx context.Context, hash string) (*domain.Analysis, error)
        List(ctx context.Context, filter AnalysisFilter, offset, limit int) ([]domain.Analysis, int, error)
    }
    ```
  - Implementação `PostgresAnalysisRepository` usando `pgx/v5`
- [ ] Criar `backend/internal/repository/analysis_repository_test.go`:
  - Integration tests com testcontainers-go (PostgreSQL)
  - Create: insere e retorna sem erro
  - GetByID: encontra análise existente
  - GetByID: retorna ErrNotFound para ID inexistente
  - GetByDiffHash: encontra análise com mesmo hash
  - List: paginação funciona, total correto

### Routes

- [ ] Atualizar `backend/internal/server/routes.go`:
  - `POST /api/v1/analyses/commit` → `analysisHandler.AnalyzeCommit`
  - `POST /api/v1/analyses/range` → `analysisHandler.AnalyzeRange`
  - `POST /api/v1/analyses/pr` → `analysisHandler.AnalyzePR`
  - `GET /api/v1/analyses/{id}` → `analysisHandler.GetAnalysis`

### Dependency Injection

- [ ] Atualizar `backend/cmd/server/main.go`:
  - Criar BitbucketClient com config
  - Criar OpenAIGenerator com config
  - Criar InMemoryCache com TTL
  - Criar PostgresAnalysisRepository com pool de conexões
  - Criar AnalysisService com todas dependências
  - Criar AnalysisHandler com service
  - Registrar rotas

---

## Checklist Frontend

### API Types e Endpoints

- [ ] Criar `frontend/src/lib/api/types.ts`:
  ```typescript
  interface AnalyzeCommitRequest {
    workspace?: string;
    repo_slug?: string;
    commit_hash?: string;
    raw_diff?: string;
  }

  interface AnalyzeRangeRequest {
    workspace: string;
    repo_slug: string;
    from_hash: string;
    to_hash: string;
  }

  interface AnalyzePRRequest {
    workspace?: string;
    repo_slug?: string;
    pr_id?: number;
    raw_diff?: string;
    pr_title?: string;
    pr_description?: string;
  }

  interface AnalysisResponse {
    id: string;
    type: string;
    description: string;
    model_used: string;
    tokens_used: number;
    created_at: string;
  }

  interface ErrorResponse {
    error: string;
    message: string;
    details?: string;
  }
  ```
- [ ] Criar `frontend/src/lib/api/endpoints.ts`:
  - `analyzeCommit(req: AnalyzeCommitRequest): Promise<AnalysisResponse>`
  - `analyzeRange(req: AnalyzeRangeRequest): Promise<AnalysisResponse>`
  - `analyzePR(req: AnalyzePRRequest): Promise<AnalysisResponse>`
  - `getAnalysis(id: string): Promise<AnalysisResponse>`

### Hooks

- [ ] Criar `frontend/src/lib/hooks/useAnalysis.ts`:
  - `useAnalyzeCommit()` — `useMutation` wrapping `analyzeCommit`
  - `useAnalyzeRange()` — `useMutation` wrapping `analyzeRange`
  - `useAnalyzePR()` — `useMutation` wrapping `analyzePR`
  - Cada um expõe: `mutate`, `isPending`, `isError`, `error`, `data`
- [ ] Criar `frontend/src/lib/hooks/useClipboard.ts`:
  - `useClipboard()` retorna `{ copy: (text: string) => void, copied: boolean }`
  - Usa `navigator.clipboard.writeText`
  - `copied` volta para false após 2 segundos

### Componentes — Commit

- [ ] Criar `frontend/src/features/commit/CommitAnalysis.tsx`:
  - Container: `CommitForm` + `ResultDisplay`
  - Gerencia estado da mutation
  - Passa resultado para ResultDisplay
- [ ] Criar `frontend/src/features/commit/CommitForm.tsx`:
  - Campos: workspace (text), repo slug (text), commit hash (text)
  - OU: raw diff (textarea grande)
  - Separador visual "OU cole o diff manualmente"
  - Botão "Gerar Descrição" com ícone Zap
  - Validação: hash+workspace+repo OU raw_diff obrigatório
  - Loading state no botão durante mutation

### Componentes — Range

- [ ] Criar `frontend/src/features/range/RangeAnalysis.tsx`:
  - Container: `RangeForm` + `ResultDisplay`
- [ ] Criar `frontend/src/features/range/RangeForm.tsx`:
  - Campos: workspace, repo slug, hash inicial (from), hash final (to)
  - Todos obrigatórios
  - Botão "Gerar Descrição"

### Componentes — Pull Request

- [ ] Criar `frontend/src/features/pull-request/PrAnalysis.tsx`:
  - Container: `PrForm` + `ResultDisplay`
- [ ] Criar `frontend/src/features/pull-request/PrForm.tsx`:
  - Campos: workspace, repo slug, PR ID (number)
  - OU: raw diff (textarea) + PR título (text) + PR descrição (textarea)
  - Botão "Gerar Descrição"

### Componentes Shared

- [ ] Criar `frontend/src/features/shared/ResultDisplay.tsx`:
  - Card estilizado mostrando descrição gerada
  - Renderiza markdown da descrição
  - Botão "Copiar" (ícone Copy) — usa useClipboard
  - Botão "Refinar" (ícone RefreshCw) — navegará para tab Refine na Fase 5
  - Info: modelo usado, tokens consumidos
  - Animação de entrada suave
- [ ] Criar `frontend/src/features/shared/CopyButton.tsx`:
  - Ícone Copy → muda para Check por 2s após copiar
  - Tooltip "Copiado!" temporário

### App

- [ ] Atualizar `frontend/src/App.tsx`:
  - Renderiza componente correto baseado na tab ativa:
    - Tab Commit → `CommitAnalysis`
    - Tab Range → `RangeAnalysis`
    - Tab PR → `PrAnalysis`
    - Tab Refinar → placeholder (Fase 5)

### Testes Frontend

- [ ] Criar `frontend/src/test/mocks/handlers.ts` — MSW handlers mockando API
- [ ] Criar `frontend/src/test/mocks/server.ts` — MSW server setup
- [ ] Teste: CommitForm renderiza campos, valida required, submete corretamente
- [ ] Teste: RangeForm renderiza campos, valida hashes obrigatórios
- [ ] Teste: PrForm renderiza campos, valida PR ID é número
- [ ] Teste: ResultDisplay mostra descrição, botão copiar funciona
- [ ] Teste: useClipboard copia texto, estado `copied` alterna

---

## Testes desta Fase

| Teste | Tipo | Validação |
|-------|------|-----------|
| `analysis_handler_test.go` — POST commit | Unit | JSON válido → 200, campos faltando → 400 |
| `analysis_handler_test.go` — POST range | Unit | JSON válido → 200, hashes inválidos → 400 |
| `analysis_handler_test.go` — POST PR | Unit | JSON válido → 200, PR not found → 404 |
| `analysis_handler_test.go` — GET analysis | Unit | ID existente → 200, inexistente → 404 |
| `analysis_service_test.go` — commit hash | Unit | Bitbucket + generator + DB chamados |
| `analysis_service_test.go` — commit raw_diff | Unit | Bitbucket NÃO chamado |
| `analysis_service_test.go` — cache hit | Unit | Generator NÃO chamado |
| `analysis_service_test.go` — range | Unit | Commits listados, descrição consolidada |
| `analysis_service_test.go` — PR | Unit | Diff + metadata do PR usados |
| `analysis_repository_test.go` — CRUD | Integration | Create, GetByID, GetByDiffHash, List com testcontainers |
| CommitForm.test.tsx | Unit (Vitest) | Renderiza, valida, submete |
| RangeForm.test.tsx | Unit (Vitest) | Renderiza, valida |
| PrForm.test.tsx | Unit (Vitest) | Renderiza, valida |
| ResultDisplay.test.tsx | Unit (Vitest) | Mostra descrição, copy funciona |
| useClipboard.test.ts | Unit (Vitest) | Copia texto, estado alterna |
| E2E: commit analysis | E2E (Playwright) | Preenche form → Gerar → loading → resultado → Copy |

---

## Critérios de Aceite

- [ ] `POST /api/v1/analyses/commit` com `{"raw_diff": "..."}` retorna descrição gerada
- [ ] `POST /api/v1/analyses/commit` com workspace/repo/hash busca do Bitbucket e retorna descrição
- [ ] `POST /api/v1/analyses/range` retorna descrição consolidada do range
- [ ] `POST /api/v1/analyses/pr` retorna descrição do PR
- [ ] `GET /api/v1/analyses/{id}` retorna análise salva
- [ ] Frontend Tab Commit: preenche form → Gerar → vê resultado → copia
- [ ] Frontend Tab Range: preenche form → Gerar → vê resultado
- [ ] Frontend Tab PR: preenche form → Gerar → vê resultado
- [ ] Loading spinner aparece durante geração
- [ ] Alerta de erro aparece em falha com mensagem significativa
- [ ] Todos testes unitários e de integração passam
