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


## Antes de integrar o Claude ao repositório (usando interface web)

- **O que foi pedido:** A partir dos documentos TASKS.md, COMMENTS.md, docs/requirements.md, e C4_architecture.md gere as issues correspondentes de cada microtask. (Anexei os documentos correspondentes.)
- **O que você aceitou como veio:** Labels sugeridas; Vite para estruturar o projeto React; 
- **O que você alterou:** Descrições de cada issue, cortando prolixidade e caminhos "overengineering"
- **O que você rejeitou:** Algumas issues que já haviam sido resolvidas (Como a sintese arquitetural, proteção da branch main, etc)

## Estratégia de busca por prefixo: Trie em memória

- **O que foi pedido:** Avaliação da ideia de, em vez de rodar uma query SQL `LIKE 'prefixo%'` a cada busca, carregar os dados do Postgres uma única vez (no startup do backend) e montar uma árvore Trie em memória, respondendo as buscas de prefixo direto dessa estrutura.
- **O que você aceitou como veio:** A ideia central — o dataset é pequeno (10k jogos) e estático entre execuções, então uma busca em Trie (O(k) no tamanho do prefixo) é bem mais rápida que ida ao Postgres por tecla digitada, atacando diretamente o atributo de qualidade "baixíssima latência".
- **O que você alterou:** Nada na ideia em si; só confirmei o escopo — o Trie vive apenas no processo Go (camada `service`), o front-end continua chamando o GraphQL por busca normalmente. Não é o caso de embarcar o dataset/Trie no cliente, o que quebraria o requisito de o front-end se comunicar obrigatoriamente com o GraphQL.
- **O que você rejeitou:** N/A.

**Trade-off registrado:** o Trie é montado uma vez no startup a partir do Postgres; se os dados forem alterados diretamente no banco (fora do fluxo de seed), o Trie fica desatualizado até o processo reiniciar. Aceitável aqui porque o dataset só muda via seed automatizado.

## Limpeza de dados no seed: `backend/db/clean_data.py`

- **O que foi pedido:** Integrar ao seed um passo de "clean data" em Python, rodando antes da carga no Postgres: padronizar os nomes dos jogos removendo acentos e deixando as iniciais de cada palavra maiúsculas, sem remover hífen.
- **O que você aceitou como veio:** A estrutura geral do script (ler `games.csv`, aplicar a limpeza, descartar duplicatas por `id` mantendo a primeira ocorrência, escrever um CSV limpo consumido pelo `\copy`).
- **O que você alterou:** Minha primeira versão usava `unicodedata.normalize("NFKD", ...)` para tirar acento — **o usuário me corrigiu**, apontando que NFKD faz decomposição de *compatibilidade* (mexe também em ligaduras, símbolos e formas full-width, não só em acentos), enquanto NFD faz decomposição *canônica*, afetando somente os diacríticos. Troquei para NFD. Depois, tentei capitalizar com `name.lower().title()`; o próprio `.title()` do Python quebra em qualquer caractere não alfabético, então nomes com apóstrofo saíam errados (ex.: `"don't fall behind"` virava `"Don'T Fall Behind"`). Troquei por uma capitalização manual que só sobe a primeira letra de cada palavra separada por espaço, sem tocar no resto da palavra (preserva `Don't`, `Flumpty's`, hífen, siglas etc.).
- **O que você rejeitou:** Uma variante intermediária com tabela de tradução manual (`str.maketrans` mapeando cada caractere acentuado para o equivalente sem acento) — o usuário pediu para usar NFD em vez disso, por ser mais simples de manter do que enumerar caractere por caractere.
