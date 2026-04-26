# Banking Transactions API — Implementation Plan

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.23+ |
| Framework | Fiber v2 (`github.com/gofiber/fiber/v2`) |
| Logging | slog (`log/slog`, stdlib since Go 1.21) |
| Storage | In-memory (`sync.RWMutex`-protected map) |
| Runtime | Docker (multi-stage build) |

---

## Architecture

Hexagonal (Ports & Adapters) with `in/out` naming convention.

**Dependency rule**: arrows always point inward — adapters depend on ports, ports never depend on adapters. Domain has zero external dependencies.

| Layer | Package | Role |
|-------|---------|------|
| Domain | `internal/domain` | Pure entities, value objects, typed errors |
| Ports in | `internal/ports/in` | Use case interfaces (driven by HTTP adapter) |
| Ports out | `internal/ports/out` | Repository interfaces (implemented by storage adapter) |
| Service | `internal/service` | Business logic + validation, implements input ports |
| Adapter in | `internal/adapters/in/http` | HTTP handlers — translate HTTP ↔ service calls |
| Adapter out | `internal/adapters/out/memory` | In-memory repository, implements output ports |
| Bootstrap | `cmd/api` | Fiber app config, route registration, DI wiring |
| Cross-cutting | `pkg/middleware`, `pkg/errorhandler`, `pkg/logger` | Request ID, structured logging, error mapping |

---

## Project Structure

```
homework-1/
├── cmd/
│   └── api/
│       ├── main.go                      # DI wiring, app startup
│       └── server.go                    # Fiber app config, middleware chain, route registration
├── internal/
│   ├── domain/
│   │   ├── transaction.go               # Transaction entity, value objects
│   │   ├── account.go                   # Account balance logic
│   │   └── errors.go                    # Typed domain errors (ErrNotFound, ErrValidationFailed, ...)
│   ├── ports/
│   │   ├── in/
│   │   │   └── transaction_service.go   # Use case interfaces (context.Context as first arg)
│   │   └── out/
│   │       └── transaction_repository.go
│   ├── service/
│   │   ├── transaction_service.go       # Implements input ports, orchestrates domain
│   │   └── validator.go                 # Amount, account format, currency validation rules
│   └── adapters/
│       ├── in/
│       │   └── http/
│       │       ├── transaction_handler.go
│       │       └── account_handler.go
│       └── out/
│           └── memory/
│               └── transaction_repository.go
├── pkg/
│   ├── middleware/
│   │   ├── request_id.go                # X-Request-ID generation and propagation
│   │   └── logger.go                    # Structured per-request logging
│   ├── errorhandler/
│   │   └── error_handler.go             # Centralized Fiber error handler (domain errors → HTTP)
│   └── logger/
│       └── logger.go                    # slog setup (JSON handler, log level from env)
├── Dockerfile                           # Multi-stage: golang:1.23-alpine → distroless/static
├── docker-compose.yml                   # Port 8080:8080, LOG_LEVEL env, health check
├── .dockerignore
├── .gitignore
├── go.mod
├── go.sum
├── demo/
│   ├── run.sh                           # docker compose up wrapper
│   ├── sample-requests.http             # Sample API calls (VS Code REST Client)
│   └── sample-data.json
├── docs/
│   └── screenshots/
├── README.md
└── HOWTORUN.md
```

---

## Endpoints

### Transactions (`transaction_handler.go`)

| Method | Path | Description | Status Codes |
|--------|------|-------------|--------------|
| `POST` | `/transactions` | Create a new transaction | `201`, `400`, `422` |
| `GET` | `/transactions` | List all transactions (filterable) | `200` |
| `GET` | `/transactions/:id` | Get a specific transaction by ID | `200`, `404` |

#### POST `/transactions`

Request body:
```json
{
  "fromAccount": "ACC-12345",
  "toAccount":   "ACC-67890",
  "amount":      100.50,
  "currency":    "USD",
  "type":        "transfer"
}
```

Response `201`:
```json
{
  "id":          "uuid",
  "fromAccount": "ACC-12345",
  "toAccount":   "ACC-67890",
  "amount":      100.50,
  "currency":    "USD",
  "type":        "transfer",
  "status":      "completed",
  "timestamp":   "2024-01-15T10:30:00Z"
}
```

#### GET `/transactions` — query params

| Param | Example | Description |
|-------|---------|-------------|
| `accountId` | `ACC-12345` | Match fromAccount or toAccount |
| `type` | `transfer` | Filter by transaction type |
| `from` | `2024-01-01` | Start of date range (inclusive) |
| `to` | `2024-01-31` | End of date range (inclusive) |

Filters are combinable (AND logic).

---

### Accounts (`account_handler.go`)

| Method | Path | Description | Status Codes |
|--------|------|-------------|--------------|
| `GET` | `/accounts/:accountId/balance` | Get current balance | `200`, `404` |
| `GET` | `/accounts/:accountId/summary` | Get transaction summary | `200`, `404` |

#### GET `/accounts/:accountId/balance`

Response `200`:
```json
{
  "accountId": "ACC-12345",
  "balance":   250.00,
  "currency":  "USD"
}
```

#### GET `/accounts/:accountId/summary`

Response `200`:
```json
{
  "accountId":           "ACC-12345",
  "totalDeposits":       1000.00,
  "totalWithdrawals":    750.00,
  "transactionCount":    12,
  "lastTransactionDate": "2024-01-15T10:30:00Z"
}
```

---

### Infrastructure

| Method | Path | Description | Status Codes |
|--------|------|-------------|--------------|
| `GET` | `/health` | Liveness check (used by Docker health check) | `200` |

---

## Validation

### Rules

| Field | Rule |
|-------|------|
| `amount` | Positive number, max 2 decimal places |
| `fromAccount` / `toAccount` | Format `ACC-XXXXX` (X = alphanumeric, case-insensitive) |
| `currency` | Valid ISO 4217 code (USD, EUR, GBP, JPY, CHF, ...) |
| `type` | One of: `deposit`, `withdrawal`, `transfer` |

### Error response `422`

```json
{
  "error":      "Validation failed",
  "request_id": "uuid",
  "details": [
    {"field": "amount",      "message": "must be a positive number with max 2 decimal places"},
    {"field": "currency",    "message": "invalid ISO 4217 code"},
    {"field": "fromAccount", "message": "must match format ACC-XXXXX"}
  ]
}
```

---

## Cross-cutting Concerns

### Traceability
- `X-Request-ID` middleware: reads from incoming header or generates a UUID
- Request ID flows: `HTTP header → fiber.Locals → context.Context → service layer → repository`
- All service methods accept `context.Context` as first parameter

### Logging
- Structured JSON logs via `log/slog` with `slog.NewJSONHandler` (fits Docker log aggregators: Loki, CloudWatch, etc.)
- Log level configurable via `LOG_LEVEL` env var
- Every log line carries: `request_id`, `method`, `path`, `status`, `latency_ms`

### Error Handling

**Domain errors** (`domain/errors.go`) — typed sentinel errors:

| Error | HTTP Status |
|-------|-------------|
| `ErrNotFound` | `404 Not Found` |
| `ErrValidationFailed` | `422 Unprocessable Entity` |
| `ErrInvalidInput` | `400 Bad Request` |
| `ErrInternal` | `500 Internal Server Error` |

**Centralized error handler** (`pkg/errorhandler/error_handler.go`):
- Registered as `fiber.Config.ErrorHandler`
- Uses `errors.As` to unwrap and map domain errors → HTTP status + structured JSON response
- Catches unhandled `*fiber.Error` (e.g. `404` from unknown routes)
- Falls back to `500` for unexpected errors, logs the full error with `request_id`

**Panic recovery**:
- Fiber built-in `recover` middleware wraps all handlers
- Logs panic stacktrace with `request_id`, returns `500` without leaking internals

**Error response shape** (all errors):
```json
{
  "error":      "human-readable message",
  "request_id": "uuid",
  "details":    []
}
```

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/gofiber/fiber/v2` | HTTP framework |
| `log/slog` | Structured JSON logging (stdlib, no external dep) |
| `github.com/google/uuid` | Request ID + transaction ID generation |

---

## Implementation Order

| Step | What | Files |
|------|------|-------|
| 1 | Init module, scaffold dirs | `go.mod`, all dirs |
| 2 | Domain entities + errors | `domain/` |
| 3 | Port interfaces | `ports/in/`, `ports/out/` |
| 4 | Validator | `service/validator.go` |
| 5 | In-memory repository | `adapters/out/memory/` |
| 6 | Service (business logic) | `service/transaction_service.go` |
| 7 | HTTP handlers | `adapters/in/http/` |
| 8 | Middleware + logger | `pkg/` |
| 9 | Server bootstrap + wiring | `cmd/api/` |
| 10 | Dockerfile + Compose | `Dockerfile`, `docker-compose.yml` |
| 11 | Demo files + docs | `demo/`, `README.md`, `HOWTORUN.md` |
