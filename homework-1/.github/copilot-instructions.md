# Banking Transactions API Instructions

## 1. Prerequisites
- Docker and Docker Compose for container run
- No Go installation is required for Docker-only workflow
- Go 1.23+ only if running locally or regenerating Swagger docs
- curl for command-line testing (optional)

## 2. Quick Start (Docker)
From project root:

```bash
cd homework-1
docker compose up --build
```

API base URL:
- http://localhost:8080

Run with Swagger UI profile:

```bash
cd homework-1
docker compose --profile docs up --build
```

Swagger UI URL (docs profile):
- http://localhost:8081

Stop services:

```bash
docker compose down
```

## 3. Quick Start (Local)
From project root:

```bash
cd homework-1/src
go run ./cmd/api
```

Default URL:
- http://localhost:8080

## 4. Environment Variables
- PORT: HTTP port (default: 8080)
- LOG_LEVEL: debug | info | warn | error (default: info)

Examples:

```bash
PORT=9090 LOG_LEVEL=debug go run ./cmd/api
```

## 5. Verify the API
Health check:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

Create transaction:

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "fromAccount":"ACC-12345",
    "toAccount":"ACC-67890",
    "amount":100.50,
    "currency":"USD",
    "type":"transfer"
  }'
```

List transactions:

```bash
curl http://localhost:8080/transactions
```

Filter transactions:

```bash
curl "http://localhost:8080/transactions?accountId=ACC-12345&type=transfer"
```

Get account balance:

```bash
curl http://localhost:8080/accounts/ACC-12345/balance
```

Get account summary:

```bash
curl http://localhost:8080/accounts/ACC-12345/summary
```

## 6. Validation Checks
Example invalid payload:

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "fromAccount":"INVALID",
    "toAccount":"ACC-67890",
    "amount":-10,
    "currency":"FAKE",
    "type":"transfer"
  }'
```

Expected behavior:
- HTTP 422
- JSON with error, request_id, and details array

## 7. Development Workflow
Run formatter and checks:

```bash
cd homework-1/src
go fmt ./...
go test ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Primary implementation locations:
- src/cmd/api: application bootstrap and route registration
- src/internal/service: business logic and validation integration
- src/internal/adapters/in/http: HTTP handlers
- src/internal/adapters/out/memory: in-memory repository
- src/pkg/errorhandler: global error handling
- src/pkg/middleware: request ID and logging middleware
- scripts/generate-swagger.sh: generate OpenAPI spec from annotations
- docs/swagger.yaml: generated OpenAPI spec served by Swagger UI in compose
  (when docs profile is enabled)

## 8. API Documentation Workflow
Generate docs from code annotations:

```bash
cd homework-1
./scripts/generate-swagger.sh
```

Notes:
- Do not manually maintain docs/swagger.yaml as source of truth.
- Source of truth is Swagger annotations in Go files.
- Regenerate docs/swagger.yaml after endpoint or schema changes.
- Running this generation command requires Go.

## 9. Manual Test Workflow
1. Start API with Docker or local Go.
2. Call health endpoint.
3. Create deposit, transfer, and withdrawal transactions.
4. Query transactions with and without filters.
5. Query balance and summary for involved accounts.
6. Submit invalid payload and verify 422 format.
7. Request unknown transaction ID and verify 404.
8. Open Swagger UI and verify endpoint docs are visible.

## 10. Troubleshooting
- Port already in use:
  - Change PORT env var (local), or stop conflicting service.
- Docker build issues:
  - Run docker compose build --no-cache.
- Swagger UI not showing new endpoints:
  - Re-run ./scripts/generate-swagger.sh and restart docker compose.
- Empty account balance/summary returns 404:
  - Create at least one transaction involving that account.
- Date filter rejected:
  - Ensure date format is YYYY-MM-DD.

## 11. Supporting Demo Assets
- demo/sample-requests.http: full API request set for REST Client
- demo/sample-data.json: sample data
- demo/run.sh: helper script for container run flow
