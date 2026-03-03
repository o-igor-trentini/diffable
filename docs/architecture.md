# Arquitetura

## Decisões Técnicas

- **HTTP Router:** go-chi/chi — leve, idiomático, compatível com net/http
- **PostgreSQL Driver:** jackc/pgx — driver nativo, melhor performance que database/sql
- **Frontend State:** @tanstack/react-query — cache automático, refetch, mutations
- **Arquitetura:** Clean Architecture — handler → service → repository/client
- **CSS:** Tailwind CSS — utility-first, produtivo para UI rápida

## Estrutura

```
backend/
  cmd/server/       → entry point
  internal/
    config/         → variáveis de ambiente
    server/         → chi router + middlewares
    handler/        → HTTP handlers
    service/        → lógica de negócio
    repository/     → acesso a dados
    domain/         → modelos e erros
    bitbucket/      → cliente API Bitbucket
    openai/         → cliente OpenAI
    middleware/      → middlewares customizados

frontend/src/
  features/         → componentes por funcionalidade
  lib/              → API client, hooks
```
