Quero que você atue como um arquiteto de software sênior e gere a base de uma plataforma fullstack com frontend em React e backend em Go.

## 🎯 Objetivo da Plataforma

Criar uma plataforma web para **gerar descrições automáticas de commits ou ranges de commits**, com foco em:

* Gerar descrição clara de um único commit
* Gerar descrição consolidada de um range de commits (ex: hash X até hash Y)
* Analisar um Pull Request e gerar uma descrição resumida do que foi implementado
* Gerar uma descrição **pouco técnica**, focada em QA/PO
* Permitir que o usuário peça para **adaptar/refinar uma descrição já gerada**

O objetivo principal é acelerar o preenchimento de cards no JIRA e melhorar a comunicação entre Dev, QA e PO.

---

# 🏗 Arquitetura Obrigatória

## Frontend

* React
* Vite
* TailwindCSS
* Lucide Icons
* Sem autenticação
* Interface simples e moderna
* Estrutura escalável (feature-based folder structure)

## Backend

* Golang
* Arquitetura limpa (handler → service → repository)
* PostgreSQL
* Separação clara de camadas
* Injeção de dependência simples
* API REST

## Banco de Dados

PostgreSQL

Sugira modelagem para:

* Histórico de análises
* Tipo de análise (commit único, range, PR)
* Descrição gerada
* Descrição adaptada
* Data de criação

---

# 🚀 Funcionalidades

## 1️⃣ Análise de Commit

Input:

* Hash do commit
* Ou diff colado manualmente

Output:

* Descrição clara e resumida
* Linguagem pouco técnica
* Foco em impacto funcional

---

## 2️⃣ Análise de Range de Commits

Input:

* Hash inicial
* Hash final

Output:

* Resumo consolidado das mudanças
* Agrupamento por tipo (ex: melhorias, correções, refatorações)

---

## 3️⃣ Análise de Pull Request

Input:

* Texto do PR
* Ou diff
* Ou título + descrição

Output:

* Descrição final pronta para card do JIRA
* Linguagem adequada para QA/PO

---

## 4️⃣ Refinamento de Descrição

Input:

* Descrição já gerada
* Instrução do usuário (ex: "deixe mais simples", "mais técnico", "mais resumido")

Output:

* Nova versão adaptada

---

# 📐 Requisitos Técnicos

* Backend deve ser preparado para integração futura com API de LLM
* Criar interface clara para serviço de geração de descrição (ex: DescriptionGenerator interface)
* Separar regras de negócio da camada HTTP
* Criar DTOs bem definidos
* Criar migrations SQL
* Criar README explicando como rodar frontend e backend
* Usar variáveis de ambiente
* Estruturar projeto para produção

---

# 🎨 Frontend – UX Esperada

Criar interface com:

* Tabs:

  * Commit
  * Range
  * PR
  * Refinar descrição
* Campo de input grande para colar diff
* Campo opcional para hash
* Botão "Gerar descrição"
* Área de resultado com:

  * Botão copiar
  * Botão "Refinar"

Layout moderno com Tailwind.
Ícones com Lucide.
Sem autenticação.
Responsivo.

---

# 📦 Entregáveis Esperados

1. Estrutura completa de pastas frontend
2. Estrutura completa de pastas backend
3. Modelagem do banco
4. Exemplos de endpoints
5. Exemplos de request/response
6. Código inicial funcional
7. Sugestão de melhorias futuras
8. Estratégia para integração com LLM

---

# 🧠 Extra

Inclua sugestões como:

* Possibilidade futura de integração com GitHub API
* Webhook para PRs
* Geração automática via CI
* Exportação para Markdown
* Geração em múltiplos níveis: técnico / funcional / executivo
