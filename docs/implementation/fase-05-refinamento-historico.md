# Fase 5: Refinamento e Histórico

**Objetivo:** Implementar a funcionalidade de refinamento de descrições (usuário adapta descrição já gerada com instruções) e o painel de histórico de análises anteriores.

**Pré-requisitos:** Fase 4 concluída

---

## Checklist Backend

### Migration

- [ ] Criar `backend/migrations/000002_create_refinements_table.up.sql`:
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
- [ ] Criar `backend/migrations/000002_create_refinements_table.down.sql`:
  ```sql
  DROP TABLE IF EXISTS refinements;
  ```

### Repository

- [ ] Atualizar `backend/internal/repository/analysis_repository.go` — adicionar métodos:
  - `CreateRefinement(ctx, refinement *domain.Refinement) error`
  - `ListRefinements(ctx, analysisID string) ([]domain.Refinement, error)`
  - `ListAnalyses(ctx, filter AnalysisFilter, offset, limit int) ([]domain.Analysis, int, error)`
    - `AnalysisFilter` com campo opcional `Type` (analysis_type)
    - Retorna lista + total count para paginação
- [ ] Atualizar testes do repository:
  - CreateRefinement: insere com analysis_id válido
  - CreateRefinement: falha com analysis_id inválido (FK constraint)
  - ListRefinements: retorna refinamentos ordenados por created_at
  - ListAnalyses: paginação correta, filtro por tipo funciona

### Service — Refinamento

- [ ] Criar `backend/internal/service/refinement_service.go`:
  - Struct `RefinementService` com: `generator`, `repository`
  - `Refine(ctx, analysisID string, instruction string) (*domain.Refinement, error)`:
    1. Busca análise por ID (`repository.GetByID`)
    2. Se não encontrada → `ErrNotFound`
    3. Valida instrução não vazia
    4. Chama `generator.Refine(ctx, {OriginalDescription: analysis.GeneratedDesc, Instruction: instruction})`
    5. Cria `Refinement` com resultado
    6. Salva no banco (`repository.CreateRefinement`)
    7. Retorna refinement
- [ ] Criar `backend/internal/service/refinement_service_test.go`:
  - Refinamento válido: análise encontrada → generator chamado → salvo no DB
  - Análise não encontrada → retorna ErrNotFound
  - Instrução vazia → retorna ErrValidation
  - Generator erro → propaga erro

### Service — Histórico

- [ ] Criar `backend/internal/service/history_service.go`:
  - Struct `HistoryService` com: `repository`
  - `ListAnalyses(ctx, typeFilter string, page, pageSize int) ([]domain.Analysis, int, error)`:
    - Calcula offset = (page - 1) * pageSize
    - Chama `repository.ListAnalyses(ctx, filter, offset, pageSize)`
  - `GetRefinements(ctx, analysisID string) ([]domain.Refinement, error)`:
    - Chama `repository.ListRefinements(ctx, analysisID)`
- [ ] Criar `backend/internal/service/history_service_test.go`:
  - ListAnalyses: sem filtro retorna todos
  - ListAnalyses: filtro por tipo funciona
  - ListAnalyses: paginação correta
  - GetRefinements: retorna lista de refinamentos

### Handler — Refinamento

- [ ] Atualizar `backend/internal/handler/analysis_handler.go` — adicionar:
  - `RefineDescription(w, r)`:
    1. Extrai `{id}` do path
    2. Decode JSON body: `{"instruction": "..."}`
    3. Valida instrução não vazia
    4. Chama `refinementService.Refine(ctx, id, instruction)`
    5. Retorna 200 + `RefinementResponse` ou erro

### Handler — Histórico

- [ ] Criar `backend/internal/handler/history_handler.go`:
  - Struct `HistoryHandler` com: `historyService`
  - `ListAnalyses(w, r)`:
    1. Query params: `type` (opcional), `page` (default 1), `page_size` (default 20, max 100)
    2. Chama `historyService.ListAnalyses(ctx, type, page, pageSize)`
    3. Retorna 200 + `PaginatedAnalysesResponse`
  - `GetRefinements(w, r)`:
    1. Extrai `{id}` do path
    2. Chama `historyService.GetRefinements(ctx, id)`
    3. Retorna 200 + lista de `RefinementResponse`
- [ ] Criar `backend/internal/handler/history_handler_test.go`:
  - ListAnalyses: sem filtro → 200 com lista
  - ListAnalyses: filtro `?type=pull_request` → 200 com filtrado
  - ListAnalyses: paginação `?page=2&page_size=10` → offset correto
  - GetRefinements: ID existente → 200 com lista
  - GetRefinements: ID inexistente → 404

### DTOs

- [ ] Atualizar `backend/internal/handler/dto/request.go`:
  - `RefineRequest`:
    ```go
    type RefineRequest struct {
        Instruction string `json:"instruction"`
    }
    ```
- [ ] Atualizar `backend/internal/handler/dto/response.go`:
  - `RefinementResponse`:
    ```go
    type RefinementResponse struct {
        ID          string `json:"id"`
        AnalysisID  string `json:"analysis_id"`
        Instruction string `json:"instruction"`
        RefinedDesc string `json:"refined_description"`
        Model       string `json:"model_used"`
        TokensUsed  int    `json:"tokens_used"`
        CreatedAt   string `json:"created_at"`
    }
    ```
  - `PaginatedAnalysesResponse`:
    ```go
    type PaginatedAnalysesResponse struct {
        Data     []AnalysisResponse `json:"data"`
        Total    int                `json:"total"`
        Page     int                `json:"page"`
        PageSize int                `json:"page_size"`
    }
    ```

### Routes

- [ ] Atualizar `backend/internal/server/routes.go`:
  - `POST /api/v1/analyses/{id}/refine` → `analysisHandler.RefineDescription`
  - `GET /api/v1/analyses` → `historyHandler.ListAnalyses`
  - `GET /api/v1/analyses/{id}/refinements` → `historyHandler.GetRefinements`

### Dependency Injection

- [ ] Atualizar `backend/cmd/server/main.go`:
  - Criar `RefinementService` com generator + repository
  - Criar `HistoryService` com repository
  - Criar `HistoryHandler` com historyService
  - Registrar novas rotas

---

## Checklist Frontend

### API e Hooks

- [ ] Atualizar `frontend/src/lib/api/types.ts`:
  ```typescript
  interface RefineRequest {
    instruction: string;
  }

  interface RefinementResponse {
    id: string;
    analysis_id: string;
    instruction: string;
    refined_description: string;
    model_used: string;
    tokens_used: number;
    created_at: string;
  }

  interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
  }

  interface HistoryFilter {
    type?: string;
    page?: number;
    page_size?: number;
  }
  ```
- [ ] Atualizar `frontend/src/lib/api/endpoints.ts`:
  - `refineDescription(id: string, instruction: string): Promise<RefinementResponse>`
  - `listAnalyses(filter?: HistoryFilter): Promise<PaginatedResponse<AnalysisResponse>>`
  - `getRefinements(id: string): Promise<RefinementResponse[]>`
- [ ] Atualizar `frontend/src/lib/hooks/useAnalysis.ts`:
  - `useRefineDescription()` — `useMutation` wrapping `refineDescription`
- [ ] Criar `frontend/src/lib/hooks/useHistory.ts`:
  - `useAnalysesList(filter?: HistoryFilter)` — `useQuery`
  - `useRefinements(analysisID: string)` — `useQuery`

### Componentes — Refinar

- [ ] Criar `frontend/src/features/refine/RefineDescription.tsx`:
  - Container component
  - Aceita `initialDescription` via props (quando vindo do botão "Refinar")
  - Pode ser usado standalone (cole a descrição manualmente)
  - Mostra: original + form instrução + resultado refinado
- [ ] Criar `frontend/src/features/refine/RefineForm.tsx`:
  - Campos:
    - Descrição original (textarea, pré-preenchido se vindo de resultado)
    - Instrução (textarea com placeholder: "ex: simplifique, mais técnico, mais resumido")
  - Botão "Refinar" com ícone RefreshCw
  - Loading state durante mutation
  - Resultado mostra versão refinada com opção de refinar novamente

### Componentes — Histórico

- [ ] Criar `frontend/src/features/history/HistoryPanel.tsx`:
  - Painel lateral ou seção abaixo do resultado
  - Toggle com ícone History (Lucide)
  - Lista de análises anteriores com scroll
  - Filtro por tipo: dropdown (Todos, Commit, Range, PR)
  - Paginação (botões anterior/próximo)
  - Loading state
  - Estado vazio: "Nenhuma análise encontrada"
- [ ] Criar `frontend/src/features/history/HistoryItem.tsx`:
  - Ícone por tipo: GitCommit (commit), GitBranch (range), GitPullRequest (PR)
  - Descrição truncada (primeiros 100 caracteres)
  - Tempo relativo (ex: "há 2 horas")
  - Click carrega análise completa no ResultDisplay

### App — Integrações

- [ ] Atualizar `frontend/src/App.tsx`:
  - Tab Refinar agora renderiza `RefineDescription`
  - Adicionar toggle do painel de histórico (ícone History na toolbar)
  - Conectar botão "Refinar" do `ResultDisplay`:
    - Ao clicar, muda para tab Refinar
    - Passa descrição atual como `initialDescription`
  - State management: guardar análise atual para transferência entre tabs

### Testes Frontend

- [ ] Teste: RefineDescription renderiza form, submete, mostra resultado
- [ ] Teste: RefineForm pré-preenche descrição original, valida instrução não vazia
- [ ] Teste: HistoryPanel renderiza lista, filtro muda resultados, paginação funciona
- [ ] Teste: HistoryItem mostra ícone correto por tipo, trunca descrição, tempo relativo

---

## Testes desta Fase

| Teste | Tipo | Validação |
|-------|------|-----------|
| `refinement_service_test.go` | Unit | Refinamento válido, instrução vazia rejeitada, erro propagado |
| `history_service_test.go` | Unit | Lista com filtro, paginação, busca por ID |
| `analysis_handler_test.go` — refine | Unit | POST refine válido → 200, instrução faltando → 400, análise not found → 404 |
| `history_handler_test.go` — list | Unit | Query params parseados, filtro por tipo |
| `history_handler_test.go` — refinements | Unit | ID existente → 200, inexistente → 404 |
| Repository — refinements | Integration | Create, list com testcontainers |
| Repository — list analyses | Integration | Paginação, filtro |
| RefineDescription.test.tsx | Unit (Vitest) | Pré-preenche, submete, mostra resultado |
| HistoryPanel.test.tsx | Unit (Vitest) | Renderiza lista, filtra, pagina |
| E2E: refine flow | E2E (Playwright) | Gerar descrição → Refinar → instrução → versão refinada |
| E2E: history | E2E (Playwright) | Gerar 2 análises → histórico → ver ambas → clicar uma → detalhe |

---

## Critérios de Aceite

- [ ] `POST /api/v1/analyses/{id}/refine` com `{"instruction": "simplifique"}` retorna descrição refinada
- [ ] `GET /api/v1/analyses` retorna lista paginada de análises
- [ ] `GET /api/v1/analyses?type=pull_request` filtra corretamente
- [ ] `GET /api/v1/analyses/{id}/refinements` retorna todos refinamentos da análise
- [ ] Tab Refinar funciona standalone (colar descrição + instrução)
- [ ] Botão "Refinar" no ResultDisplay navega para tab Refinar com dados pré-preenchidos
- [ ] Painel de histórico mostra análises com ícones de tipo e timestamps
- [ ] Clicar em item do histórico carrega análise completa
- [ ] Migration 000002 executa com sucesso, tabela `refinements` existe
- [ ] Todos testes passam
