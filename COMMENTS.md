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