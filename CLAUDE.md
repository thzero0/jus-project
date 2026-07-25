# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Contexto do projeto

Este repositório implementa o desafio técnico "Justarter": uma página única de busca com autocompletar. Domínio escolhido: top 10k jogos mais populares de 2025 (dataset do Kaggle, ver `projeto.md`).

Arquitetura escolhida (`docs/requirements.md`, `docs/C4_architecture.uml`): **monolito modular**, React SPA → API GraphQL (Go) → PostgreSQL. `docs/C4_architecture.uml` é o diagrama de componentes estilo C4; `docs/C4_architecture.pdf` é a versão renderizada.

## Requisitos funcionais obrigatórios (`docs/requirements.md` / `projeto.md`)

São requisitos avaliados, não sugestões — mantenha-os em mente em qualquer trabalho relacionado ao autocompletar:

- Sugestões só aparecem após **4+ caracteres** digitados.
- Sem correspondências, nada deve ser renderizado abaixo do campo de busca (nenhum elemento de estado vazio).
- O backend retorna **no máximo 20** sugestões; o frontend exibe **10**, as demais (11-20) acessíveis via scroll dentro da lista.
- O trecho correspondente ao termo buscado deve aparecer em **negrito** em cada sugestão.
- Hover (desktop) e touch (mobile) devem destacar a sugestão.
- As sugestões atualizam dinamicamente conforme o usuário digita, na velocidade da digitação (debounce, mas sem perder responsividade).
- Clicar em uma sugestão preenche o campo principal de busca com o texto dela.
- Deve ser responsivo para mobile.
- `docker compose up` deve subir toda a stack (frontend + backend + db) — ver `docker-compose.yml`.

## Processo específico deste desafio

- Toda mudança passa por Pull Request, mesmo trabalhando sozinho — nunca push direto na `main`.
- PRs devem mapear para as microtasks em `TASKS.md` sempre que possível (PRs pequenos e coerentes).
- `starter/suggestions.js` (se presente) é código de referência fornecido para avaliação, não uma base para construir sem revisão crítica — ver discussão esperada em `COMMENTS.md`.

### Convenção de nomenclatura de branch e PR

- Branches seguem `<tipo>/<descrição-curta>` (ex.: `chore/project-setup`, `docs/initial-setup`).
- Títulos de commit/PR seguem o padrão Conventional Commits — `<tipo>: <descrição>` — usando os tipos `chore`, `feat`, `test` e `docs` (mais `fix`/`refactor` quando fizer sentido).
- O corpo do PR referencia as issues resolvidas com `Closes #N` (uma ou mais por PR, agrupando issues relacionadas de uma mesma microtask).

## Manutenção do COMMENTS.md (uso de IA)

`COMMENTS.md` é o registro vivo de decisões arquiteturais e do uso de IA ao longo do projeto. **Sempre que você (Claude) participar de uma decisão relevante de código** — arquitetura, lógica não trivial, trechos que geraram dúvida — atualize `COMMENTS.md` você mesmo, como parte do trabalho, registrando:

- **O que foi pedido** — o prompt/pergunta feita a você;
- **O que foi aceito como veio** — e por quê fazia sentido;
- **O que foi alterado** — e o motivo da alteração;
- **O que foi rejeitado** — sugestões que você deu e que não foram usadas, e por quê.

Não é necessário registrar interações triviais (autocomplete de uma linha, correções óbvias de sintaxe). Registre principalmente decisões de arquitetura, lógica mais complexa, ou pontos em que o usuário corrigiu/ajustou o que você propôs. Faça essa atualização proativamente ao final de tarefas relevantes, sem esperar o usuário pedir.

## Estrutura do repositório

- `backend/` — API GraphQL em Go. `cmd/main.go` é o entrypoint; `internal/graphql`, `internal/repository`, `internal/service` são as camadas (resolvers GraphQL, acesso a dados, lógica de negócio).
- `frontend/` — SPA em React.
- `docs/` — documentos de arquitetura e requisitos (`requirements.md`, `C4_architecture.uml`/`.pdf`).
- `TASKS.md` — divisão do projeto em microtasks (setup & CI, backend, GraphQL, frontend, documentação).
- `COMMENTS.md` — decisões arquiteturais e registro de uso de IA (ver seção acima).
- `docker-compose.yml` — sobe `db` (Postgres 16), `api` (backend, porta 8080), `web` (frontend, porta 3000).

## Comandos

O CI (`.github/workflows/ci.yml`) define os comandos canônicos por componente; rode os mesmos localmente antes de dar push.

Backend (a partir de `backend/`):
```
go vet ./...
go test ./...
```

Frontend (a partir de `frontend/`):
```
npm install
npm run lint
npm test
```

Stack completa:
```
docker compose up
```
