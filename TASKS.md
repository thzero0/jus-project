# Escopo geral do projeto

1. Leitura detalhada do projeto, registrando os atributos e requisitos de qualidade exigidos
2. Escolher dominio das sugestoes
3. Design da arquitetura de software:
    - Escolha de padrao arquitetural 
    - Escolha das tecnologias
    - Diagramacao da arquitetura
4. Configuracao inicial do repo github (Criacao de issues para as micro-tasks)
5. Setup de ambiente e configuracao de CI
6. Desenvolvimento do Backend 
7. Desenvolvimento do GraphQL 
8. Desenvolvimento do FrontEnd 
9. Testes (diluidos ao longo do desenvolvimento)
10. Documentacao 

# Micro-tasks

## Setup & CI
1. Criar estrutura inicial de pastas & docker compose "esqueleto"
2. Configurar CI no Github Actions. Um job para cada componente.
3. Proteger a branch main.
4. Criar plano de desenvolvimento e arquivos necessarios para o uso do Claude Code ao longo do projeto.

## Backend
1. Definir modelo de dados e esquema de armazenamento
2. Escrever script/seed automatizado para popular o banco de dados
3. Implementar logica de armazenamento
4. Implementar logica de busca de sugestoes
5. Configurar lint do backend
6. Escrever testes unitarios 

## GraphQL
1. Definir esquema do GraphQL
2. Implementar conexao ao servico do backend
3. Testes 

## Frontend
1. Estruturar projeto React
2. Configurar lint do frontend
3. Criar componente de busca (sem o autocomplete ainda)
4. Integrar com o GraphQL
5. Incluir debounce
6. Renderizar lista de sugestoes 
    - negrito no trecho correspondente, scroll pros itens 11-20, hover highlight, e "nenhum elemento se nao houver sugestoes"
7. Implementar clique que preenche a palavra
8. Implementar e garantir responsividade mobile
8. Testar componentes

## Documentacao
1. Escrever revisao sobre o starter/suggestions.js
2. Escrever README.md com instrucoes de como rodar e quais dependencias sao necessarias
3. Consolidar COMMENTS.md com as decisoes arquiteturais + uso de IA
4. Teste final end-to-end