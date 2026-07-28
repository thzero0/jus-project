## Domínio da aplicação
Domínio escolhido: Top 10k jogos de vídeo games mais populares 2025. Com dados sendo retirados de um dataset do Kaggle: https://www.kaggle.com/datasets/onlypythondatasheet/10k-most-popular-gaming-2025?resource=download

## Projeto Arquitetural

### Síntese Arquitetural

Padrão arquitetural escolhido: Monolito Modular
- Dado o prazo de desenvolvimento, e a simplicidade do sistema, o padrão se justifica pela rapidez que proporciona no desenvolvimento sem gerar overhead de overengineering, ao mesmo tempo que separa logicamente as componentes envolvidas.

Modelo de representação escolhido: C4
- Possibilita difentes níveis da visão arquitetural, sendo suficientemente formal ao mesmo tempo que consegue atender difentes níveis de comunicação.

### Tecnologias escolhidas
Backend: Golang
- Pela familiaridade pessoal + eficiência da linguagem

Base de Dados: PostgreSQL
- Embora um JSON supra as necessidades atuais, o PostgreSQL é considerado pensando na eventual continuidade do projeto, em que novas funcionalidades poderiam ser incorporadas, tirando proveito dessa base de dados.


## Avaliação do `starter/suggestions.js`

O que eu manteria: a checagem de 4 caracteres antes de qualquer request, mesma lógica que usarei no cliente (antes do debounce) e no servidor (Como uma medida denfensiva e redundante).

O que eu faria diferente:
- O código não usa GraphQL na prática - `fetch('/graphql?q=...')` é um REST "disfarçado". Servidor GraphQL na realidade espera `{query, variables}`, não um `q` solto na URL. Por isso usaria um cliente de verdade (`graphql-request`) mandando a query certa pro schema.
- Não há cancelamento de requisição. Respostas antiga podem chegar depois da nova e sobrescrever a tela com dado velho.
- Não há debounce, cada tecla dispara uma requisição ao backend, como o app é leve não impacta tanto. Colocaria um debounce rápido o suficiente pra parecer instantâneo, mas que evite as multiplas requisições ao digitar rápido.
- A parte de `.slice(0, 10)` descarta os itens 11-20 pra sempre, isso quebraria o requisito de scroll. Faria o corte ser visual (CSS), não na lógica de busca.

O que eu mudaria antes de produção: Não existe tratamento de erro nenhum, se o fetch falhar, a resposta não for OK, ou o corpo não for o JSON esperado, a UI quebra sem nenhum tipo de "feedback".

## Antes de integrar o Claude ao repositório (usando interface web)

- **O que pedi ao Claude:** A partir dos documentos TASKS.md, COMMENTS.md, docs/requirements.md e C4_architecture.md, gerar as issues correspondentes a cada microtask (anexei os documentos).
- **O que aceitei como veio:** As labels que o Claude sugeriu; o uso do Vite pra estruturar o projeto React.
- **O que alterei:** As descrições de cada issue geradas pelo Claude, cortando prolixidade e caminhos que soavam a "overengineering".
- **O que rejeitei:** Algumas issues que o Claude sugeriu e que já estavam resolvidas (ex.: síntese arquitetural, proteção da branch main).

## Estratégia de busca por prefixo: Trie em memória

- **O que pedi ao Claude:** Avaliar minha ideia de, em vez de rodar uma query SQL `LIKE 'prefixo%'` a cada busca, carregar os dados do Postgres uma única vez (no startup do backend) e montar uma Trie em memória, respondendo buscas de prefixo direto dessa estrutura.
- **O que aceitei como veio:** A validação do Claude sobre o raciocínio central — o dataset é pequeno (10k jogos) e estático entre execuções, então uma busca em Trie (O(k) no tamanho do prefixo) é bem mais rápida que uma ida ao Postgres por tecla digitada, atacando diretamente o atributo de qualidade "baixíssima latência".
- **O que alterei:** Nada na ideia em si — só confirmei com o Claude o escopo: a Trie vive só no processo Go (camada `service`); o front-end continua chamando o GraphQL normalmente, sem receber a Trie/dataset diretamente (o que quebraria o requisito de o front-end se comunicar só via GraphQL).
- **O que rejeitei:** N/A.

**Trade-off registrado:** a Trie é montada uma vez no startup a partir do Postgres; se os dados forem alterados diretamente no banco (fora do fluxo de seed), a Trie fica desatualizada até o processo reiniciar. Aceitável aqui porque o dataset só muda via seed automatizado.

## Limpeza de dados no seed: `backend/db/clean_data.py`

- **O que pedi ao Claude:** Integrar ao seed um passo de "clean data" em Python, rodando antes da carga no Postgres — padronizar os nomes dos jogos removendo acentos e deixando as iniciais de cada palavra maiúsculas, sem remover hífen.
- **O que aceitei como veio:** A estrutura geral do script proposta pelo Claude (ler `games.csv`, aplicar a limpeza, descartar duplicatas por `id` mantendo a primeira ocorrência, escrever um CSV limpo consumido pelo `\copy`).
- **O que alterei:** A primeira versão do Claude usava `unicodedata.normalize("NFKD", ...)` pra tirar acento — eu corrigi, apontando que NFKD faz decomposição de *compatibilidade* (mexe também em ligaduras, símbolos e formas full-width, não só em acentos), enquanto NFD faz decomposição *canônica*, afetando só os diacríticos; pedi a troca pra NFD. Depois o Claude tentou capitalizar com `name.lower().title()` — o próprio `.title()` do Python quebra em qualquer caractere não alfabético, então nomes com apóstrofo saíam errados (ex.: `"don't fall behind"` virava `"Don'T Fall Behind"`); pedi a troca por uma capitalização manual que só sobe a primeira letra de cada palavra separada por espaço, sem tocar no resto (preserva `Don't`, `Flumpty's`, hífen, siglas etc.).
- **O que rejeitei:** Uma variante intermediária que o Claude propôs com tabela de tradução manual (`str.maketrans` mapeando cada caractere acentuado pro equivalente sem acento) — pedi NFD no lugar, por ser mais simples de manter do que enumerar caractere por caractere.

## Implementação da Trie em Go: simplicidade vs. otimização (`backend/internal/service/trie.go`)

- **O que pedi ao Claude:** Implementar a lógica de busca por prefixo na camada `service` (sem dependência de GraphQL — ver decisão acima).
- **O que aceitei como veio:** A versão que o Claude propôs em que cada nó guarda, em `matches []string`, todos os nomes que passam por ele (preenchido no `insert`); a busca desce k nós até o nó do prefixo e já tem a lista pronta, ordenando com `sort.Strings` a cada chamada de `search`.
- **O que alterei:** O Claude tinha implementado antes uma variante mais sofisticada — DFS clássico a partir do nó final de cada nome, com os filhos de cada nó pré-ordenados uma única vez no startup em vez de ordenar a cada busca. Identifiquei que essa complexidade extra não se justificava para um dataset estático de ~10k nomes (a diferença de performance entre as duas abordagens é desprezível) e pedi pra voltar pra versão mais simples.
- **O que rejeitei:** N/A — a variante com DFS e pré-ordenação não chegou a ser commitada; foi substituída antes disso.

**Trade-off registrado:** `search()` refaz `sort.Strings` a cada chamada em vez de ordenar uma única vez no startup. Aceitável porque o dataset é pequeno e estático (10k nomes) — o custo por tecla digitada é irrelevante na prática, e a simplicidade do código pesa mais que o ganho de performance nessa escala.

## Camada GraphQL: escolha de biblioteca e wiring (`backend/internal/graphql/`, `backend/cmd/main.go`)

- **O que pedi ao Claude:** Implementar a camada GraphQL sobre o `SuggestionService` já existente, explicando trade-offs ao longo do caminho.
- **O que aceitei como veio:** `gqlgen` (schema-first, geração de código, type-safe em compile-time) em vez de `graph-gophers/graphql-go` (sem codegen, erros só em runtime) — troquei segurança em compile-time pelo custo de um `go generate` a mais. Schema minimalista: uma query `suggestions(term: String!): [String!]!`, sem tipo customizado.
- **O que alterei:** Nada diretamente — mas o Claude, testando de verdade (não só lendo o código), encontrou e corrigiu sozinho: (1) `gqlgen` mais recente exigia Go 1.25, incompatível com o 1.24 fixado no CI/`Dockerfile` — pinou `v0.17.80`, compatível com 1.24, em vez de subir a versão do projeto; (2) um campo `Resolver.Suggestions` com o mesmo nome do método gerado `Suggestions()` quebraria a chamada (Go prioriza método declarado sobre campo promovido por embedding) — renomeou pra `SuggestionService` antes de compilar; (3) um teste assumia que `"the e"` traria `"The Elder Scrolls VI"` no top 20 — rodando contra o Postgres real, falhou: dezenas de outros títulos "The Elder Scrolls ..." empurram esse pra fora do cap alfabético. O backend estava certo, o teste que errou; ajustou o prefixo pra `"the elder scrolls vi"`.
- **O que rejeitei:** N/A.

**Trade-off registrado:** GraphQL Playground montado em `/`, queries em `/graphql` — troca o placeholder "ok" por uma IDE utilizável no navegador; útil pra avaliação, seria superfície a mais pra revisar em produção real.

## CORS entre `web` e `api` (`backend/cmd/main.go`, `docker-compose.yml`)

- **O que pedi ao Claude:** Perguntei se faltava algo em backend + GraphQL antes de considerar essa parte pronta.
- **O que aceitei como veio:** O Claude identificou (e confirmou testando `curl` com preflight `OPTIONS` de verdade contra o container) que `api` (porta 8080) e `web` (porta 3000) são origens diferentes pro navegador, e o servidor não mandava nenhum cabeçalho CORS — o fetch do React pro `/graphql` seria bloqueado assim que o frontend existisse. Aceitei corrigir agora, com `github.com/rs/cors` (em vez de escrever os cabeçalhos na mão — CORS tem detalhes fáceis de errar, como o `Vary` e a resposta ao preflight) restringindo a origem via env var `CORS_ORIGIN` (default `http://localhost:3000`, setada explicitamente no `docker-compose.yml`).
- **O que alterei:** N/A.
- **O que rejeitei:** A alternativa que o Claude também levantou — proxy reverso no nginx do `web` pro `api`, evitando CORS por completo — preferi CORS no Go mesmo, por manter o frontend sem precisar conhecer o roteamento da API.

## Setup do frontend React (`frontend/`)

- **O que pedi ao Claude:** Implementar só o setup do projeto React, sem a busca ainda.
- **O que aceitei como veio:** TypeScript (Vite já tinha sido decidido antes) pela tipagem em compile-time, consistente com o rigor do backend em Go. ESLint em vez do `oxlint`, que o `create-vite` mais recente já traz por padrão — preferi o padrão mais estabelecido, que qualquer revisor reconhece de cara. Vitest + Testing Library pro `npm test`, com só um smoke test de renderização por ora.
- **O que alterei:** N/A.
- **O que rejeitei:** N/A.

**Trade-off registrado:** o TypeScript aqui só tipa o interior dos componentes React — ainda não tem nenhuma geração de tipos a partir do `schema.graphqls`, então uma mudança no schema do backend não seria pega automaticamente pelo compilador do frontend. Fica pra quando a chamada GraphQL de verdade for implementada.

## Debounce e cancelamento de requisições obsoletas (`frontend/src/App.tsx`)

- **O que pedi ao Claude:** Debounce na busca — sem requisição abaixo de 4 caracteres, calibrado pra parecer instantâneo, descartando respostas obsoletas.
- **O que aceitei como veio:** `AbortController` pra cancelar a requisição anterior quando o termo muda, em vez de um contador de sequência.
- **O que alterei:** Editei `DEBOUNCE_MS` pra 3s pra testar visualmente e não vi diferença — o Claude explicou que eu estava olhando um container Docker com build antigo (edição em disco não atualiza imagem já buildada). Em vez de escolher entre as duas opções que ele deu (`npm run dev` ou rebuild pontual), pedi um jeito de `docker compose up` sempre refletir o código atual.
- **O que rejeitei:** O teste automatizado que o Claude propôs pra isso (`App.debounce.test.tsx`, com fake timers) — não achei necessário.

## Responsividade mobile (`frontend/src/App.module.css`, `SuggestionList.module.css`)

- **O que pedi ao Claude:** Campo de busca e lista exibidos corretamente em telas pequenas, testado num viewport de dispositivo móvel real.
- **O que aceitei como veio:** A lista virou um dropdown ancorado (`position: absolute` sob o campo) em vez de empurrar o resto da página — evita o layout "pulando" quando a lista abre/fecha, o que importa ainda mais com o teclado virtual ocupando parte da tela.
- **O que alterei:** N/A — mas o Claude errou e precisou ser corrigido: a primeira tentativa de alinhar verticalmente o texto dos itens usou `display: flex`, que quebrou o negrito e o resto do texto em colunas separadas em vez de fluir como texto normal. Só apareceu testando de verdade num viewport estreito (títulos longos quebram em duas linhas com mais frequência ali) — não teria aparecido só lendo o CSS. O Claude trocou por padding vertical simétrico, que resolve sem mexer no fluxo do texto.
- **O que rejeitei:** N/A.

## Escopo e organização dos testes unitários (`service_test.go` / `trie_test.go`)

- **O que pedi ao Claude:** Implementar a lógica de busca de sugestões com uma interface pronta pra mock em testes.
- **O que aceitei como veio:** O `fakeRepository` escrito à mão (sem lib de mock) e funções de teste separadas por cenário, em vez de tabela única — ambos sugeridos pelo Claude.
- **O que alterei:** O Claude escreveu os testes todos misturados em `service_test.go`, testando a Trie e as regras de negócio (min. 4 caracteres, máx. 20) no mesmo lugar. Também pedi pra separar as responsabilidades: um corte entre `trie_test.go` (estrutura de dados isolada — prefixo, case-insensitive, limite, ordem alfabética) e `service_test.go` (regras de negócio — tamanho mínimo, cap de 20, erro do repositório).
- **O que rejeitei:** N/A.
