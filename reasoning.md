# Reasoning and Design Notes

This document summarizes my approach, design decisions, and how the service works end‑to‑end. It is intended to make the code review and discussion efficient.

## Goals

- Store a user's `name` and `dob` in a relational database.
- Expose REST endpoints to create, read, update, delete, and list users.
- Return `age` dynamically when fetching users (not stored in DB).
- Keep the implementation simple, observable, and testable.

## Architecture

- Layered design with clear boundaries:
  - Handler (HTTP): request/response, validation, HTTP statuses
  - Service (Business): input parsing, orchestration, age calculation
  - Repository (Data): SQLC‑generated typed queries
  - DB (Infrastructure): pgx pool wiring, connection management
- Framework: Fiber for fast HTTP with a light ergonomic API.
- Logging: Uber Zap for structured logs; logs include requestId, method, path, status, duration.

### Data Model

- Table: `users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, dob DATE NOT NULL)`
- Domain model returned to clients augments DB model with calculated `age`.

## Age Calculation Strategy

- Implemented in `internal/models/user.go`.
- `CalculateAge(time.Time) int` uses current local time and subtracts birth year, decrementing if the birthday has not occurred this year.
- Covered edge cases in tests: birthday today/yesterday/tomorrow, leap year DOBs (Feb 29), newborn.
- Rationale: Simple, correct for calendar arithmetic at day granularity; no time zone reliance in persistence (DATE), presentation via `YYYY-MM-DD`.

## Validation & Error Handling

- Validation via `go-playground/validator` in handlers:
  - `name`: required, length bounds
  - `dob`: required, `datetime=2006-01-02`
- Centralized error handler returns consistent JSON shape and attaches `request_id`.
- Status codes:
  - 201 Create
  - 200 OK, 204 No Content
  - 400 Validation/parse errors
  - 404 When ID not found (service layer sentinel)
  - 405 Method not allowed (framework)
  - 500 Internal errors (DB or unexpected)

## SQLC + Database Access

- SQL: `db/queries/users.sql`; SQLC config: `sqlc.yaml`.
- Generated code is checked in under `db/sqlc/` to avoid requiring SQLC during review.
- Driver: pgx/v5 with a connection pool (`pgxpool.Pool`). The SQLC `DBTX` is satisfied directly by `pgxpool.Pool`.
- Queries use placeholders and typed params; outputs map to generated `User` model.

## Middleware

- `RequestID`: adds `X-Request-ID` to responses and stores the value in context for logs.
- `RequestLogger`: logs method, path, status, duration, IP, user agent.
- `CORS`: permissive default for demo/dev.
- `ErrorHandler`: consistent JSON errors and structured logs.

## Routes & Handlers

- `/users` POST: create user; returns id/name/dob.
- `/users/:id` GET: return id/name/dob/age.
- `/users` GET: list users with ages. Supports `page`, `page_size`, and optional `paginated=true` to include metadata.
- `/users/:id` PUT: update name/dob.
- `/users/:id` DELETE: 204 on success.
- Compatibility: routes are available both under root and `/api/v1`.

## Configuration

- `.env` (via `godotenv`) with sensible defaults.
- Example file `.env.example` uses `DB_PORT=5433` because Windows local PostgreSQL often runs on 5433; adjust as needed.
- `config.DatabaseConfig.GetDSN()` shows how DSNs are built; connection is created in `db/database.go` and validated with a ping.

## Docker

- Multi‑stage Dockerfile builds a static binary for a minimal alpine runtime image.
- `docker-compose.yml` spins Postgres and the app; applies init migration via bind‑mounted SQL.
- Rationale: provides a consistent, reproducible environment.

## Testing

- `internal/models/user_test.go` exercises:
  - `CalculateAge` across edge cases
  - `ParseDate` and `FormatDate`
- Focused where deterministic, fast, and highest ROI for correctness.

## Observability

- Structured Zap logs at info/error with context.
- Request correlation via `X-Request-ID`.
- Central error handler emits errors with stack traces.

## Security Considerations

- Parameterized SQL via SQLC eliminates injection risks.
- Input validation rejects malformed payloads early.
- Secrets via environment variables; no secrets committed.
- CORS set to `*` for demo; constrain for production.
- TLS termination not included; assume reverse proxy/load balancer.

## Performance Considerations

- Fiber provides efficient routing; pgx uses a pooled connection model.
- Queries are simple and indexed by primary key.
- Pagination with `LIMIT/OFFSET`; sufficient for small/medium datasets.

## Trade‑offs & Alternatives

- Stored procedures vs plain SQL: chose plain SQL for clarity and type‑safety with SQLC.
- ORM vs SQLC: SQLC yields compile‑time safety and minimal overhead.
- Time zones: persisted as DATE, formatted and parsed as `YYYY-MM-DD` strings to avoid TZ ambiguity.
- Migrations: simple SQL migration file and docker init; a tool like Goose/Migrate could be added for versioning.

## Future Improvements

- OpenAPI/Swagger for API documentation.
- Health endpoint that also checks DB readiness.
- More tests (handlers/services) with DB test container.
- CI workflow (fmt, vet, tests, build, docker) and linting.
- Rate limiting and more robust CORS configuration for production.
- Unique constraints (e.g., name + dob) and additional indexes depending on access patterns.

## End‑to‑End Flow (Create → Read)

1. Client POSTs JSON to `/users`.
2. Handler parses & validates; service parses DOB; repository calls SQLC `CreateUser`.
3. On success, return 201 with id/name/dob.
4. Client GETs `/users/:id`; service fetches user, computes `age`, and returns id/name/dob/age.

---

This design keeps concerns separated, promotes testability, and uses widely adopted Go tooling (Fiber, pgx, SQLC, Zap, validator) to meet the requirements with clear, maintainable code.
