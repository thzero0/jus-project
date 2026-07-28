# Justarter

Busca com autocompletar de jogos — React (frontend) → GraphQL (backend em Go) → PostgreSQL.

## Domínio

Top 10k jogos mais populares de 2025, a partir de um [dataset do Kaggle](https://www.kaggle.com/datasets/onlypythondatasheet/10k-most-popular-gaming-2025). Os dados são limpos (acentos, capitalização, duplicatas) e versionados no repositório, carregados no Postgres por um seed automatizado.

## Stack e dependências

**Backend** (`backend/`, Go 1.24)
- [`gqlgen`](https://github.com/99designs/gqlgen) — servidor GraphQL (schema-first, geração de código)
- [`lib/pq`](https://github.com/lib/pq) — driver Postgres
- [`rs/cors`](https://github.com/rs/cors) — CORS
- PostgreSQL 16

**Frontend** (`frontend/`, Node 20)
- React 19 + TypeScript, via Vite
- [`graphql-request`](https://github.com/jasonkuhrt/graphql-request) — cliente GraphQL
- Vitest + Testing Library — testes
- ESLint — lint

**Infra**
- Docker + Docker Compose (sobe `db`, `seed`, `api`, `web`)

## Como rodar

Pré-requisito: Docker + Docker Compose.

```bash
docker compose up --build
```

- Frontend: http://localhost:3000
- API GraphQL: http://localhost:8080/graphql (Playground em http://localhost:8080)
- Postgres: `localhost:5432` (`admin`/`123`, banco `games-db`)

Sempre use `--build` após alterar código — sem ele o Compose reaproveita a imagem já buildada.

Docker + Docker Compose é tudo que é preciso ter instalado — Go, Node e até o Python do seed rodam dentro das próprias imagens (`golang:1.24-alpine`, `node:20-alpine`, `postgres:16-alpine` + `python3`), nunca no host.

### Rodando local (sem Docker), pra desenvolver

Só necessário se for editar código com hot-reload em vez de rebuildar a imagem a cada mudança — nesse caso sim, precisa de Go e Node instalados na máquina.

Backend:
```bash
cd backend
go run ./cmd   # precisa de DATABASE_URL apontando pro Postgres (docker compose up -d db seed)
```

Frontend:
```bash
cd frontend
npm install
npm run dev    # usa VITE_GRAPHQL_URL de .env (default: localhost:8080/graphql)
```

## Testes

```bash
# backend
cd backend && go vet ./... && go test ./...

# frontend
cd frontend && npm run lint && npm test && npm run build
```

Detalhes de cada camada de teste (unitário, integração com Postgres real, smoke do binário) em [`docs/testing-backend.md`](docs/testing-backend.md).

## Mais documentação

- [`TASKS.md`](TASKS.md) — divisão do projeto em microtasks
- [`COMMENTS.md`](COMMENTS.md) — decisões arquiteturais e uso de IA
- [`docs/requirements.md`](docs/requirements.md) — requisitos extraídos do desafio
- [`docs/C4_architecture.pdf`](docs/C4_architecture.pdf) — diagrama de arquitetura (C4)
