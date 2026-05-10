# Intelligent Customer Support System - Developer Reference

## Documentation

| Document | Audience | Content |
|----------|----------|---------|
| [docs/README.md](docs/README.md) | Developers | Project overview, quick start, project structure |
| [docs/API_REFERENCE.md](docs/API_REFERENCE.md) | API consumers | All endpoints, request/response schemas, cURL examples |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Technical leads | Component design, data flow diagrams, design decisions |
| [docs/TESTING_GUIDE.md](docs/TESTING_GUIDE.md) | QA engineers | Test inventory, coverage report, manual checklist, benchmarks |

---

## Architecture

Hexagonal architecture — domain at the center, everything else plugged in through interfaces.

```
Domain Layer  →  Ports (interfaces)  →  Adapters (implementations)
```

```
src/
├── cmd/api/                   # Entry point — DI wiring, Fiber server setup
├── internal/
│   ├── domain/                # Ticket entity, enums, typed errors (zero external imports)
│   ├── ports/
│   │   ├── in/                # TicketService interface (use cases)
│   │   └── out/               # TicketRepository interface (storage)
│   ├── service/               # Validator · Classifier · TicketServiceImpl
│   └── adapters/
│       ├── in/http/           # Fiber HTTP handlers (no business logic)
│       └── out/memory/        # sync.RWMutex map repository
├── pkg/
│   ├── importer/              # CSV / JSON / XML parsers
│   ├── errorhandler/          # Global Fiber error mapper
│   └── middleware/            # RequestID, logger
└── tests/                     # Black-box integration + performance tests
    └── fixtures/              # Sample data for tests
```

---

## Development Commands

```bash
# Run locally
go run ./cmd/api

# Build binary
go build -o server ./cmd/api

# Run all tests
go test ./...

# Cross-package coverage (recommended)
go test -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
go tool cover -html=coverage.out -o docs/screenshots/coverage.html

# Integration tests only
go test ./tests/... -v

# Race detector
go test ./... -race

# Benchmarks
go test ./tests/... -bench=. -benchmem

# Generate Swagger docs
./scripts/generate-swagger.sh

# Docker
docker compose up --build              # API on :8080
docker compose --profile docs up       # API + Swagger UI on :8080/:8081
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP listen port |
| `LOG_LEVEL` | `info` | slog level: `debug`, `info`, `warn`, `error` |

---

## Test Coverage

Overall: **89.1%** (cross-package). See [docs/TESTING_GUIDE.md](docs/TESTING_GUIDE.md) for the full breakdown and test inventory.

| Package | Coverage |
|---------|----------|
| `internal/domain` | 100.0% |
| `pkg/importer` | 95.5% |
| `internal/service` | 89.1% |
| `internal/adapters/in/http` | 83.1% |
| `internal/adapters/out/memory` | 78.1% |
| `cmd/api` | 63.0% |

---

## Key Design Decisions

1. **No mocks in tests** — real implementations wired via `newApp()`; catches integration bugs
2. **Partial success on import** — per-row errors collected; bad rows don't abort the batch
3. **Hexagonal ports** — storage and HTTP are swappable without touching business logic
4. **Keyword classifier** — deterministic, zero-latency, fully testable offline

Full rationale and trade-offs: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
