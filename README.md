# curriculum-go

## 📖 Descrição
API RESTful em **Go**, migrada do projeto [`go-api`](https://github.com/joelsouza82/go-api) e reescrita seguindo a **Arquitetura Hexagonal (Ports & Adapters)**. A aplicação oferece um CRUD completo para o gerenciamento de **informações pessoais/perfil** e de **credenciais de login**.

---

## 🛠️ Linguagem e Tecnologias Utilizadas

- **Linguagem:** [Go (Golang)](https://go.dev/) (v1.26+)
- **Framework Web:** [Gin Gonic](https://github.com/gin-gonic/gin) — usado apenas no adapter de entrada HTTP.
- **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/), acessado via `database/sql` + driver `github.com/lib/pq` no adapter de saída.
- **Testes:** `testify` (asserts e mocks) + `go-sqlmock` para os testes do adapter de banco.
- **Containerização:** Docker e Docker Compose.

---

## 🏗️ Arquitetura Hexagonal (Ports & Adapters)

O núcleo da aplicação (domínio + regras de negócio) não depende de nenhum framework, banco de dados ou protocolo HTTP. Toda comunicação com o mundo externo passa por **portas** (interfaces), implementadas por **adapters**:

```
                         ┌───────────────────────────┐
                         │   Adapter de Entrada (in)  │
                         │   internal/adapter/in/http │
                         │   (Gin: handlers + router) │
                         └──────────────┬─────────────┘
                                        │ implementa/chama
                                        ▼
                         ┌───────────────────────────┐
                         │   Porta de Entrada (in)    │
                         │  internal/core/port/in     │
                         │  (PersonalService,         │
                         │   LoginService)            │
                         └──────────────┬─────────────┘
                                        │
                         ┌──────────────▼─────────────┐
                         │   Núcleo / Domínio          │
                         │  internal/core/domain       │
                         │  internal/core/service       │
                         │  (regras de negócio puras)  │
                         └──────────────┬─────────────┘
                                        │
                         ┌──────────────▼─────────────┐
                         │   Porta de Saída (out)      │
                         │  internal/core/port/out     │
                         │  (PersonalRepository,       │
                         │   LoginRepository)          │
                         └──────────────┬─────────────┘
                                        │ implementa
                                        ▼
                         ┌───────────────────────────┐
                         │  Adapter de Saída (out)    │
                         │ internal/adapter/out/postgres│
                         │ (database/sql + lib/pq)    │
                         └───────────────────────────┘
```

### Camadas

1. **`internal/core/domain`** — Entidades (`Personal`, `Login`) e erros de domínio (`ErrPersonalNotFound`, `ErrLoginNotFound`). Não conhece HTTP, SQL ou nenhum framework.
2. **`internal/core/port/in`** — Portas de entrada (primárias): interfaces que o núcleo expõe para o mundo externo (`PersonalService`, `LoginService`).
3. **`internal/core/port/out`** — Portas de saída (secundárias): interfaces que o núcleo espera que o mundo externo implemente (`PersonalRepository`, `LoginRepository`).
4. **`internal/core/service`** — Implementação das portas de entrada, contendo as regras de negócio. Depende apenas das portas de saída (nunca de um adapter concreto).
5. **`internal/adapter/in/http`** — Adapter de entrada: handlers Gin + roteador, traduzem requisições HTTP em chamadas às portas de entrada.
6. **`internal/adapter/out/postgres`** — Adapter de saída: implementação concreta das portas de saída usando PostgreSQL.
7. **`internal/config`** — Carrega configuração (porta do servidor, credenciais do banco) a partir de variáveis de ambiente.
8. **`internal/mocks`** — Mocks (testify) das portas de entrada e saída, usados nos testes unitários.
9. **`cmd/api/main.go`** — *Composition root*: é o único ponto do sistema que conhece e conecta todas as implementações concretas (config → conexão → repository → service → handler → router).

A regra de dependência é sempre **de fora para dentro**: adapters dependem do núcleo através das portas; o núcleo nunca depende de um adapter.

---

## 🗄️ Banco de Dados

### Estrutura das Tabelas

#### Tabela `personal`
- `id` (INTEGER / SERIAL - Chave primária)
- `name`, `rg`, `document` (VARCHAR)
- `address`, `city`, `neighborhood`, `state`, `cep`, `phone`, `email`, `website`, `linkedin`, `github` (VARCHAR)
- `birthdate` (TIMESTAMP)
- `login_id` (INTEGER - Chave estrangeira para `login.id`)

#### Tabela `login`
- `id` (INTEGER / SERIAL - Chave primária)
- `email`, `password` (VARCHAR)

---

## 📂 Estrutura do Projeto

```text
curriculum-go/
├── cmd/
│   └── api/
│       └── main.go                        # Composition root
├── internal/
│   ├── config/
│   │   └── config.go                      # Configuração via variáveis de ambiente
│   ├── core/
│   │   ├── domain/
│   │   │   ├── personal.go
│   │   │   ├── login.go
│   │   │   └── errors.go
│   │   ├── port/
│   │   │   ├── in/                        # Portas de entrada (services)
│   │   │   └── out/                       # Portas de saída (repositories)
│   │   └── service/                       # Casos de uso / regras de negócio
│   ├── adapter/
│   │   ├── in/
│   │   │   └── http/                      # Handlers Gin + router
│   │   └── out/
│   │       └── postgres/                  # Implementação PostgreSQL
│   └── mocks/                             # Mocks das portas para testes
├── Dockerfile
├── docker-compose.yml
├── .env.example                           # Modelo das variáveis de ambiente
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Configuração

Nenhuma credencial fica commitada no repositório. A aplicação lê a configuração de variáveis de ambiente:

| Variável       | Obrigatória | Padrão    | Descrição                          |
|----------------|:-----------:|-----------|-------------------------------------|
| `SERVER_PORT`  | não         | `8000`    | Porta HTTP do servidor              |
| `DB_HOST`      | sim         | —         | Host do PostgreSQL                  |
| `DB_PORT`      | sim         | —         | Porta do PostgreSQL                 |
| `DB_USER`      | sim         | —         | Usuário do banco                    |
| `DB_PASSWORD`  | sim         | —         | Senha do banco                      |
| `DB_NAME`      | sim         | —         | Nome do banco                       |
| `DB_SSLMODE`   | não         | `require` | Modo SSL da conexão                 |

Copie `.env.example` para `.env` e preencha com suas credenciais antes de rodar a aplicação.

---

## 🚀 Como Executar o Projeto

### Opção 1: Executar Localmente com Go

1. **Baixar as dependências:**
   ```bash
   go mod download
   ```

2. **Configurar as variáveis de ambiente:**
   ```bash
   cp .env.example .env
   # edite o .env com suas credenciais
   export $(grep -v '^#' .env | xargs)
   ```

3. **Iniciar a aplicação:**
   ```bash
   go run cmd/api/main.go
   ```
   *A API estará acessível na porta 8000: `http://localhost:8000`*

---

### Opção 2: Executar via Docker Compose

1. **Configurar as variáveis de ambiente:**
   ```bash
   cp .env.example .env
   # edite o .env com suas credenciais
   ```

2. **Build e inicialização do container:**
   ```bash
   docker-compose up --build
   ```
   *A API estará acessível em `http://localhost:8000`*

---

### 🧪 Executando os Testes

```bash
go test ./...
```

---

## 📡 Endpoints da API

### 🟢 Health Check

#### `GET /ping`
- **Resposta Sucesso (`200 OK`):**
  ```json
  { "message": "pong" }
  ```

---

### 👤 Informações Pessoais (`/personal` e `/personals`)

#### `GET /personals`
Retorna todas as informações de perfis pessoais armazenadas.

#### `GET /personal/:personalId`
Retorna um perfil pessoal específico pelo ID.
- **Erro (`400`):** ID inválido.
- **Erro (`404`):** registro não encontrado.

#### `POST /personal`
Cria um novo perfil pessoal. Resposta `201 Created`.

#### `PUT /personal/:personalId`
Atualiza um perfil pessoal existente. Resposta `200 OK` ou `404 Not Found`.

#### `DELETE /personal/:personalId`
Remove um perfil pessoal pelo ID. Resposta `204 No Content` ou `404 Not Found`.

---

### 🔐 Login (`/login` e `/logins`)

#### `GET /logins`
Retorna todos os registros de login armazenados.

#### `GET /login/:loginId`
Retorna um registro de login específico pelo ID.
- **Erro (`400`):** ID inválido.
- **Erro (`404`):** registro não encontrado.

#### `POST /login`
Cria um novo registro de login. Resposta `201 Created`.

#### `PUT /login/:loginId`
Atualiza um registro de login existente. Resposta `200 OK` ou `404 Not Found`.

#### `DELETE /login/:loginId`
Remove um registro de login pelo ID. Resposta `204 No Content` ou `404 Not Found`.

---

## 🛠️ Contribuição
Sinta-se à vontade para enviar *pull requests* ou abrir *issues* com sugestões e melhorias.
