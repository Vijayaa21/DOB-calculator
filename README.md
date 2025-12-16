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
   ```

4. **Run the Application**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:3000`

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
