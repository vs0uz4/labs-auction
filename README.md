# Auction

> [!IMPORTANT]
> Para poder executar o projeto contido neste repositório é necessário que se tenha o Go instalado no computador. Para maiores informações siga o site <https://go.dev/>

- [Auction](#auction)
  - [Desafio Golang Pós GoExpert - Labs Auction](#desafio-golang-pós-goexpert---labs-auction)
  - [Requisitos a Serem Seguidos](#requisitos-a-serem-seguidos)
  - [Entrega](#entrega)
  - [Extras](#extras)
    - [Usuários Default\`s](#usuários-defaults)
  - [Executando a Aplicação](#executando-a-aplicação)
    - [Subindo a Aplicação](#subindo-a-aplicação)
    - [Encerrando a Aplicação](#encerrando-a-aplicação)
  - [Leilões (`auctions`)](#leilões-auctions)
    - [Criando um Leilão](#criando-um-leilão)
    - [Consultando Leilões](#consultando-leilões)
      - [Consultando leilão por ID](#consultando-leilão-por-id)
      - [Consultando leilão por Status, Category ou ProductName](#consultando-leilão-por-status-category-ou-productname)
  - [Lances (`bids`)](#lances-bids)
    - [Dando um Lance em um item](#dando-um-lance-em-um-item)
    - [Consultando Lances](#consultando-lances)
      - [Consultando lances pelo ID do leilão](#consultando-lances-pelo-id-do-leilão)
      - [Consultando lance vencedor do leilão](#consultando-lance-vencedor-do-leilão)
  - [Usuários (`users`)](#usuários-users)

## Desafio Golang Pós GoExpert - Labs Auction

Este projeto faz parte da Pós GoExpert como desafio, nele são cobertos os conhecimentos em APIRest, channels, tratamentos de erros, packages, Clean Architecture, DI, Banco de Dados, Go Rotines, Multithreading

**Objetivo**: Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo determinado.

## Requisitos a Serem Seguidos

- Clonar o seguinte repositório: https://github.com/devfullcycle/labs-auction-goexpert;
- Adicionar a rotina de fechamento automático a partir de um determinado tempo;
- Utilizar `go routines` para a implementação da rotina de fechamento automático;

**Nós Devemos Desenvolver**:

- Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
- Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
- Um teste para validar se o fechamento está acontecendo de forma automatizada;

> [!TIP]
> Algumas dicas para ajudar no desenvolvimento

- Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
- Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
- Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
- Para mais informações de como funciona uma `goroutine`, foque em nosso módulo de Multithreading no curso Go Expert;

## Entrega

- O código-fonte completo da implementação;
- Documentação explicando como rodar o projeto em ambiente dev;
- Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.

## Extras

Como ao clonar o repositório, constatei que algumas funcionalidades estavam com problemas, então decidi realizar a correção das mesmas para que todas as funcionalidades já presentes no projeto estejam funcionais. Abaixo,
segue uma lista de algumas correções e adições realizadas no projeto com o intúito de melhorar o mesmo.

> [!WARNING]
> A listagem de correções e implementações abaixo não faziam parte do enunciado do desafio e nem tão pouco eram obrigatórias!

- Correção do `FindBidByAuctionId` no respositório de Bids, fazendo com que seja possível pesquisar os lances por `auctionId`;
- Correção do `FindAuctions` no repositório de Auctions, fazendo com que seja possível pesquisar também por `productName`;
- Implementação de camada de tratamento para `productName` evitando que caracteres especiais possam ser interpretados pelo Mongo como operadores `regex`;
- Configurações dos `paths` de `output` para os Logs do framework ZapCore, de forma que os mesmos sejam direcionados para Std;
- Adição de `script` de inicialização do MongoDB criando `collections`, `indíces` e `seeds` de usuários defaults para testar a aplicação;
- Adição de `health-check` para o serviço do MongoDB no docker-compose.yml e configuração de `dependency` para a aplicação.

### Usuários Default`s

Os dados dos usuários que serão adicionados automaticamente ao MongoDB através do `seed` implementado, são os seguintes:

```json
{
  _id: "d290f1ee-6c54-4b01-90e6-d701748f0851",
  name: "John Doe",
},
{
  _id: "93fb1e9c-523f-4d92-80b4-0f7ba12fef56",
  name: "Jane Smith",
},
{
  _id: "4be43d3d-5f47-4881-a07b-8b5d3c5296c1",
  name: "Alice Johnson",
},
```

> Utilize os `id's` destes usuários para poder criar os lances(`bids`) do leilão.

## Executando a Aplicação

Para executar a aplicação, primeiramente você deverá se certificar de atender ao pré-requisito básico que é ter o `Docker` instalado em sua máquina, desta forma você poderá rodar a aplicação sem a necessidade de instalar mais nada em sua máquina.

### Subindo a Aplicação

Estando dentro do diretório da aplicação, de modo que você já tenha `clonado` o repositório, basta executar os seguintes comandos abaixo:

```shell
❯ docker compose up -d --build
```

Na janela do terminal, você deverá ver uma mensagem parecida com o exemplo abaixo quando:

```shell
[+] Running 2/2
 ✔ Container mongodb     Healthy   0.5s 
 ✔ Container auctionapp  Started   0.7s
```

> A aplicação estará disponível no seguinte endereço `http://localhost:8080`

### Encerrando a Aplicação

Para encerrar a aplicação estando na janela do terminal onde iniciamos a mesma, bastar executarmos o seguinte comando abaixo:

```shell
❯ docker compose stop
```

Na janela do terminal, você deverá ver uma mensagem parecida com o exemplo abaixo quando:

```shell
[+] Stopping 2/2
 ✔ Container auctionapp  Stopped  0.1s 
 ✔ Container mongodb     Stopped  0.2s
```

## Leilões (`auctions`)

### Criando um Leilão

Para cadastrar um novo item para leiloar, deverá ser realizada uma requisição do tipo `POST` para o seguinte endereço `http://localhost:8080/auction` contendo em seu corpo um JSON parecido com o abaixo:

```json
{
    "product_name": "Kit de Chave de Fenda Mijia Wiha Para uso Diário",
    "category": "Ferramentas",
    "description": "Xiaomi Wiha Chave de Fenda de Precisão 24 em 1 - Modelo: JXLSD01XH",
    "condition": 1
}
```

> Em um cenário de sucesso, você deverá receber uma resposta sem `conteúdo` com um **HTTP Status Code** igual à **201**.

> [!NOTE]
> Lembrando que `condition` pode variar entre `1` à `3` sendo respectivamente seus valores: \
> 1 - Novo; \
> 2 - Usado; \
> 3 - Recondicionado

### Consultando Leilões

Temos duas formas de consultar os leilões, sendo elas as seguintes:

- Consultar leilão por ID;
- Consultar leilão por Status, Category ou ProductName.

#### Consultando leilão por ID

Para consultar um leilão por seu ID, deverá ser realizada uma requisição do tipo `GET` para o seguinte endereço `http://localhost:8080/auction/:auctionId` onde `:auctionId` deverá ser substituido pelo ID do leilão a ser consultado, conforme o exemplo abaixo:

```http
http://localhost:8080/auction/d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa
```

> Em um cenário de sucesso, você deverá receber um **HTTP Status Code** igual a **200** e um conteúdo precido com o seguinte abaixo:

```json
{
  "id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
  "product_name": "Kit de Chave de Fenda Mijia Wiha Para uso Diário",
  "category": "Ferramentas",
  "description": "Xiaomi Wiha Chave de Fenda de Precisão 24 em 1 - Modelo: JXLSD01XH",
  "condition": 0,
  "status": 0,
  "timestamp": "2024-12-19T18:40:48-03:00"
}
```

#### Consultando leilão por Status, Category ou ProductName

Para consultar um leilão por seu ID, deverá ser realizada uma requisição do tipo `GET` para o seguinte endereço `http://localhost:8080/auction` onde os parâmetros de busca deverão ser adicionados a URL como `queryParams`, por exemplo: Em uma consulta onde queremos consultar os leilões com `status` equivalente à **0** e `productName` contenha o termo **Mijia**, o endereço ficaria conforme abaixo:

```http
http://localhost:8080/auction?status=0&productName=Mijia
```

> Em um cenário de sucesso, você deverá receber um **HTTP Status Code** igual a **200** e um conteúdo precido com o seguinte abaixo:

```json
[
 {
  "id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
  "product_name": "Kit de Chave de Fenda Mijia Wiha Para uso Diário",
  "category": "Ferramentas",
  "description": "Xiaomi Wiha Chave de Fenda de Precisão 24 em 1 - Modelo: JXLSD01XH",
  "condition": 0,
  "status": 0,
  "timestamp": "2024-12-19T18:40:48-03:00"
 },
 {
  "id": "e0684cea-0889-459e-9ea6-c0aa04d6ebf4",
  "product_name": "Câmera Ip Xiaomi Mijia Wifi",
  "category": "Segurança",
  "description": "Camera de Seguranca IP Xiaomi 360? Hd 1080p Wifi",
  "condition": 0,
  "status": 0,
  "timestamp": "2024-12-19T18:41:03-03:00"
 },
 {
  "id": "25bb30d6-89bd-452c-9377-6d9249634074",
  "product_name": "Caneta Stylus Xiaomi Mijia - 0,5mm",
  "category": "Escritório",
  "description": "Caneta de tinta em gel Xiaomi Mijia",
  "condition": 0,
  "status": 0,
  "timestamp": "2024-12-19T18:41:08-03:00"
 }
]
```

## Lances (`bids`)

### Dando um Lance em um item

Para dar um lance em um item sendo leiloado, deverá ser realizada um requisição do tipo `POST` para o seguinte endereço `http://localhost:8080/bid` contendo em seu corpo um JSON parecido com o abaixo:

```json
{
 "user_id": "93fb1e9c-523f-4d92-80b4-0f7ba12fef56",
 "auction_id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
 "amount": 300
}
```

Onde:

- user_id : Deverá ser equivalente a um dos `id's` dos usuários disponibilizados nesta documentação na seção [Usuários Default\`s](#usuários-defaults)
- auction_id : Deverá ser equivalente a um `id` de um item adicionado para leilão;
- amount : Deverá ser um `float` representando o valor do lance a ser dado.

> Em um cenário de sucesso, você deverá receber uma resposta sem `conteúdo` com um **HTTP Status Code** igual à **201**.

### Consultando Lances

Temos duas formas de consultar os lances, sendo elas as seguintes:

- Consultar lances pelo ID do leilão;
- Consultar lance vencedor pelo ID do leilão.

#### Consultando lances pelo ID do leilão

Para consultar os lances de um determinado leilão, deverá ser realizada uma requisição do tipo `GET` para o seguinte endereço `http://localhost:8080/bid/:auctionId` onde `:auctionId` deverá ser substituido pelo ID do leilão a ser consultado os lances, conforme o exemplo abaixo:

```http
http://localhost:8080/bid/d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa
```

> Em um cenário de sucesso, você deverá receber um **HTTP Status Code** igual a **200** e um conteúdo precido com o seguinte abaixo:

```json
[
 {
  "id": "cf7b0f4f-1c06-46e9-9324-03cf356ecd6b",
  "user_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
  "auction_id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
  "amount": 150.45,
  "timestamp": "2024-12-19T21:15:31-03:00"
 },
 {
  "id": "cdb608c6-a833-4df7-bb99-231682c2d0a4",
  "user_id": "4be43d3d-5f47-4881-a07b-8b5d3c5296c1",
  "auction_id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
  "amount": 215,
  "timestamp": "2024-12-19T21:23:41-03:00"
 },
 {
  "id": "cebef59f-3226-4426-b8c4-e764706c2edc",
  "user_id": "93fb1e9c-523f-4d92-80b4-0f7ba12fef56",
  "auction_id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
  "amount": 300,
  "timestamp": "2024-12-19T21:31:15-03:00"
 }
]
```

#### Consultando lance vencedor do leilão

Para consultar o lance vencedor de um leilão, deverá ser realizada uma requisição do tipo `GET` para o seguinte endereço `http://localhost:8080/auction/winner/:auctionId` onde `:auctionId` deverá ser substituido pelo ID do leilão a ser consultado o lance vencedor, conforme o exemplo abaixo:

```http
http://localhost:8080/auction/winner/3b2dcc69-5ca1-4e39-9fe2-7c08ab7588b6
```

> Em um cenário de sucesso, você deverá receber um **HTTP Status Code** igual a **200** e um conteúdo precido com o seguinte abaixo:

```json
{
  "auction": {
  "id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
  "product_name": "Kit de Chave de Fenda Mijia Wiha Para uso Diário",
  "category": "Ferramentas",
  "description": "Xiaomi Wiha Chave de Fenda de Precisão 24 em 1 - Modelo: JXLSD01XH",
  "condition": 1,
  "status": 0,
  "timestamp": "2024-12-19T18:40:48-03:00"
 },
 "bid": {
  "id": "cebef59f-3226-4426-b8c4-e764706c2edc",
  "user_id": "93fb1e9c-523f-4d92-80b4-0f7ba12fef56",
  "auction_id": "d0d0fe5c-5c7c-400c-ad6d-75d8f5daa9aa",
  "amount": 300,
  "timestamp": "2024-12-19T21:31:15-03:00"
 }
}
```

## Usuários (`users`)

Para consultar as informações de um usuário, deverá ser realizada uma requisição do tipo `GET` para o seguinte endereço `http://localhost:8080/user/:userId` onde `:userId` deverá ser substituido pelo ID do usuário a ser consultado, conforme o exemplo abaixo:

```http
http://localhost:8080/user/d290f1ee-6c54-4b01-90e6-d701748f0851
```

> Em um cenário de sucesso, você deverá receber um **HTTP Status Code** igual a **200** e um conteúdo precido com o seguinte abaixo:

```json
{
 "id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
 "name": "John Doe"
}
```