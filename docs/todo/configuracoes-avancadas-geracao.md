# Configuracoes Avancadas de Geracao via Frontend

## Contexto

Atualmente os parametros de geracao da OpenAI (max_tokens, temperature, modelo) sao fixos via env vars no backend. O usuario so controla o `level` (funcional/tecnico/executivo). Este plano torna esses parametros configuraveis por request, diretamente nos formularios do frontend, com defaults pre-selecionados e descricoes explicativas para o usuario.

### Parametros a expor

| Parametro | Tipo UI | Range | Default | Descricao |
|---|---|---|---|---|
| Temperature | slider | 0.0 - 1.0 (step 0.1) | 0.3 | Criatividade da resposta |
| Max Tokens | grupo de botoes | 256, 512, 1024, 2048, 4096 | 1024 | Tamanho maximo da resposta |
| Modelo | cards selecionaveis | Auto, GPT-4o Mini, GPT-4o | Auto | Modelo de IA utilizado |

---

## Etapa 1: Backend - Struct de overrides no DTO

**Arquivo:** `backend/internal/handler/dto/request.go`

- [ ] Criar struct `GenerationOverrides` com campos ponteiro (`*int`, `*float64`, `*string`)
  - `MaxTokens *int` (json: `max_tokens,omitempty`)
  - `Temperature *float64` (json: `temperature,omitempty`)
  - `Model *string` (json: `model,omitempty`) — valores: `"auto"`, `"gpt-4o-mini"`, `"gpt-4o"`
- [ ] Criar mapa `validModels` com os modelos permitidos
- [ ] Criar metodo `Validate()` na struct com regras:
  - max_tokens: entre 64 e 4096
  - temperature: entre 0.0 e 2.0
  - model: deve estar no mapa `validModels`
  - nil e aceito (significa "usar default do servidor")
- [ ] Adicionar campo `Overrides *GenerationOverrides` (json: `overrides,omitempty`) em:
  - `AnalyzeCommitRequest`
  - `AnalyzeRangeRequest`
  - `AnalyzePRRequest`
- [ ] Chamar `r.Overrides.Validate()` dentro do `Validate()` de cada DTO (quando Overrides != nil)

---

## Etapa 2: Backend - Campos de override no GenerationInput

**Arquivo:** `backend/internal/openai/generator.go`

- [x] Adicionar campos opcionais ao `GenerationInput`:
  - `MaxTokensOverride *int`
  - `TemperatureOverride *float64`
  - `ModelOverride *string`
- [x] Criar metodo helper `HasOverrides() bool` no `GenerationInput` que retorna true se qualquer override esta ativo (diferente de nil, e Model diferente de "auto")
- [x] Criar metodo `resolveModel(input GenerationInput, tokenCount int) string` no generator:
  - Se `ModelOverride != nil && != "auto"` → retorna o modelo diretamente
  - Senao → delega para `SelectModel()` (comportamento atual)
- [x] Modificar metodo `Generate()`:
  - Resolver `maxTokens`: usar `input.MaxTokensOverride` se presente, senao `g.config.MaxTokens`
  - Resolver `temperature`: usar `float32(*input.TemperatureOverride)` se presente, senao `g.config.Temperature`
  - Resolver `model`: usar `resolveModel()` ao inves de `SelectModel()` direto
  - Cache: quando `input.HasOverrides()` == true, pular leitura e escrita de cache

---

## Etapa 3: Backend - Threading dos overrides no Service

**Arquivo:** `backend/internal/service/analysis_service.go`

- [x] Criar funcao helper `applyOverrides(input *openai.GenerationInput, overrides *dto.GenerationOverrides)`:
  - Copia `overrides.MaxTokens` → `input.MaxTokensOverride`
  - Copia `overrides.Temperature` → `input.TemperatureOverride`
  - Copia `overrides.Model` → `input.ModelOverride`
  - Se overrides == nil, nao faz nada
- [x] Criar funcao helper `hasOverrides(overrides *dto.GenerationOverrides) bool`:
  - Retorna true se algum campo e non-nil (e Model != "auto")
- [x] Em `AnalyzeCommit()`:
  - Envolver cache do service (`s.cache.Get` na linha 79) e DB lookup (`s.repository.GetByDiffHash` na linha 98) em `if !hasOverrides(req.Overrides)`
  - Chamar `applyOverrides(&genInput, req.Overrides)` antes de `s.generator.Generate()`
- [x] Em `AnalyzeRange()`:
  - Mesmo tratamento: guardar cache (linha 171) e DB (linha 190) com `!hasOverrides()`
  - Chamar `applyOverrides` antes de `Generate()`
- [x] Em `AnalyzePR()`:
  - Mesmo tratamento: guardar cache (linha 262) e DB (linha 281) com `!hasOverrides()`
  - Chamar `applyOverrides` antes de `Generate()`

---

## Etapa 4: Frontend - Tipo compartilhado

**Arquivo:** `frontend/src/lib/api/types.ts`

- [x] Criar interface `GenerationOverrides`:
  - `max_tokens?: number`
  - `temperature?: number`
  - `model?: string`
- [x] Adicionar campo `overrides?: GenerationOverrides` em:
  - `AnalyzeCommitRequest`
  - `AnalyzeRangeRequest`
  - `AnalyzePRRequest`

---

## Etapa 5: Frontend - Componente AdvancedSettings

**Novo arquivo:** `frontend/src/features/shared/AdvancedSettings.tsx`

- [x] Criar componente controlado com props `{ value: GenerationOverrides, onChange, disabled? }`
- [x] Toggle colapsavel (comeca fechado):
  - Icone `Settings2` + texto "Configuracoes avancadas" + chevron de direcao
  - Ao clicar, expande/recolhe o painel de configs
- [x] Secao **Temperature**:
  - Slider HTML range (0.0 a 1.0, step 0.1)
  - Label mostrando valor atual (ex: "Temperature: 0.3")
  - Descricao: "Controla a criatividade. Valores baixos = mais preciso, altos = mais criativo."
  - Escala visual: "Preciso (0.0)" ←→ "Criativo (1.0)"
- [x] Secao **Max Tokens**:
  - Grupo de 5 botoes: 256, 512, 1024, 2048, 4096
  - Default 1024 marcado com estilo violet
  - Descricao: "Limite de tokens na resposta gerada. Mais tokens = descricao mais longa."
- [x] Secao **Modelo**:
  - 3 cards estilo `LevelSelector` (grid 3 colunas):
    - Auto: "Selecao automatica baseada no tamanho e tipo"
    - GPT-4o Mini: "Rapido e eficiente para diffs simples"
    - GPT-4o: "Mais capaz, ideal para diffs complexos"
  - Default "Auto" marcado
  - Descricao: "Escolha o modelo de IA. 'Auto' seleciona automaticamente com base no tipo de analise."
- [x] Estilo consistente com `LevelSelector.tsx`:
  - Violet para selecionado, stone para inativo
  - Suporte dark mode com prefixos `dark:`
  - `rounded-xl`, borders, transitions

---

## Etapa 6: Frontend - Integracao nos formularios

### CommitForm (`frontend/src/features/commit/CommitForm.tsx`)

- [x] Importar `AdvancedSettings` e `GenerationOverrides`
- [x] Adicionar estado: `const [overrides, setOverrides] = useState<GenerationOverrides>({})`
- [x] Colocar `<AdvancedSettings value={overrides} onChange={setOverrides} disabled={isPending} />` entre `<LevelSelector>` e `<Button>`
- [x] No `handleSubmit`: montar `effectiveOverrides` somente com valores que diferem dos defaults (temperature != 0.3, max_tokens != 1024, model != "auto"); incluir `overrides` no payload apenas se houver algum override ativo

### PrForm (`frontend/src/features/pull-request/PrForm.tsx`)

- [x] Mesma integracao do CommitForm
- [x] Importar, estado, JSX, handleSubmit

### RangeForm (`frontend/src/features/range/RangeForm.tsx`)

- [x] Mesma integracao do CommitForm
- [x] Importar, estado, JSX, handleSubmit

---

## Etapa 7: Testes

### Backend

- [x] **DTO**: testar `GenerationOverrides.Validate()` — limites invalidos, modelo invalido, nil (passa), valores validos
- [x] **Generator**: testar que overrides sao usados no `ChatCompletionRequest` enviado ao OpenAI
- [x] **Generator**: testar que cache e bypassado quando overrides estao presentes
- [x] **Generator**: testar que `model = "auto"` ainda delega para `SelectModel()`
- [x] **Service**: testar que overrides fluem do DTO ate o `GenerationInput` passado ao generator

### Frontend

- [x] **AdvancedSettings.test.tsx**: toggle abre/fecha, slider atualiza callback, botoes de max_tokens funcionam, cards de modelo funcionam, disabled propaga
- [x] **CommitForm.test.tsx**: overrides nao incluidos com defaults; overrides incluidos quando alterados
- [x] **PrForm.test.tsx**: mesmo que CommitForm
- [x] **RangeForm.test.tsx**: mesmo que CommitForm

---

## Verificacao final

- [x] `cd backend && go test ./...` — todos os testes passam
- [x] `cd frontend && npm test` — todos os testes passam
- [ ] Teste manual: abrir form, expandir configs avancadas, alterar parametros, submeter e verificar `model_used` na resposta
- [ ] Teste de retrocompatibilidade: submeter sem abrir configs → payload sem `overrides`, comportamento identico ao atual
