# How to Run

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose (for container run)
- No Go installation is required if you run with Docker Compose only
- [Go 1.23+](https://go.dev/dl/) only if:
  - you run the API locally, or
  - you regenerate Swagger docs from code annotations

---

## Run with Docker (recommended)

```bash
cd homework-1
docker compose up --build
```

The API will be available at `http://localhost:8080`.

Run with Swagger UI too:

```bash
cd homework-1
docker compose --profile docs up --build
```

With docs profile enabled:
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8081`

To stop:
```bash
docker compose down
```

---

## Run locally

```bash
cd homework-1/src
go run ./cmd/api
```

Optional environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP port |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |

---

## Quick test

```bash
# Health check
curl http://localhost:8080/health

# Create a transaction
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"fromAccount":"ACC-12345","toAccount":"ACC-67890","amount":100.50,"currency":"USD","type":"transfer"}'

# List all transactions
curl http://localhost:8080/transactions

# Get account balance
curl http://localhost:8080/accounts/ACC-12345/balance
```

See `demo/sample-requests.http` for a full set of sample requests (VS Code REST Client compatible).

---

## Swagger API docs

Generate OpenAPI spec from code annotations:

```bash
cd homework-1
./scripts/generate-swagger.sh
```

Note:
- This generation step requires Go on your machine.
- If you run only with Docker Compose and do not need to regenerate docs,
  Go is not required.

Generated file:
- `docs/swagger.yaml`

When running with Docker Compose, Swagger UI is started only with profile `docs`.

Run Swagger UI only (without compose) with Docker:

```bash
cd homework-1
docker run --rm -p 8081:8080 \
  -e SWAGGER_JSON=/app/swagger.yaml \
  -v "$PWD/docs/swagger.yaml:/app/swagger.yaml:ro" \
  swaggerapi/swagger-ui
```

Then open:
- `http://localhost:8081`

Stop Swagger UI:
- Press `Ctrl+C` in the terminal running the container
