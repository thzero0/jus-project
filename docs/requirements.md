# Requisitos de projeto explícitos no texto
### Estrutural arquitetural
- Uma **única página**
- Front-end em **ReactJS**
- Front-end se comunica com GraphQL, GraphQL se comunica com Back-end
- **Responsividade mobile**

### Sugestões do autocompletar
- Sugestões só começam a ser exibidas após digitar, no mínimo, **4 caracteres**;
- Se não existem sugestões para o termo digitado, nenhum elemento deve ser exibido abaixo do campo de busca; 
- Retornar no **máximo 20 sugestões, mas exibir no máximo 10**, as outras sendo acessadas por scroll
- As sugestões precisam **manter em negrito** a parte do termo que corresponde ao termo inicial da busca; 
- Ao passar o mouse por alguma sugestão (hover) ou tocar no elemento (touch no mobile), o elemento deve ser destacado; 
- O usuário pode continuar digitando — as sugestões vão **mudando dinamicamente** conforme o termo é atualizado; 
- As sugestões precisam ser exibidas **na mesma (ou próxima) velocidade que o usuário digita**; 
- Ao clicar em alguma sugestão, o campo de busca principal deve ser atualizado com o texto da sugestão.

# Atributos de qualidade desejáveis (não explícitos)
- Sistema leve, baixo consumo de recursos dos usuários;
- Desacoplamento das componentes
- Interface simples, intuitiva e esteticamente agradável
- Baixíssima latência
