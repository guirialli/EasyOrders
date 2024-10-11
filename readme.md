## Desafio: Listagem de Orders

Tecnologias Utilizadas
REST API: Endpoint GET /order
gRPC: Serviço ListOrders
GraphQL: Query ListOrders
Banco de Dados: MySQL com Docker
Arquivo api.http: Para facilitar a criação e listagem das orders via requests.
Pré-requisitos
Docker e Docker Compose
Go
MySQL (via Docker)
Configurações de Ambiente
Crie um arquivo .env na raiz do projeto com o seguinte conteúdo:

```dotenv
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
WEB_SERVER_PORT=:8081
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8080
```
Setup do Projeto
1. Clone o repositório
   ```bash
   git clone https://github.com/guirialli/EasyOrders.git
   cd EasyOrders
   ```
2. Subir o banco de dados com Docker
   O Docker Compose será usado para subir um container com o MySQL configurado.

Execute o comando abaixo para subir o banco de dados:
```bash
docker compose up
```
3. Executar o projeto
   Para rodar a aplicação, use o comando:

```bash
go run cmd/ordersystem/main.go wire_gen.go
```
Endpoints
1. REST API (GET /order)
   Método: GET
   URL: http://localhost:8081/order
2. gRPC (ListOrders Service)
   Serviço gRPC rodando na porta 50051. Para testar o serviço, use um cliente gRPC como o Evans ou outro de sua preferência.
3. GraphQL Query
   A query ListOrders pode ser executada no servidor GraphQL, rodando na porta 8080.
   Exemplo de query:

```graphql
query {
   listOrders {
      id
      price
      tax
      final_price
   }
}
```
Arquivo api.http
O projeto inclui um arquivo api.http com exemplos de requests para criar e listar orders. Utilize-o em clientes HTTP como o VSCode ou Insomnia.
