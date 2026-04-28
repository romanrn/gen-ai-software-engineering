# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All Go commands must be run from `src/` (the module root):

```bash
# Run the server locally
cd src && go run ./cmd/api

# Build binary
cd src && go build -o server ./cmd/api

# Fetch / tidy dependencies
cd src && go mod tidy

# Format
cd src && go fmt ./...

# Run all tests
cd src && go test ./...

# Run a single package's tests
cd src && go test ./internal/service/...

# Run a single test by name
cd src && go test ./internal/service/... -run TestGetAccountBalance_ComputesFromTransactionTypes

# Run tests with coverage
cd src && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# Run with Docker (API only)
docker compose up --build

# Run with Docker + Swagger UI on :8081
docker compose --profile docs up --build

# Generate OpenAPI spec (installs swag CLI if missing, requires Go)
./scripts/generate-swagger.sh
```

Environment variables: `PORT` (default `8080`), `LOG_LEVEL` (`debug|info|warn|error`, default `info`).

## Architecture

Hexagonal (Ports & Adapters). Dependency rule: everything points inward — adapters depend on ports, ports depend only on domain.

```
domain  ←  ports/in|out  ←  service  ←  adapters/in|out
```

**`internal/domain/`** — pure entities (`Transaction`, `AccountBalance`, `AccountSummary`), typed errors (`ErrNotFound`, `ErrInvalidInput`, `ErrValidation`). Zero external dependencies. `ErrValidation` carries per-field `[]ValidationDetail` — this is what the error handler serializes into the `details` array.

**`internal/ports/in/`** — `TransactionService` interface (use cases). Also defines `CreateTransactionInput` and `TransactionFilter` (used by handlers and service).

**`internal/ports/out/`** — `TransactionRepository` interface (`Save`, `FindByID`, `FindAll`). Filtering is intentionally NOT part of the repository contract — it is applied in the service layer after `FindAll`.

**`internal/service/`** — implements `ports/in.TransactionService`. `validator.go` contains all validation rules (amount precision via string formatting, account regex `ACC-[A-Z0-9]{5}`, ISO 4217 currency allowlist, type enum). Returns `*domain.ErrValidation` on failure, which flows to the centralized error handler.

**`internal/adapters/in/http/`** — Fiber handlers (package name `handler`). Translate HTTP ↔ service calls only — no business logic. Return errors directly; Fiber routes them to `pkg/errorhandler.Handler`. `swagger_models.go` contains response structs used exclusively for Swagger annotation (`HealthResponse`, `ErrorResponse`, `ValidationErrorResponse`).

**`internal/adapters/out/memory/`** — `sync.RWMutex`-protected map. Implements `ports/out.TransactionRepository`.

**`cmd/api/`** — wiring only. `main.go` instantiates layers and injects dependencies. `server.go` configures Fiber (`ErrorHandler`, `recover` middleware, routes).

**`pkg/`** — cross-cutting concerns importable by any layer:
- `pkg/errorhandler` — single Fiber `ErrorHandler`: unwraps domain errors → HTTP status + JSON. Uses `errors.As`/`errors.Is` for matching.
- `pkg/middleware` — `RequestID()` (reads or generates `X-Request-ID`, stores in `fiber.Locals` under key `RequestIDKey`) and `Logger()` (structured `slog` per-request line).
- `pkg/logger` — configures `slog.SetDefault` with `JSONHandler` and level from `LOG_LEVEL` env.

## Testing approach

Tests use `github.com/stretchr/testify` (`assert` / `require`). Service tests instantiate a real `memory.TransactionRepository` — no mocks. This is intentional: the in-memory store is lightweight enough that mocking it would only hide bugs.

## Swagger

`docs/swagger.yaml` is generated from code annotations via `swag`. Run `./scripts/generate-swagger.sh` from the project root to regenerate after adding or changing annotations. The script auto-installs `swag` if not present. Swagger UI is **not** started by default — use `docker compose --profile docs up` to include it on `http://localhost:8081`. The `swagger-ui` service waits for the `api` health check to pass before starting.

## Key conventions

- **Error flow**: handlers `return err` → Fiber calls `errorhandler.Handler` → maps to HTTP status. Never write `c.Status(...).JSON(...)` for errors inside handlers.
- **Adding a new domain error**: define it in `domain/errors.go`, then add a matching `errors.Is`/`errors.As` branch in `pkg/errorhandler/error_handler.go`.
- **Adding a new endpoint**: add the use case method to `ports/in/transaction_service.go`, implement in `service/transaction_service.go`, add a handler method in `adapters/in/http/`, register the route in `cmd/api/server.go`.
- **Request ID** is available inside any handler via `c.Locals(middleware.RequestIDKey).(string)` and is included in every error response automatically by the error handler.
- HTTP handler package is named `handler` (not `http`) to avoid shadowing the stdlib package.
