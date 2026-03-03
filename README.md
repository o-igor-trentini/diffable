# Diffable

Plataforma fullstack (React + Go) para gerar descrições automáticas de commits e PRs do Bitbucket Cloud, usando OpenAI para geração de texto focado em QA/PO e JIRA.

## Quick Start

```bash
# Subir todos os serviços
docker compose up --build

# Backend: http://localhost:8080
# Frontend: http://localhost:3000
```

## Desenvolvimento

```bash
# Modo dev com hot reload
make dev

# Rodar todos os testes
make test-all
```

## Variáveis de Ambiente

Consulte `backend/.env.example` e `frontend/.env.example`.
