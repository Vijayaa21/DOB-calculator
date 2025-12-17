# DOB Calculator API

A RESTful API built with Go to manage users with their name and date of birth. The API calculates and returns a user's age dynamically.

## 🔧 Tech Stack

- [GoFiber](https://gofiber.io/) - Fast HTTP framework
- PostgreSQL with [SQLC](https://sqlc.dev/) - Type-safe SQL
- [Uber Zap](https://github.com/uber-go/zap) - Structured logging
- [go-playground/validator](https://github.com/go-playground/validator) - Input validation
- Docker & Docker Compose

## 📁 Project Structure

```
.
├── cmd/server/main.go          # Application entry point
├── config/                     # Configuration management
├── db/
│   ├── migrations/             # SQL migrations
│   ├── queries/                # SQLC queries
│   └── sqlc/                   # Generated SQLC code
├── internal/
│   ├── handler/                # HTTP handlers
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic
│   ├── routes/                 # Route definitions
│   ├── middleware/             # HTTP middleware
│   ├── models/                 # Data models
│   └── logger/                 # Zap logger setup
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Running with Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Running Locally

1. **Setup PostgreSQL**
   
   Create a database named `dob_calculator` and run the migration:
   ```bash
   psql -U postgres -c "CREATE DATABASE dob_calculator;"
   psql -U postgres -d dob_calculator -f db/migrations/001_create_users_table.sql
   ```

2. **Configure Environment**
   
   Copy the example env file and update values:
   ```bash
   cp .env.example .env
   ```

3. **Install Dependencies**
   ```bash
   go mod download
  # optional, if you plan to modify SQL in db/queries
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

4. **Run the Application**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:3000`

> Note (Windows / local Postgres): If your local PostgreSQL runs on a different port (e.g. 5433), update `.env` accordingly. This repo’s `.env.example` uses 5433 by default.

## 📊 API Endpoints

### Health Check
```
GET /health
```

### Create User
```
POST /users
Content-Type: application/json

{
  "name": "Alice",
  "dob": "1990-05-10"
}

Response: 201 Created
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User by ID
```
GET /users/:id

Response: 200 OK
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

### List All Users
```
GET /users?page=1&page_size=10

Response: 200 OK
[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10",
    "age": 35
  }
]

# For paginated response with metadata
GET /users?page=1&page_size=10&paginated=true

Response: 200 OK
{
  "users": [...],
  "total": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

### Update User
```
PUT /users/:id
Content-Type: application/json

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}

Response: 200 OK
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Delete User
```
DELETE /users/:id

Response: 204 No Content
```

## 🧪 Running Tests

```bash
go test ./... -v
```

## 📝 Regenerating SQLC Code

If you modify the SQL queries, regenerate the SQLC code:

```bash
sqlc generate
```

If you don’t have `sqlc` yet, install it first:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## 🔒 Middleware Features

- **Request ID**: Injects unique `X-Request-ID` header in responses
- **Request Logger**: Logs request method, path, status, and duration
- **CORS**: Handles Cross-Origin Resource Sharing
- **Error Handler**: Consistent JSON error responses

## 📄 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_PORT | 3000 | HTTP server port |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5433 | PostgreSQL port |
| DB_USER | postgres | PostgreSQL user |
| DB_PASSWORD | postgres | PostgreSQL password |
| DB_NAME | dob_calculator | Database name |
| DB_SSLMODE | disable | SSL mode |

## 📜 License

MIT

---

## ✅ Deliverables Checklist

- Store dob in DB: users.dob (DATE) is persisted via SQLC queries.
- Age calculated dynamically: handled in `internal/models/user.go` and returned by GET endpoints.
- SQLC used for DB access: queries in `db/queries/users.sql`, generated code in `db/sqlc/*`, config `sqlc.yaml`.
- Input validation: `go-playground/validator` in `internal/handler/user_handler.go`.
- Logging: Uber Zap set up in `internal/logger`, request logs via middleware.
- Clean HTTP status codes and error handling: centralized error handler + appropriate codes in handlers.

### Bonus (Included)
- Docker support: `Dockerfile` + `docker-compose.yml` (Postgres + app).
- Pagination: `/users?page=1&page_size=10` and `paginated=true` for metadata.
- Unit tests: age/date helpers covered in `internal/models/user_test.go`.
- Middleware: `X-Request-ID` injection and request duration logging in `internal/middleware`.

---

## 🔗 Submission

1. Initialize a Git repository and push to GitHub:

```bash
git init
git add .
git commit -m "Initial submission: DOB Calculator API"
git branch -M main
git remote add origin https://github.com/<your-gh-username>/<your-repo-name>.git
git push -u origin main
```

2. Share the repository link.

---

## 🛠️ Troubleshooting

- Port already in use: stop existing process on 3000 or set `SERVER_PORT`.
  ```powershell
  Get-NetTCPConnection -LocalPort 3000 -State Listen | Select-Object -ExpandProperty OwningProcess | % { Stop-Process -Id $_ -Force }
  ```
- Local Postgres on a non-default port (e.g. 5433): update `DB_PORT` in `.env`.
- Creating DB/table manually (Windows example):
  ```powershell
  $env:PGPASSWORD='your_password'
  & "C:\\Program Files\\PostgreSQL\\18\\bin\\psql.exe" -U postgres -p 5433 -c "CREATE DATABASE dob_calculator;"
  & "C:\\Program Files\\PostgreSQL\\18\\bin\\psql.exe" -U postgres -p 5433 -d dob_calculator -c "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, dob DATE NOT NULL);"
  ```
