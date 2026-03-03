# Diffable

Plataforma fullstack (React + Go) para gerar descrições automáticas de commits e PRs do Bitbucket Cloud, usando OpenAI para geração de texto focado em QA/PO e JIRA.

## Features

- Análise de commits individuais, ranges e Pull Requests
- Geração de descrições em linguagem não-técnica via OpenAI
- Refinamento iterativo com instruções customizadas
- Histórico de análises com filtros e paginação
- Cache inteligente (in-memory + deduplicação por hash)
- Dark mode com persistência de preferência
- Atalhos de teclado (Ctrl+Enter, Ctrl+Shift+C)
- Rate limiting por IP (configurável)
- Logs estruturados JSON com correlação por Request ID

## Arquitetura

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   Frontend   │────▶│   Backend    │────▶│  PostgreSQL  │
│  React/Vite  │     │   Go/Chi     │     │   16-alpine  │
│  :3000       │     │   :8080      │     │   :5432      │
└─────────────┘     └──────┬───────┘     └──────────────┘
                           │
                    ┌──────┴───────┐
                    │              │
              ┌─────▼─────┐ ┌─────▼─────┐
              │ Bitbucket  │ │  OpenAI   │
              │ Cloud API  │ │  API      │
              └───────────┘ └───────────┘
```

## Quick Start

```bash
# 1. Configure as variáveis de ambiente
cp backend/.env.example backend/.env
# Edite backend/.env com suas credenciais

# 2. Suba todos os serviços
docker compose up --build

# 3. Acesse
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080
```

## Desenvolvimento

```bash
# Modo dev com hot reload
make dev

# Rodar todos os testes
make test-all

# Ver logs
make logs

# Parar serviços
make down
```

### Sem Docker

**Backend:**
```bash
cd backend
cp .env.example .env
# Edite .env com DATABASE_URL e credenciais
go run cmd/server/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Configuração

### Backend

| Variável | Default | Descrição |
|----------|---------|-----------|
| `PORT` | 8080 | Porta do servidor HTTP |
| `DATABASE_URL` | — | Connection string PostgreSQL |
| `FRONTEND_URL` | http://localhost:3000 | URL do frontend (CORS) |
| `LOG_LEVEL` | info | Nível de log (debug, info, warn, error) |
| `SHUTDOWN_TIMEOUT` | 30s | Timeout para graceful shutdown |
| `BITBUCKET_BASE_URL` | https://api.bitbucket.org/2.0 | URL base da API Bitbucket |
| `BITBUCKET_EMAIL` | — | Email para autenticação Bitbucket |
| `BITBUCKET_API_TOKEN` | — | App Password do Bitbucket |
| `BITBUCKET_TIMEOUT` | 30s | Timeout das requisições ao Bitbucket |
| `OPENAI_API_KEY` | — | Chave da API OpenAI |
| `OPENAI_DEFAULT_MODEL` | gpt-4o-mini | Modelo para diffs simples |
| `OPENAI_COMPLEX_MODEL` | gpt-4o | Modelo para diffs complexos |
| `OPENAI_MAX_TOKENS` | 1024 | Máximo de tokens na resposta |
| `OPENAI_TEMPERATURE` | 0.3 | Temperatura de geração |
| `OPENAI_TOKEN_THRESHOLD` | 4000 | Limite para usar modelo complexo |
| `CACHE_TTL` | 24h | TTL do cache in-memory |
| `RATE_LIMIT_RPM` | 60 | Requisições por minuto por IP |
| `DB_MAX_CONNS` | 25 | Máximo de conexões no pool |
| `DB_MIN_CONNS` | 5 | Mínimo de conexões no pool |

### Frontend

| Variável | Default | Descrição |
|----------|---------|-----------|
| `VITE_API_URL` | /api/v1 | URL base da API |

## Testes

```bash
# Todos os testes
make test-all

# Apenas backend
cd backend && go test ./...

# Apenas frontend
cd frontend && npm test -- --run
```

## Estrutura do Projeto

```
diffable/
├── backend/
│   ├── cmd/server/main.go          # Entry point
│   ├── internal/
│   │   ├── config/                 # Configuração via env vars
│   │   ├── server/                 # Router e middleware setup
│   │   ├── handler/                # HTTP handlers
│   │   │   └── dto/                # Request/Response DTOs
│   │   ├── service/                # Lógica de negócio
│   │   ├── repository/            # Camada de persistência
│   │   ├── domain/                 # Modelos e erros
│   │   ├── bitbucket/              # Cliente Bitbucket Cloud
│   │   ├── openai/                 # Gerador de descrições
│   │   ├── cache/                  # Cache in-memory
│   │   └── middleware/             # CORS, logging, rate limit, request ID
│   └── migrations/                 # SQL migrations
├── frontend/
│   └── src/
│       ├── features/               # Componentes por funcionalidade
│       │   ├── commit/             # Análise de commit
│       │   ├── range/              # Análise de range
│       │   ├── pull-request/       # Análise de PR
│       │   ├── refine/             # Refinamento
│       │   ├── history/            # Histórico
│       │   └── shared/             # Componentes compartilhados
│       └── lib/
│           ├── api/                # Cliente HTTP e tipos
│           └── hooks/              # React hooks customizados
├── docs/                           # Documentação
├── docker-compose.yml              # Produção
└── docker-compose.dev.yml          # Desenvolvimento
```

## API

Documentação completa em [docs/api.md](docs/api.md).

| Método | Path | Descrição |
|--------|------|-----------|
| GET | `/healthz` | Liveness probe |
| GET | `/readyz` | Readiness probe |
| POST | `/api/v1/analyses/commit` | Análise de commit único |
| POST | `/api/v1/analyses/range` | Análise de range de commits |
| POST | `/api/v1/analyses/pr` | Análise de PR |
| GET | `/api/v1/analyses` | Listar histórico (paginado) |
| GET | `/api/v1/analyses/{id}` | Buscar análise por ID |
| POST | `/api/v1/analyses/{id}/refine` | Refinar descrição |
| GET | `/api/v1/analyses/{id}/refinements` | Listar refinamentos |

## Decisões Arquiteturais

- [Arquitetura](docs/architecture.md)
- [Integração Bitbucket](docs/bitbucket-integration.md)
- [Integração OpenAI](docs/openai-integration.md)
- [Fases de desenvolvimento](docs/phases.md)
- [Referência da API](docs/api.md)
