# 🏦 Homework 1: Banking Transactions API

> **Student Name**: Roman Reznik
> **Date Submitted**: 26-04-2026
> **AI Tools Used**: [Claude Code, GitHub Copilot]

---

## 📋 Project Overview

I built a **Banking Transactions REST API** using **Go + Fiber** with a clean
**Hexagonal (Ports & Adapters)** architecture and in-memory storage.

The API supports creating transactions, listing/filtering transaction history,
retrieving a transaction by ID, and computing account balance/summary.

The project also includes:
- ✅ Centralized error handling with consistent JSON error contract
- ✅ Request tracing with `X-Request-ID`
- ✅ Structured JSON logging via `slog`
- ✅ Unit tests with `testify` and high coverage
- ✅ OpenAPI/Swagger documentation generated from code annotations
- ✅ Docker Compose setup for both API and Swagger UI

---

## 🚀 Features Implemented

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/transactions` | Create a new transaction |
| `GET` | `/transactions` | List transactions (supports filters) |
| `GET` | `/transactions/:id` | Get transaction by ID |
| `GET` | `/accounts/:accountId/balance` | Get account balance |
| `GET` | `/accounts/:accountId/summary` | Get account summary |
| `GET` | `/health` | Health check |

### Validation Rules
- Amount must be positive
- Amount supports max 2 decimal places
- Account format must match `ACC-XXXXX`
- Currency must be from ISO 4217 allowlist
- Type must be one of: `deposit`, `withdrawal`, `transfer`

### Transaction Filters
- `accountId`
- `type`
- `from` date (`YYYY-MM-DD`)
- `to` date (`YYYY-MM-DD`)
- Combined filtering supported

### Documentation & DevOps
- OpenAPI spec generated via `./scripts/generate-swagger.sh`
- Generated spec location: `docs/swagger.yaml`
- Docker Compose runs:
	- API at `http://localhost:8080` (default)
	- Swagger UI at `http://localhost:8081` (with `docs` profile)

---

## 🧱 Architecture Decisions

### 1) Hexagonal Architecture (Ports & Adapters)
**Decision**: Separate domain logic from delivery/storage layers.

**Why**:
- Keeps business logic independent from framework specifics
- Makes testing easier and cleaner
- Enables easier replacement of adapters in future

### 2) In-Memory Repository for Homework Scope
**Decision**: Use `sync.RWMutex` protected map as storage.

**Why**:
- Matches homework requirement (no DB)
- Simple and fast for local testing/demo

**Trade-off**:
- Data is not persistent across restarts

### 3) Centralized Error Handling
**Decision**: Route all handler errors through global error handler.

**Why**:
- Consistent response contract (`error`, `request_id`, `details`)
- Cleaner handlers with less repetitive code

### 4) Structured Logging + Request Correlation
**Decision**: Use `slog` JSON logs + request ID middleware.

**Why**:
- Easier debugging and observability
- Better traceability in container environments

### 5) Swagger Generated from Code Annotations
**Decision**: Treat annotations in Go code as source of truth.

**Why**:
- Reduces drift between implementation and documentation
- Keeps documentation maintenance repeatable and automated

---

## 🧪 Testing Summary

- Unit tests added across service, handlers, repository, middleware, logger, and error handler
- Assertions use `testify` (`assert` / `require`)
- Coverage improved to target level (>= 85%)

---

## 🤖 AI Workflow And Evidence

### Claude Code (Sonnet 4.6)

| Point | Activity | Screenshot Evidence |
|------|----------|---------------------|
| 1 | Create plan | [1_Plan.jpg](docs/screenshots/1_Plan.jpg) |
| 1.1 | Edit the created plan | [1_1Plan_Edit.jpg](docs/screenshots/1_1Plan_Edit.jpg) |
| 2 | Implement the created plan | [2_Plan_Implementation.jpg](docs/screenshots/2_Plan_Implementation.jpg) |
| 2.1 | Review implementation and edit requests | [2_1_Review_and_Edit.jpg](docs/screenshots/2_1_Review_and_Edit.jpg) |
| 3 | Check the result | [3_Run_Check_results.jpg](docs/screenshots/3_Run_Check_results.jpg), [3_1_Check_results.jpg](docs/screenshots/3_1_Check_results.jpg), [3_2_Check_results.jpg](docs/screenshots/3_2_Check_results.jpg) |
| 4 | Create claude init context | [4_Create_claude_init.jpg](docs/screenshots/4_Create_claude_init.jpg) |

### Co-Pilot (GPT-5.3 Codex)

| Point | Activity | Screenshot Evidence |
|------|----------|---------------------|
| 5 | Create specification and instructions files | [5_Co-Pilot_Create_instructions_specifications.jpg](docs/screenshots/5_Co-Pilot_Create_instructions_specifications.jpg) |
| 5.1 | Verify and fix result | [5_1_Co-Pilot_Check_Edit_instructions_specifications.jpg](docs/screenshots/5_1_Co-Pilot_Check_Edit_instructions_specifications.jpg) |
| 6 | Add unit tests | [6_Co-Pilot_Create_unitTests.jpg](docs/screenshots/6_Co-Pilot_Create_unitTests.jpg), [6_1_Co-Pilot_UnitTests_improve_coverage.jpg](docs/screenshots/6_1_Co-Pilot_UnitTests_improve_coverage.jpg) |
| 7 | Create Swagger API documentation | [7_Co-Pilot_Create_Swagger_docs.jpg](docs/screenshots/7_Co-Pilot_Create_Swagger_docs.jpg), [7_1_Co-Pilot_Result_Swagger.jpg](docs/screenshots/7_1_Co-Pilot_Result_Swagger.jpg) |
| 8 | Update README file | This README update |

---

## ▶️ Quick Start

```bash
cd homework-1
docker compose up --build
```

- API: `http://localhost:8080`

For full run instructions, see [HOWTORUN.md](HOWTORUN.md).


<div align="center">

*This project was completed as part of the AI-Assisted Development course.*

</div>
