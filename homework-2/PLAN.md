# Homework-2: Intelligent Customer Support System — Implementation Plan

## Context

Build a customer support ticket management system in Go following the same hexagonal architecture established in homework-1. The system supports multi-format ticket imports (CSV/JSON/XML), auto-classification by keyword analysis, and a full CRUD REST API. Tech stack mirrors homework-1: Fiber, in-memory storage, testify, Docker/Compose, Swagger.

---

## Project Structure

```
homework-2/
├── src/
│   ├── cmd/api/
│   │   ├── main.go               # Entry point, DI wiring, Swagger annotations
│   │   └── server.go             # Fiber setup, middleware, routes
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── ticket.go         # Ticket entity, all enums/types, constructor
│   │   │   └── errors.go         # Typed domain errors (ErrNotFound, ErrValidation)
│   │   ├── ports/
│   │   │   ├── in/ticket_service.go      # Use case interface
│   │   │   └── out/ticket_repository.go  # Storage contract
│   │   ├── service/
│   │   │   ├── ticket_service.go  # CRUD + import + classify orchestration
│   │   │   ├── classifier.go      # Keyword-based auto-classification
│   │   │   └── validator.go       # All field validation rules
│   │   └── adapters/
│   │       ├── in/http/
│   │       │   ├── handler.go         # HTTP handlers (no business logic)
│   │       │   └── swagger_models.go  # Response types for Swagger annotations
│   │       └── out/memory/
│   │           └── ticket_repository.go  # sync.RWMutex map storage
│   └── pkg/
│       ├── errorhandler/   # Fiber global error handler
│       ├── logger/         # slog JSON setup
│       ├── middleware/     # RequestID + Logger middleware
│       └── importer/
│           ├── importer.go  # TicketImporter interface + ImportResult type
│           ├── csv.go       # encoding/csv parser
│           ├── json.go      # encoding/json parser
│           └── xml.go       # encoding/xml parser
├── docs/
│   ├── swagger.yaml
│   └── screenshots/
├── sample_data/
│   ├── sample_tickets.csv    # 50 tickets
│   ├── sample_tickets.json   # 20 tickets
│   ├── sample_tickets.xml    # 30 tickets
│   ├── invalid_tickets.csv   # Negative test data
│   ├── invalid_tickets.json
│   └── invalid_tickets.xml
├── scripts/
│   └── generate-swagger.sh
├── Dockerfile                # Multi-stage build (golang:1.23-alpine → alpine:3.20)
├── docker-compose.yml        # API on :8080 + Swagger UI on :8081 via profile
└── CLAUDE.md                 # Architecture guide + dev commands
```

---

## Domain Model (`internal/domain/ticket.go`)

```go
type Category   string  // account_access | technical_issue | billing_question | feature_request | bug_report | other
type Priority   string  // urgent | high | medium | low
type Status     string  // new | in_progress | waiting_customer | resolved | closed
type Source     string  // web_form | email | api | chat | phone
type DeviceType string  // desktop | mobile | tablet

type TicketMetadata struct {
    Source     Source     `json:"source"`
    Browser    string     `json:"browser"`
    DeviceType DeviceType `json:"device_type"`
}

type Classification struct {
    Category   Category `json:"category"`
    Priority   Priority `json:"priority"`
    Confidence float64  `json:"confidence"`      // 0.0–1.0
    Reasoning  string   `json:"reasoning"`
    Keywords   []string `json:"keywords_found"`
}

type Ticket struct {
    ID             uuid.UUID       `json:"id"`
    CustomerID     string          `json:"customer_id"`
    CustomerEmail  string          `json:"customer_email"`
    CustomerName   string          `json:"customer_name"`
    Subject        string          `json:"subject"`        // 1–200 chars
    Description    string          `json:"description"`    // 10–2000 chars
    Category       Category        `json:"category"`
    Priority       Priority        `json:"priority"`
    Status         Status          `json:"status"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
    ResolvedAt     *time.Time      `json:"resolved_at"`
    AssignedTo     *string         `json:"assigned_to"`
    Tags           []string        `json:"tags"`
    Metadata       TicketMetadata  `json:"metadata"`
    Classification *Classification `json:"classification,omitempty"`
}
```

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET    | `/health` | Health check |
| POST   | `/tickets` | Create ticket (`?auto_classify=true` optional) |
| POST   | `/tickets/import` | Bulk import multipart file (CSV/JSON/XML) |
| GET    | `/tickets` | List tickets (filter: status, category, priority, customer_id) |
| GET    | `/tickets/:id` | Get ticket by ID |
| PUT    | `/tickets/:id` | Update ticket |
| DELETE | `/tickets/:id` | Delete ticket |
| POST   | `/tickets/:id/auto-classify` | Run auto-classification on a ticket |

---

## Implementation Phases

### Phase 1: Project Scaffold
- Create `src/` with `go.mod` (module `support-tickets`, go 1.23)
- `go get` dependencies: fiber/v2, gofiber/swagger, google/uuid, testify, swaggo/swag
- Dockerfile (multi-stage) and docker-compose.yml copied/adapted from homework-1
- `scripts/generate-swagger.sh`

### Phase 2: Domain Layer
- `internal/domain/ticket.go` — types, enum constants, `NewTicket()` constructor
- `internal/domain/errors.go` — `ErrNotFound`, `ErrInvalidInput`, `ErrValidation{Details []ValidationDetail}`

### Phase 3: Ports
- `internal/ports/in/ticket_service.go` — `TicketService` interface: `Create`, `BulkImport`, `List`, `GetByID`, `Update`, `Delete`, `AutoClassify`
- `internal/ports/out/ticket_repository.go` — `TicketRepository` interface: `Save`, `FindByID`, `FindAll`, `Update`, `Delete`

### Phase 4: Importers (`pkg/importer/`)
- `importer.go` — `TicketImporter` interface, `ImportRecord`, `ImportResult{Total, Successful, Failed int, Errors []ImportError}`
- `csv.go` — `encoding/csv`, header-based column mapping
- `json.go` — `encoding/json`, expects `[]ImportRecord` array
- `xml.go` — `encoding/xml`, `<tickets><ticket>...</ticket></tickets>` structure
- Each parser collects per-row errors; malformed files return partial success

### Phase 5: Application Service
- `internal/service/validator.go` — email regex, subject 1–200, description 10–2000, enum validation
- `internal/service/classifier.go` — keyword tables per category and priority; confidence = `matched / len(keywords)` capped at 1.0

  **Keyword tables:**
  ```
  account_access:   login, password, 2fa, access, account, locked out, sign in
  technical_issue:  error, crash, bug, broken, not working, 500, exception, fail
  billing_question: payment, invoice, refund, charge, billing, subscription
  feature_request:  feature, enhancement, suggestion, request, would like
  bug_report:       reproduce, steps to, expected, actual, defect, regression

  urgent: can't access, critical, production down, security, outage
  high:   important, blocking, asap, immediately
  low:    minor, cosmetic, nice to have, when possible
  medium: (default when no priority keywords match)
  ```

- `internal/service/ticket_service.go` — orchestrates repo + validator + classifier; `BulkImport` calls the appropriate importer then validates+saves each record

### Phase 6: In-Memory Repository
- `internal/adapters/out/memory/ticket_repository.go`
- `sync.RWMutex`-protected `map[uuid.UUID]*domain.Ticket`
- `FindAll` accepts a `TicketFilter` struct and applies filters in-memory

### Phase 7: HTTP Handlers + Swagger
- `internal/adapters/in/http/handler.go` — bind request → call service → return JSON; no business logic
- `POST /tickets/import` — reads multipart file, detects format by Content-Type or file extension
- Swagger `//` annotations on every handler method
- `swagger_models.go` — response wrapper types for Swagger only

### Phase 8: Wiring
- `cmd/api/main.go` — construct `memory.Repo → service → handler → server.Run()`
- `cmd/api/server.go` — Fiber config with global error handler, middleware stack, route registration
- Run `./scripts/generate-swagger.sh` to produce `docs/swagger.yaml`

### Phase 9: Sample Data
- `sample_data/sample_tickets.csv` — 50 rows, varied categories/priorities/statuses
- `sample_data/sample_tickets.json` — 20 ticket objects
- `sample_data/sample_tickets.xml` — 30 tickets under `<tickets>` root
- `sample_data/invalid_*.{csv,json,xml}` — bad email, missing fields, invalid enum values

### Phase 10: AI-Generated Test Suite (target >85% coverage)

Tests live in `src/internal/` and `src/pkg/` alongside the code they cover (Go convention).

| Test File | Group | Count | What it covers |
|-----------|-------|-------|----------------|
| `adapters/in/http/handler_test.go` | test_ticket_api | 14 | Every endpoint via `app.Test()` — 2xx and 4xx paths |
| `domain/ticket_test.go` | test_ticket_model | 9 | Entity construction, all enum validators |
| `pkg/importer/csv_test.go` | test_import_csv | 6 | Valid CSV, missing headers, bad rows, empty file |
| `pkg/importer/json_test.go` | test_import_json | 5 | Valid JSON, malformed JSON, empty array |
| `pkg/importer/xml_test.go` | test_import_xml | 5 | Valid XML, malformed XML, bad structure |
| `service/classifier_test.go` | test_categorization | 12 | Each category, each priority, confidence bounds |
| `service/ticket_service_test.go` | test_integration | 10 | CRUD, BulkImport summary counts, AutoClassify |
| `service/validator_test.go` | test_ticket_model | 10 | All field validation rules |
| `adapters/out/memory/ticket_repository_test.go` | — | 8 | Concurrent reads/writes, filter combos |
| `cmd/api/server_test.go` | — | 2 | Server wiring |
| `service/benchmarks_test.go` | test_performance | 6 | Benchmarks for validator, classifier, importer, service |

**Fixtures** reuse `sample_data/` files (CSV 50, JSON 20, XML 30, plus invalid variants).

**Coverage achieved:** overall 89.1% ✓

### Phase 11: Integration & Performance Tests (`tests/` package)

Black-box tests in `src/tests/` that treat the system as a whole. All tests use a real in-memory app instance via `newApp()`.

| Test File | Group | Count | What it covers |
|-----------|-------|-------|----------------|
| `tests/test_ticket_api_test.go` | test_ticket_api | 11 | HTTP endpoints, 2xx and 4xx paths |
| `tests/test_ticket_model_test.go` | test_ticket_model | 9 | Domain entity, all enums, typed errors |
| `tests/test_import_csv_test.go` | test_import_csv | 6 | CSV fixture file, field mapping, edge cases |
| `tests/test_import_json_test.go` | test_import_json | 5 | JSON fixture file, field mapping, edge cases |
| `tests/test_import_xml_test.go` | test_import_xml | 5 | XML fixture file, field mapping, edge cases |
| `tests/test_categorization_test.go` | test_categorization | 10 | All categories, all priorities, confidence/reasoning |
| `tests/test_integration_test.go` | test_integration | 8 | Lifecycle, classify, concurrent (25 reqs), combined filter |
| `tests/test_performance_test.go` | test_performance | 7 | Benchmarks: validator, classifier, importers, HTTP, concurrent, list-with-filter |

**New integration scenarios (Phase 11):**
- `TestIntegration_BulkImportAndAutoClassify` — imports CSV batch, runs auto-classify on imported ticket, verifies classification persisted
- `TestIntegration_ConcurrentCreate` — 25 goroutines POST simultaneously; asserts all 25 stored
- `TestIntegration_FilterByCategoryAndPriority` — seeds 5 tickets across two categories/priorities, validates all four filter combos

**New performance benchmarks (Phase 11):**
- `BenchmarkPerf_ConcurrentHTTPCreate` — `b.RunParallel` stress across all CPU cores
- `BenchmarkPerf_ListWithFilter` — GET with `category+priority` query on 50-ticket dataset

---

## Go Module

```
module support-tickets
go 1.23
```

Direct dependencies:
- `github.com/gofiber/fiber/v2`
- `github.com/gofiber/swagger`
- `github.com/google/uuid`
- `github.com/stretchr/testify`
- `github.com/swaggo/swag` (CLI, dev only)

---

## Verification

```bash
cd homework-2/src
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total    # expect ≥ 85%
go tool cover -html=coverage.out -o coverage.html

docker compose up --build                        # API on :8080
docker compose --profile docs up --build         # + Swagger UI on :8081
curl http://localhost:8080/health
curl -X POST http://localhost:8080/tickets/import \
     -F "file=@../sample_data/sample_tickets.csv"
```
