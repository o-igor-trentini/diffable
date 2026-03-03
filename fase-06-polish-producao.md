# Fase 6: Polish e Produção

**Objetivo:** Preparar a aplicação para uso em produção com logging estruturado, rate limiting, tratamento de erros robusto, UX polida (responsivo, dark mode), documentação completa e otimizações de performance.

**Pré-requisitos:** Fase 5 concluída

---

## Checklist Backend

### Middleware — Rate Limiting

- [ ] Criar `backend/internal/middleware/ratelimit.go`:
  - Token bucket rate limiter por IP
  - Configurável: `RATE_LIMIT_RPM` (default: 60 req/min por IP)
  - Retorna 429 com header `Retry-After` quando limite excedido
  - Usa `sync.Map` para tracking por IP
  - Limpeza periódica de IPs antigos (goroutine com ticker)
- [ ] Testes: permite sob limite, bloqueia acima, retorna 429 com Retry-After

### Middleware — Request ID

- [ ] Criar `backend/internal/middleware/requestid.go`:
  - Garante que request ID propaga via `context.Context`
  - Todos logs de service/repository incluem request ID
  - Header `X-Request-ID` na resposta

### Logging Aprimorado

- [ ] Atualizar `backend/internal/middleware/logging.go`:
  - Log estruturado JSON com `slog`:
    - Request: method, path, user-agent, content-length
    - Response: status code, duration, bytes written
    - Request ID em cada entrada
  - Log levels configuráveis via `LOG_LEVEL`

### Error Handling Robusto

- [ ] Atualizar `backend/internal/handler/analysis_handler.go`:
  - Mapeamento centralizado de erros domain → HTTP:
    - `ErrNotFound` → 404
    - `ErrValidation` → 400
    - `ErrRateLimited` → 429
    - `ErrExternalService` → 502
    - `ErrTimeout` → 504
    - `ErrTokenLimitExceeded` → 422
    - Qualquer outro → 500
  - Sempre retorna `ErrorResponse` JSON com mensagem amigável
  - Log do erro completo internamente, retorna mensagem sanitizada ao client
- [ ] Atualizar `backend/internal/domain/errors.go`:
  - Adicionar `Is()` e `As()` para unwrapping
  - Adicionar `ErrTimeout`, `ErrTokenLimitExceeded`
  - Error wrapping com contexto em cada camada

### Graceful Shutdown

- [ ] Atualizar `backend/internal/server/server.go`:
  - Listen para `SIGTERM` e `SIGINT`
  - `server.Shutdown(ctx)` com timeout configurável (`SHUTDOWN_TIMEOUT`, default 30s)
  - Drena conexões em andamento
  - Fecha pool de conexões DB
  - Log de início e fim do shutdown

### Connection Pool Tuning

- [ ] Atualizar `backend/cmd/server/main.go`:
  - pgx pool config:
    - `MaxConns` = `DB_MAX_CONNS` (default 25)
    - `MinConns` = `DB_MIN_CONNS` (default 5)
    - `MaxConnLifetime` = 30 minutos
    - `MaxConnIdleTime` = 5 minutos
    - Health check period = 1 minuto

### Configuração

- [ ] Atualizar `backend/internal/config/config.go`:
  - `RateLimitRPM` (default: 60)
  - `DBMaxConns` (default: 25)
  - `DBMinConns` (default: 5)
  - `ShutdownTimeout` (default: 30s)

---

## Checklist Frontend

### Error Boundary

- [ ] Criar `frontend/src/features/shared/ErrorBoundary.tsx`:
  - React Error Boundary wrapping toda app
  - Mostra página amigável com:
    - Mensagem: "Algo deu errado"
    - Botão "Tentar novamente" (recarrega componente)
    - Detalhes do erro em modo dev
  - Log do erro no console

### Retry de Mutations

- [ ] Criar `frontend/src/features/shared/RetryButton.tsx`:
  - Botão que executa retry de mutations falhadas
  - Usa `reset()` do TanStack Query
  - Texto: "Tentar novamente"

### API Client Aprimorado

- [ ] Atualizar `frontend/src/lib/api/client.ts`:
  - Retry interceptor: 1 retry automático em 5xx
  - Request timeout: 30s (abort controller)
  - Extração melhorada de mensagem de erro do response
  - Header `X-Request-ID` para correlação

### Dark Mode

- [ ] Atualizar `frontend/tailwind.config.ts`:
  - `darkMode: 'class'`
  - Cores dark customizadas para todos tokens
- [ ] Criar toggle dark mode no `frontend/src/App.tsx`:
  - Ícone Moon/Sun (Lucide)
  - Persiste preferência no `localStorage`
  - Aplica classe `dark` no elemento root
  - Respeita `prefers-color-scheme` como default inicial
- [ ] Adicionar variantes `dark:` em todos componentes:
  - Backgrounds, textos, bordas, cards, inputs, buttons
  - ResultDisplay, TabNavigation, formulários, alertas

### Responsividade

- [ ] Verificar e ajustar todos componentes para:
  - Mobile (375px): formulários empilhados verticalmente, botões full-width
  - Tablet (768px): layout intermediário
  - Desktop (1280px): layout lado a lado quando aplicável
- [ ] Tab Navigation: horizontal em desktop, scrollável em mobile
- [ ] Formulários: campos empilham em mobile
- [ ] ResultDisplay: full-width em todos breakpoints
- [ ] HistoryPanel: overlay/drawer em mobile, sidebar em desktop

### Keyboard Shortcuts

- [ ] Criar `frontend/src/features/shared/KeyboardShortcuts.tsx`:
  - `Ctrl+Enter`: submete formulário ativo
  - `Ctrl+Shift+C`: copia resultado
  - Tooltip com dicas de atalho nos botões relevantes

### Testes Frontend

- [ ] Teste: ErrorBoundary captura erro de render, mostra botão retry
- [ ] Teste: Dark mode toggle adiciona/remove classe `dark`, persiste no localStorage
- [ ] Teste: Responsivo — layout empilha em mobile viewport

---

## Checklist Raiz

### Docker Compose Produção

- [ ] Atualizar `docker-compose.yml`:
  - Resource limits para cada serviço:
    ```yaml
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          memory: 256M
    ```
  - Health checks para backend: `curl -f http://localhost:8080/healthz`
  - Restart policies: `restart: unless-stopped`
  - Logging config: json-file com max-size

### Documentação Completa

- [ ] Criar/atualizar `docs/api.md`:
  - Todos endpoints documentados com:
    - URL, método, descrição
    - Request body (JSON schema)
    - Response body (JSON schema)
    - Códigos de erro possíveis
    - Exemplos curl
  - Rate limiting: limites e headers
  - Autenticação Bitbucket
- [ ] Atualizar `README.md` — versão final:
  - Overview do projeto e features
  - Diagrama de arquitetura (ASCII)
  - Quick start com Docker Compose (3 comandos)
  - Configuração completa (tabela de env vars)
  - Desenvolvimento local (com e sem Docker)
  - Rodando testes
  - Estrutura do projeto
  - Decisões arquiteturais (link para docs/)
  - Contributing guide básico

---

## Testes desta Fase

| Teste | Tipo | Validação |
|-------|------|-----------|
| Rate limiter middleware | Unit | Permite sob limite, bloqueia acima, 429 + Retry-After |
| Error mapping | Unit | Cada domain error mapeia para HTTP status correto |
| Graceful shutdown | Integration | Server drena requests antes de parar |
| Request ID propagation | Unit | ID presente em logs e response header |
| ErrorBoundary.test.tsx | Unit (Vitest) | Captura erro, mostra retry |
| Dark mode toggle | Unit (Vitest) | Toggle altera classe, persiste localStorage |
| Responsivo | E2E (Playwright) | Mobile: forms empilham. Desktop: side-by-side |
| Full E2E happy path | E2E (Playwright) | Gerar commit → copy → refinar → histórico → dark mode |
| Keyboard shortcuts | Unit (Vitest) | Ctrl+Enter submete form |

---

## Critérios de Aceite

- [ ] Todos erros da API retornam `ErrorResponse` JSON consistente
- [ ] Rate limiting: 60 req/min por IP (configurável), retorna 429 com `Retry-After`
- [ ] Logs JSON estruturados com request ID correlacionando handler → service → repository
- [ ] Graceful shutdown: requests em andamento completam, novos requests rejeitados
- [ ] Frontend: Error boundary captura crashes, mostra botão retry
- [ ] Frontend: Dark mode toggle funciona, preferência persiste entre sessões
- [ ] Frontend: Responsivo em mobile (375px), tablet (768px), desktop (1280px)
- [ ] `Ctrl+Enter` submete formulário ativo
- [ ] `docs/api.md` documenta todos endpoints com exemplos curl
- [ ] `README.md` tem diagrama de arquitetura, quick start e referência completa
- [ ] Docker Compose produção com health checks e resource limits
- [ ] Todos testes passam (unit, integration, E2E)
- [ ] `make test-all` executa suite completa sem falhas
