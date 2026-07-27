# Como testar a validade do backend

Guia rápido para rodar os testes já implementados em `backend/internal/service` — a camada de lógica de negócio do autocompletar (Trie em memória + regras de tamanho mínimo/máximo, ver `COMMENTS.md`).

## Rodando

A partir de `backend/`:

```bash
go vet ./...
go test ./... -v
```

Saída esperada:

```
=== RUN   TestMain_ServesHTTP
    main_test.go:26: DATABASE_URL not set; see docs/testing-backend.md to run this test
--- SKIP: TestMain_ServesHTTP (0.00s)
PASS
=== RUN   TestPostgresRepository_ListGames
    postgres_integration_test.go:23: DATABASE_URL not set; see docs/testing-backend.md to run this test
--- SKIP: TestPostgresRepository_ListGames (0.00s)
PASS
=== RUN   TestNewSuggestionService_RepositoryError
--- PASS: TestNewSuggestionService_RepositoryError (0.00s)
=== RUN   TestSearch_BelowMinLengthReturnsEmpty
--- PASS: TestSearch_BelowMinLengthReturnsEmpty (0.00s)
=== RUN   TestNewSuggestionService_Count
--- PASS: TestNewSuggestionService_Count (0.00s)
=== RUN   TestSearch_CapsAt20Results
--- PASS: TestSearch_CapsAt20Results (0.00s)
=== RUN   TestTrieSearch_PrefixMatch
--- PASS: TestTrieSearch_PrefixMatch (0.00s)
=== RUN   TestTrieSearch_CaseInsensitive
--- PASS: TestTrieSearch_CaseInsensitive (0.00s)
=== RUN   TestTrieSearch_NoMatchReturnsNil
--- PASS: TestTrieSearch_NoMatchReturnsNil (0.00s)
=== RUN   TestTrieSearch_RespectsLimit
--- PASS: TestTrieSearch_RespectsLimit (0.00s)
=== RUN   TestTrieSearch_AlphabeticalOrder
--- PASS: TestTrieSearch_AlphabeticalOrder (0.00s)
PASS
```

Os `SKIP` são esperados rodando localmente sem `DATABASE_URL` — ver "Testando contra dados reais" abaixo. No CI (`.github/workflows/ci.yml`), o job `backend` sobe um Postgres real (`services:`) e roda o seed antes do `go test`, então lá esses dois testes **rodam de verdade** em todo PR, sem `SKIP`.

Lint (mesma versão do CI, `.github/workflows/ci.yml`):

```bash
golangci-lint run ./...
# ou, sem instalar nada:
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v2.12.2 golangci-lint run ./...
```

Flags úteis: `-cover` (cobertura), `-run TestTrieSearch` (roda um subconjunto), `-count=1` (ignora cache).

## O que cada teste garante

**`trie_test.go`** — estrutura de dados isolada:
- `PrefixMatch` — retorna só nomes cujo prefixo bate
- `CaseInsensitive` — `"red "` encontra `"Red Dead Redemption"`
- `NoMatchReturnsNil` — prefixo sem correspondência retorna `nil`
- `RespectsLimit` — `limit` corta o resultado
- `AlphabeticalOrder` — resultado vem ordenado

**`service_test.go`** — regras de negócio, com `fakeRepository` (mock manual, sem lib):
- `RepositoryError` — erro do repositório é propagado, não engolido
- `BelowMinLengthReturnsEmpty` — termos com menos de 4 caracteres retornam vazio (requisito obrigatório)
- `Count` — reflete quantos jogos o repositório retornou no startup
- `CapsAt20Results` — nunca retorna mais que 20 resultados (requisito obrigatório)

**`postgres_integration_test.go`** (`internal/repository`) — fala com um Postgres de verdade, exercitando o mesmo caminho que `cmd/main.go` usa no startup (`NewPostgresRepository` → `ListGames`):
- confirma a contagem de jogos carregados após a limpeza/deduplicação do seed (8978)
- confirma que um jogo conhecido do dataset (`"The Elder Scrolls VI"`) está presente

**`main_test.go`** (`cmd`) — smoke test do binário como um todo: builda `cmd/main.go`, roda o binário como subprocesso de verdade (não em-processo, já que `main()` bloqueia em `ListenAndServe` e chama `log.Fatal` em erro), aponta pra uma porta livre via a env var `PORT`, espera o servidor responder, e faz uma query GraphQL de verdade em `POST /graphql`. É o único teste que exercita o processo inteiro — parsing de env, conexão com o banco, wiring do `SuggestionService`, e a camada GraphQL — de ponta a ponta.

**`resolver_test.go`** (`internal/graphql`) — testa só a tradução GraphQL → `SuggestionService`, com `fakeRepository` (sem Postgres, sem HTTP, roda em milissegundos):
- `ReturnsMatches` — `resolver.Query().Suggestions(ctx, term)` delega pro `Search()` do serviço e devolve o resultado
- `BelowMinLengthReturnsEmpty` — a regra de tamanho mínimo continua valendo passando pelo resolver

## Testando contra dados reais

`TestPostgresRepository_ListGames` e `TestMain_ServesHTTP` só rodam se a variável `DATABASE_URL` estiver setada — sem ela, são pulados (`SKIP`), então não quebram o CI (que não sobe Postgres pro job de backend). Pra rodar de verdade, contra os dados reais do `db/games.csv`:

```bash
# na raiz do repo — sobe o Postgres (porta 5432 exposta no host) e roda o seed
docker compose up -d db seed

export DATABASE_URL="postgres://admin:123@localhost:5432/games-db?sslmode=disable"
cd backend
go test ./... -v
```

Isso é também a forma mais direta de confirmar que o `cmd/main.go` real (não um fake) carrega o dataset: suba a stack inteira e cheque o log do serviço `api`:

```bash
docker compose logs api
# api-1  | .../... suggestion service ready: 8978 games loaded
# api-1  | .../... listening on :8080
```

## Escrevendo um novo teste

- Trie: use `newTestTrie(...)` para montar e `tr.search(prefixo, limite)` para consultar, comparando com `slices.Equal`.
- Regras de negócio: monte um serviço com `newTestService(t, games)` passando um `[]repository.Game` fake, chame `svc.Search(termo)`.
- Contra Postgres real ou subprocesso do binário: siga o padrão de `postgres_integration_test.go`/`main_test.go` — `t.Skip` se `DATABASE_URL` não estiver setada, pra não quebrar `go test ./...` sem banco disponível.

## Testando a camada GraphQL manualmente

Com `docker compose up` (ou só `db seed api`) rodando:

**Playground (navegador)** — abra `http://localhost:8080/`, uma IDE de GraphQL onde dá pra escrever a query, ver o schema (aba "Docs") e rodar direto:

```graphql
query {
  suggestions(term: "minecraf")
}
```

**curl** — o endpoint de queries é `POST /graphql`, corpo JSON com `query` e (opcional) `variables`:

```bash
curl -s http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"query($t: String!) { suggestions(term: $t) }","variables":{"t":"minecraf"}}'
```

```json
{"data":{"suggestions":["Minecraft","Minecraft: Pocket Edition","Minecraft: Story Mode","Minecraft: Story Mode — Season Two"]}}
```

Vale testar também os dois requisitos que vivem no `service` mas se manifestam pelo GraphQL: termo com menos de 4 caracteres (`"min"`) deve devolver `"suggestions":[]`; termo sem nenhum jogo correspondente (`"zzzz"`) também.

**CORS** — confere que a origem do `web` (`http://localhost:3000` por padrão, configurável via `CORS_ORIGIN`) recebe `Access-Control-Allow-Origin`, e que outras origens não recebem:

```bash
curl -s -i -X OPTIONS http://localhost:8080/graphql \
  -H 'Origin: http://localhost:3000' \
  -H 'Access-Control-Request-Method: POST' \
  -H 'Access-Control-Request-Headers: content-type' | grep -i access-control
```
