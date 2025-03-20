##### Projeto Braip - Teste Técnico

Este projeto é uma API RESTful desenvolvida em Go (Golang) para gerenciar produtos. Ele inclui operações CRUD, consultas personalizadas e um script para importar dados de uma API externa.

Desde já, agradeço a oportunidade de estar participando do processo👨‍💻

Desenvolvedor:  Thiago Fernandes Xavier de Souza
Linkedin:       www.linkedin.com/in/thiago-fernandes-b7bb64252


## 📚 Endpoints da API

# Produtos
GET /products - Lista todos os produtos.
GET /products/{id} - Retorna um produto pelo ID.
POST /products - Cria um novo produto.
PUT /products/{id} - Atualiza um produto existente.
DELETE /products/{id} - Remove um produto.

# Consultas Personalizadas
GET /products/search/categoryandname - Busca produtos por nome e categoria.
GET /products/search/category - Busca produtos por categoria.
GET /products/search/image - Busca produtos com ou sem imagem.


## Importação de produtos de uma API externa

# Para obter todos os produtos de uma API externa:
go run cmd/importer/main.go

# Para obter um produto de uma API externa com ID específico
go run cmd/importer/main.go --id=12    -> exemplo de id

## Estrutura do projeto

/braip
├── /cmd
│   └── /importer
│       └── main.go                 # Script para importar dados externos
├── /internal
│   ├── /api
│   │   └── product_handler.go      # Handlers da API
│   ├── /database
│   │   └── db.go                   # Configuração do banco de dados
│   ├── /models
│   │   └── products.go             # Definição dos modelos
│   ├── /repository
│   │   └── product_repository.go   # Acesso ao banco de dados
│   └── /services
│       └── product_service.go      # Lógica de negócio
├── database.db
├── go.mod
├── go.sum
├── main.go                         # Ponto de entrada da aplicação
└── README.md                       # Documentação do projeto