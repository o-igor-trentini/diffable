# Versionamento de Analises Refinadas

## Objetivo

Implementar versionamento com encadeamento nos refinamentos de analises. Cada refinamento gera uma nova versao (v1, v2, v3...) e por padrao encadeia a partir da versao mais recente. O usuario pode tambem refinar a partir de qualquer versao anterior. O frontend exibe uma timeline visual de todas as versoes.

---

## Etapa 1: Migracao de Banco de Dados

Adicionar colunas `version` e `parent_refinement_id` na tabela `refinements` existente.

### Checklist

- [ ] Criar arquivo `backend/migrations/000005_add_versioning_to_refinements.up.sql`
  - `ALTER TABLE refinements ADD COLUMN version INTEGER NOT NULL DEFAULT 1`
  - `ALTER TABLE refinements ADD COLUMN parent_refinement_id UUID REFERENCES refinements(id) ON DELETE SET NULL`
  - `CREATE UNIQUE INDEX idx_refinements_analysis_version ON refinements(analysis_id, version)`
  - `CREATE INDEX idx_refinements_parent ON refinements(parent_refinement_id)`
- [ ] Criar arquivo `backend/migrations/000005_add_versioning_to_refinements.down.sql`
  - Drop dos indices e colunas adicionados
- [ ] Rodar migracao e validar que refinamentos existentes receberam `version=1` e `parent_refinement_id=NULL`

### Notas
- `DEFAULT 1` atribui automaticamente versao 1 aos registros existentes
- `UNIQUE(analysis_id, version)` previne versoes duplicadas em concorrencia
- `ON DELETE SET NULL` preserva a cadeia se versao intermediaria for removida

---

## Etapa 2: Backend — Domain Model

Atualizar o struct `Refinement` com os novos campos.

### Checklist

- [ ] Editar `backend/internal/domain/model.go`
  - Adicionar campo `Version int` ao struct `Refinement`
  - Adicionar campo `ParentRefinementID *string` ao struct `Refinement`

---

## Etapa 3: Backend — Repository

Adicionar novos metodos e atualizar existentes no repository.

### Checklist

- [ ] Editar `backend/internal/repository/analysis_repository.go`
- [ ] Adicionar metodo `GetRefinementByID(ctx, id) → *Refinement` na interface e implementacao
  - SELECT por ID na tabela refinements incluindo `version` e `parent_refinement_id`
- [ ] Adicionar metodo `GetLatestRefinement(ctx, analysisID) → *Refinement` na interface e implementacao
  - `SELECT ... WHERE analysis_id = $1 ORDER BY version DESC LIMIT 1`
- [ ] Adicionar metodo `GetNextRefinementVersion(ctx, analysisID) → int` na interface e implementacao
  - `SELECT COALESCE(MAX(version), 0) + 1 FROM refinements WHERE analysis_id = $1`
- [ ] Atualizar `CreateRefinement` para incluir `version` e `parent_refinement_id` no INSERT
- [ ] Atualizar `ListRefinements` para incluir `version` e `parent_refinement_id` no SELECT e Scan
- [ ] Verificar que scan de `parent_refinement_id` como `*string` funciona com NULL do PostgreSQL

---

## Etapa 4: Backend — Service

Reescrever a logica de refinamento para suportar encadeamento de versoes.

### Checklist

- [ ] Editar `backend/internal/service/refinement_service.go`
- [ ] Alterar assinatura da interface `RefinementService.Refine` para aceitar `fromRefinementID *string`
- [ ] Implementar logica de encadeamento:
  - Se `fromRefinementID` informado: buscar esse refinamento e usar seu `refined_desc` como base
  - Se nao informado: buscar ultimo refinamento via `GetLatestRefinement`; se existir, usar `refined_desc` dele
  - Se nao existir nenhum refinamento: usar `analysis.GeneratedDesc` (primeiro refinamento)
- [ ] Calcular `nextVersion` via `GetNextRefinementVersion`
- [ ] Preencher campos `Version` e `ParentRefinementID` no refinement antes de salvar

---

## Etapa 5: Backend — DTOs e Handler

Atualizar os DTOs de request/response e o handler para expor os novos campos.

### Checklist

- [ ] Editar `backend/internal/handler/dto/request.go`
  - Adicionar `FromRefinementID *string` ao `RefineRequest`
- [ ] Editar `backend/internal/handler/dto/response.go`
  - Adicionar `OriginalDesc string` ao `RefinementResponse`
  - Adicionar `Version int` ao `RefinementResponse`
  - Adicionar `ParentRefinementID *string` ao `RefinementResponse`
  - Atualizar `RefinementToResponse()` para mapear os novos campos
- [ ] Editar `backend/internal/handler/analysis_handler.go`
  - Atualizar `RefineDescription()` para passar `req.FromRefinementID` ao service
- [ ] Testar endpoints via curl/Postman:
  - `POST /api/v1/analyses/{id}/refine` sem `from_refinement_id` — encadeia automaticamente
  - `POST /api/v1/analyses/{id}/refine` com `from_refinement_id` — refina a partir de versao especifica
  - `GET /api/v1/analyses/{id}/refinements` — retorna `version` e `parent_refinement_id`

---

## Etapa 6: Frontend — Types, Endpoints e Hooks

Atualizar tipagens, chamadas de API e hooks do React Query.

### Checklist

- [ ] Editar `frontend/src/lib/api/types.ts`
  - Adicionar `from_refinement_id?: string` ao `RefineRequest` (ou criar se nao existir)
  - Adicionar `original_description: string` ao `RefinementResponse`
  - Adicionar `version: number` ao `RefinementResponse`
  - Adicionar `parent_refinement_id?: string` ao `RefinementResponse`
- [ ] Editar `frontend/src/lib/api/endpoints.ts`
  - Atualizar `refineDescription()` para aceitar `fromRefinementId?: string` e incluir no body
- [ ] Editar `frontend/src/lib/hooks/useAnalysis.ts`
  - Atualizar `useRefineDescription` para aceitar `fromRefinementId` no objeto de parametros

---

## Etapa 7: Frontend — Componente VersionTimeline

Criar o componente de timeline visual de versoes.

### Checklist

- [ ] Criar `frontend/src/features/refine/VersionTimeline.tsx`
  - Receber props: `analysisDescription`, `analysisCreatedAt`, `refinements[]`, `selectedVersionId`, callbacks
  - Renderizar entrada "v0 — Original" representando a descricao da analise
  - Renderizar entradas "v1", "v2", etc. para cada refinamento
  - Exibir instrucao usada e timestamp relativo em cada entrada
  - Destacar versao selecionada com borda violet
  - Incluir botao "Refinar a partir daqui" em cada versao
  - Seguir design system existente: cores stone/violet, dark mode, rounded-xl, font-mono

---

## Etapa 8: Frontend — Refatorar RefineDescription

Integrar timeline + selecao de versao + formulario de refinamento.

### Checklist

- [ ] Editar `frontend/src/features/refine/RefineDescription.tsx`
  - Adicionar estado `selectedVersionId: string | null` (null = original)
  - Adicionar estado `refiningFromId: string | null` (qual versao sera base do proximo refinamento)
  - Buscar refinamentos via `useRefinements(analysis.id)`
  - Renderizar `VersionTimeline` no topo com a lista de refinamentos
  - Renderizar painel de detalhe mostrando descricao completa da versao selecionada + CopyButton
  - Renderizar `RefineForm` na parte inferior, passando descricao-fonte correta baseada em `refiningFromId`
  - Ao criar refinamento com sucesso: refetch da lista e selecionar nova versao
- [ ] Editar `frontend/src/features/refine/RefineForm.tsx`
  - Adicionar prop opcional `sourceLabel?: string` (ex: "v0 — Original", "v2 — Mais tecnico")
  - Exibir label acima da "Descricao Original" para indicar de qual versao esta refinando

---

## Verificacao Final

- [ ] Migracao roda sem erros e dados existentes sao preservados com `version=1`
- [ ] Primeiro refinamento de uma analise usa a descricao original (v0) como base
- [ ] Segundo refinamento encadeia automaticamente a partir do primeiro (v1)
- [ ] Enviar `from_refinement_id` permite refinar a partir de versao especifica
- [ ] `GET /analyses/{id}/refinements` retorna todas as versoes com `version` e `parent_refinement_id`
- [ ] Timeline no frontend exibe v0 + todas as versoes refinadas
- [ ] Clicar em uma versao na timeline mostra sua descricao completa
- [ ] "Refinar a partir daqui" popula o formulario com a descricao correta
- [ ] Novo refinamento aparece na timeline imediatamente apos criacao
- [ ] Dark mode funciona corretamente no novo componente
