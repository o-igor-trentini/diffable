# Nivel QA Detalhado + Contexto Adicional

## Contexto

As descricoes geradas pela plataforma sao insuficientes para uso em cards Jira para validacao de QA. O prompt atual do nivel "funcional" limita a saida a ~200 palavras com formato generico (Resumo + Mudancas + Impacto). O QA precisa de descricoes completas que cubram: mudancas de banco de dados, fluxos de negocio (com cenarios alternativos e de borda), regras de validacao, mudancas de API/integracoes, e cenarios de teste sugeridos.

Alem disso, o preprocessamento do diff remove TODAS as linhas de contexto, dificultando a compreensao do codigo ao redor das mudancas. Nao existe campo para o usuario fornecer contexto adicional sobre a tarefa.

### Resultado esperado

Ao selecionar o nivel "QA Detalhado", a descricao gerada deve conter secoes como:
- **Contexto**: O que a mudanca faz e por que foi necessaria
- **Mudancas no Banco de Dados**: Tabelas, colunas, migracoes
- **Mudancas de API/Integracoes**: Endpoints, servicos externos, novos campos
- **Regras de Negocio**: Validacoes, logica condicional, restricoes
- **Fluxos Afetados**: Caminho feliz, alternativos, casos de borda
- **Cenarios de Teste Sugeridos**: Cenarios com dados de entrada e resultado esperado
- **Observacoes**: Retrocompatibilidade, dependencias, riscos

---

## Etapa 1: Backend - Novo nivel `qa_detailed` no DTO

**Arquivo:** `backend/internal/handler/dto/request.go`

- [x] Adicionar `"qa_detailed": true` no mapa `validLevels` (linha 9)
- [x] Atualizar mensagem de erro em `validateLevel` (linha 49) para incluir `qa_detailed`
- [x] Adicionar campo `UserContext string` (json: `user_context,omitempty`) em:
  - `AnalyzePRRequest`
  - `AnalyzeCommitRequest`
  - `AnalyzeRangeRequest`
- [x] Aumentar limite superior de `max_tokens` de 4096 para 8192 em `GenerationOverrides.Validate()` (linha 32)

---

## Etapa 2: Backend - Campo UserContext no GenerationInput

**Arquivo:** `backend/internal/openai/generator.go`

- [x] Adicionar campo `UserContext string` ao `GenerationInput` (apos `Level`, linha 23)
- [x] No metodo `Generate()` (linha 110): trocar `PreprocessDiff(input.Diff)` por `PreprocessDiffForLevel(input.Diff, level)`
- [x] No metodo `Generate()` (linhas 114-117): adicionar default de `maxTokens = 4096` quando `level == "qa_detailed"` e nao houver override:
  ```go
  maxTokens := g.config.MaxTokens
  if input.MaxTokensOverride != nil {
      maxTokens = *input.MaxTokensOverride
  } else if level == "qa_detailed" {
      maxTokens = 4096
  }
  ```
- [x] No metodo `buildMessages()` (linha 201): trocar `buildFewShotExamples()` por `buildFewShotExamplesForLevel(level)`
- [x] No metodo `buildMessages()` (linha 208): passar `input.UserContext` para `buildUserPrompt()`

---

## Etapa 3: Backend - Preprocessamento com contexto para QA

**Arquivo:** `backend/internal/openai/tokenizer.go`

- [x] Criar funcao `PreprocessDiffForLevel(rawDiff, level string) string`:
  - Para `qa_detailed`: chama `reduceContextWithLines(section, 5)` ao inves de `reduceContext(section)`
  - Para demais niveis: chama `reduceContext(section)` (comportamento atual inalterado)
- [x] Criar funcao `reduceContextWithLines(section string, contextLines int) string`:
  - Identifica indices das linhas com mudancas (+/-)
  - Monta um set de indices a manter (N linhas antes e depois de cada mudanca)
  - Sempre mantem headers (diff --git, ---, +++, @@)
  - Mantem linhas de contexto (sem prefixo + ou -) que estejam dentro do range

**Racional:** 5 linhas de contexto permitem que a IA veja assinaturas de funcao, schemas de tabela, nomes de variaveis ao redor das mudancas. Aumenta tokens em ~30-50% mas e essencial para a qualidade do nivel QA.

---

## Etapa 4: Backend - System prompt QA e few-shot examples

**Arquivo:** `backend/internal/openai/prompts.go`

- [x] Adicionar case `"qa_detailed"` em `buildSystemPromptForLevel()` com prompt detalhado:
  - Persona: analista de qualidade senior
  - Instrucoes para analisar cada aspecto (banco, API, regras de negocio, fluxos, cenarios de teste)
  - Sem limite rigido de palavras: priorizar completude sobre brevidade
  - Formato obrigatorio com 7 secoes (omitir secoes que nao se aplicam)
- [x] Refatorar `buildFewShotExamples()` em `buildFewShotExamplesForLevel(level string)`:
  - Manter `buildFewShotExamples()` delegando para `buildFewShotExamplesForLevel("functional")` (retrocompatibilidade)
  - `"qa_detailed"` retorna exemplo dedicado com todas as secoes preenchidas (cenario de integracao CNH/Nexus como referencia)
  - Demais niveis retornam os exemplos atuais (login lockout + PIX)
- [x] Modificar assinatura de `buildUserPrompt()` para aceitar parametro `userContext string`:
  - Quando nao vazio, insere antes do titulo do PR: `\n\nContexto adicional fornecido pelo desenvolvedor:\n{userContext}`
  - Quando vazio, comportamento identico ao atual

---

## Etapa 5: Backend - Threading do UserContext no Service

**Arquivo:** `backend/internal/service/analysis_service.go`

- [x] Em `AnalyzePR()` (linha 323): adicionar `UserContext: req.UserContext` no `GenerationInput`
- [x] Em `AnalyzeCommit()`: mesmo tratamento — adicionar `UserContext: req.UserContext`
- [x] Em `AnalyzeRange()`: mesmo tratamento — adicionar `UserContext: req.UserContext`

---

## Etapa 6: Frontend - Tipos da API

**Arquivo:** `frontend/src/lib/api/types.ts`

- [x] Adicionar `user_context?: string` em:
  - `AnalyzePRRequest`
  - `AnalyzeCommitRequest`
  - `AnalyzeRangeRequest`

---

## Etapa 7: Frontend - Level Selector

**Arquivo:** `frontend/src/features/shared/LevelSelector.tsx`

- [ ] Importar icone `TestTube2` de `lucide-react`
- [ ] Adicionar 4a opcao no array `levels` (entre Funcional e Tecnico):
  - value: `qa_detailed`
  - label: `QA Detalhado`
  - description: `Para QA validar cards Jira. Descricao completa com fluxos, regras de negocio e cenarios de teste.`
  - icon: `<TestTube2 size={18} />`
- [ ] Ajustar grid de `sm:grid-cols-3` para `sm:grid-cols-2 lg:grid-cols-4`

---

## Etapa 8: Frontend - Campo Contexto Adicional nos formularios

### PrForm (`frontend/src/features/pull-request/PrForm.tsx`)

- [ ] Adicionar estado `const [userContext, setUserContext] = useState('')`
- [ ] Adicionar `<TextArea>` para "Contexto Adicional" entre `LevelSelector` e `AdvancedSettings`:
  - label: `Contexto Adicional`
  - placeholder: `Ex: Este PR integra a data de emissao da CNH retornada pela API Nexus. O campo e opcional e pode vir vazio...`
  - hint: `Opcional. Forneca contexto sobre a tarefa para melhorar a descricao gerada. Especialmente util no nivel QA Detalhado.`
  - rows: 3
- [ ] No `handleSubmit`: incluir `user_context: userContext.trim()` no payload quando nao vazio

### CommitForm (`frontend/src/features/commit/CommitForm.tsx`)

- [ ] Mesma integracao do PrForm (estado, TextArea, handleSubmit)

### RangeForm (`frontend/src/features/range/RangeForm.tsx`)

- [ ] Mesma integracao do PrForm (estado, TextArea, handleSubmit)

---

## Etapa 9: Frontend - Opcao de 8192 tokens

**Arquivo:** `frontend/src/features/shared/AdvancedSettings.tsx`

- [ ] Adicionar `8192` ao array `TOKEN_OPTIONS`

---

## Etapa 10: Testes

### Backend

- [ ] **DTO** (`request_test.go`): testar que `qa_detailed` e aceito como nivel valido
- [ ] **DTO** (`request_test.go`): testar que `max_tokens = 8192` passa validacao e `8193` falha
- [ ] **Tokenizer** (`tokenizer_test.go`): testar `PreprocessDiffForLevel` com `qa_detailed` mantem linhas de contexto
- [ ] **Tokenizer** (`tokenizer_test.go`): testar `PreprocessDiffForLevel` com `functional` remove contexto (regressao)
- [ ] **Prompts** (`prompts_test.go`): testar `buildSystemPromptForLevel("qa_detailed")` contem secoes esperadas
- [ ] **Prompts** (`prompts_test.go`): testar `buildFewShotExamplesForLevel("qa_detailed")` retorna exemplos diferentes
- [ ] **Prompts** (`prompts_test.go`): testar `buildUserPrompt` com `userContext` nao vazio inclui contexto
- [ ] **Generator** (`generator_test.go`): testar que `qa_detailed` usa 4096 max_tokens por default
- [ ] **Service** (`analysis_service_test.go`): testar que `UserContext` flui ate o `GenerationInput`

### Frontend

- [ ] **LevelSelector**: testar que 4 opcoes sao renderizadas, incluindo "QA Detalhado"
- [ ] **PrForm** (`PrForm.test.tsx`): testar que `user_context` e incluido no payload quando preenchido
- [ ] **PrForm** (`PrForm.test.tsx`): testar que `user_context` e omitido quando vazio
- [ ] **CommitForm** (`CommitForm.test.tsx`): mesmo que PrForm
- [ ] **RangeForm** (`RangeForm.test.tsx`): mesmo que PrForm

---

## Verificacao final

- [ ] `cd backend && go test ./...` — todos os testes passam
- [ ] `cd frontend && npm test` — todos os testes passam
- [ ] Teste manual: selecionar nivel "QA Detalhado", preencher contexto adicional, submeter PR e verificar que a descricao contem todas as secoes
- [ ] Teste de retrocompatibilidade: niveis existentes (funcional, tecnico, executivo) continuam funcionando identicamente
- [ ] Teste sem contexto adicional: submeter sem preencher contexto — comportamento deve ser identico ao anterior
