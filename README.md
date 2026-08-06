# Animal Rescue App

API REST simples para gerenciamento de animais em resgate/adoção, escrita em Go com `net/http` e SQLite.

## Funcionalidades

- listar animais
- buscar um animal por ID
- criar um animal
- atualizar um animal
- remover um animal com soft delete
- testes automatizados para os handlers

## Stack

- Go
- `net/http`
- SQLite com `modernc.org/sqlite`
- `go-sqlmock` para testes

## Estrutura do projeto

```text
animalrescueapp/
├── cmd/
│   └── api/
│       └── main.go
├── db/
│   └── migrations/
├── internal/
│   ├── database/
│   ├── handlers/
│   └── models/
├── go.mod
├── go.sum
└── README.md
```

### Pastas principais

- `cmd/api`: ponto de entrada da aplicação
- `internal/database`: acesso ao banco e queries
- `internal/handlers`: handlers HTTP da API
- `internal/models`: modelos e validações de domínio
- `db/migrations`: migrations SQL do banco

## Pré-requisitos

- Go instalado
- versão compatível com o `go.mod`

> Observação: o projeto usa `modernc.org/sqlite`, então não é necessário instalar SQLite nativo para compilar a aplicação.

## Como rodar

Na raiz do projeto:

```bash
go run ./cmd/api
```

A aplicação sobe em:

```text
http://localhost:8080
```

## Banco de dados

A aplicação abre um banco SQLite local chamado:

```text
app.db
```

O projeto possui migrations em `db/migrations`, no formato compatível com Goose.

### Schema atual

A tabela principal é `animals`, com os campos:

- `id`
- `name`
- `species`
- `age`
- `breed`
- `status`
- `created_at`
- `updated_at`
- `deleted_at`

### Soft delete

A exclusão é lógica. Quando um animal é removido, o campo `deleted_at` é preenchido e o registro deixa de aparecer nas consultas normais.

## Rodando os testes

```bash
go test ./...
```

## Endpoints

### `GET /animals`
Lista todos os animais não removidos.

**Resposta:** `200 OK`

---

### `GET /animals/{id}`
Busca um animal por ID.

**Resposta:**
- `200 OK`
- `404 Not Found`

---

### `POST /animals`
Cria um novo animal.

**Body de exemplo:**

```json
{
  "name": "Bob",
  "species": "Dog",
  "age": 2,
  "breed": "Labrador",
  "status": 0
}
```

**Resposta:**
- `201 Created`
- `400 Bad Request`
- `500 Internal Server Error`

---

### `PUT /animals/{id}`
Atualiza um animal existente.

**Body de exemplo:**

```json
{
  "name": "Bob",
  "species": "Dog",
  "age": 3,
  "breed": "Labrador",
  "status": 1
}
```

**Resposta:**
- `200 OK`
- `400 Bad Request`
- `500 Internal Server Error`

---

### `DELETE /animals/{id}`
Remove um animal com soft delete.

**Resposta:**
- `204 No Content`
- `400 Bad Request`
- `500 Internal Server Error`

## Status do animal

Atualmente o campo `status` é representado por inteiro:

- `0` = `available`
- `1` = `adopted`

## Exemplo com `curl`

### Criar um animal

```bash
curl -X POST http://localhost:8080/animals \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob",
    "species": "Dog",
    "age": 2,
    "breed": "Labrador",
    "status": 0
  }'
```

### Listar animais

```bash
curl http://localhost:8080/animals
```

### Buscar por ID

```bash
curl http://localhost:8080/animals/1
```

### Atualizar um animal

```bash
curl -X PUT http://localhost:8080/animals/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob",
    "species": "Dog",
    "age": 3,
    "breed": "Labrador",
    "status": 1
  }'
```

### Deletar um animal

```bash
curl -X DELETE http://localhost:8080/animals/1
```

## Autor

Projeto de estudo/prática para evolução em Go, APIs REST e organização de camadas (`handlers`, `database`, `models`).
