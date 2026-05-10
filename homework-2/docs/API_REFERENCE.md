# API Reference

Base URL: `http://localhost:8080`

All request and response bodies use `application/json` unless noted otherwise. Timestamps are RFC 3339 UTC.

---

## Data Models

### Ticket

```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "customer_id": "CUST-001",
  "customer_email": "user@example.com",
  "customer_name": "Alice Smith",
  "subject": "Cannot login to account",
  "description": "Getting password error every time I try to sign in.",
  "category": "account_access",
  "priority": "high",
  "status": "new",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "resolved_at": null,
  "assigned_to": null,
  "tags": [],
  "metadata": {
    "source": "web_form",
    "browser": "Chrome",
    "device_type": "desktop"
  },
  "classification": null
}
```

**Field constraints:**

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| `customer_id` | string | yes | non-empty |
| `customer_email` | string | yes | valid email format |
| `customer_name` | string | yes | non-empty |
| `subject` | string | yes | 1–200 characters |
| `description` | string | yes | 10–2000 characters |
| `category` | string | yes | see Category enum |
| `priority` | string | yes | see Priority enum |
| `metadata.source` | string | no | see Source enum |
| `metadata.device_type` | string | no | see DeviceType enum |

### Enums

| Type | Values |
|------|--------|
| Category | `account_access`, `technical_issue`, `billing_question`, `feature_request`, `bug_report`, `other` |
| Priority | `urgent`, `high`, `medium`, `low` |
| Status | `new`, `in_progress`, `waiting_customer`, `resolved`, `closed` |
| Source | `web_form`, `email`, `api`, `chat`, `phone` |
| DeviceType | `desktop`, `mobile`, `tablet` |

### Classification

Returned by `POST /tickets/:id/auto-classify` and stored on the ticket.

```json
{
  "category": "billing_question",
  "priority": "medium",
  "confidence": 0.71,
  "reasoning": "Matched 5 of 7 billing_question keywords",
  "keywords_found": ["refund", "charge", "billing", "invoice", "payment"]
}
```

### BulkImportResult

Returned by `POST /tickets/import`.

```json
{
  "Total": 50,
  "Successful": 48,
  "Failed": 2,
  "Errors": [
    { "row": 3, "error": "invalid email format" },
    { "row": 17, "error": "description too short (min 10 chars)" }
  ]
}
```

### Error Response

```json
{
  "error": "ticket not found",
  "request_id": "req-abc123",
  "details": [
    { "field": "customer_email", "message": "invalid email format" }
  ]
}
```

`details` is present only for validation errors (HTTP 400).

---

## Endpoints

### GET /health

Health check. Returns `200 OK` when the service is running.

**Response**
```json
{ "status": "ok" }
```

**cURL**
```bash
curl http://localhost:8080/health
```

---

### POST /tickets

Create a new support ticket. Optionally auto-classify it in the same request.

**Query parameters**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `auto_classify` | boolean | `false` | Run auto-classification after creation |

**Request body**
```json
{
  "customer_id": "CUST-001",
  "customer_email": "alice@example.com",
  "customer_name": "Alice Smith",
  "subject": "Cannot login to account",
  "description": "Getting password error every time I try to sign in. Account is locked.",
  "category": "account_access",
  "priority": "high",
  "metadata": {
    "source": "web_form",
    "browser": "Chrome",
    "device_type": "desktop"
  }
}
```

**Responses**

| Status | Description |
|--------|-------------|
| `201 Created` | Ticket created; body contains the new `Ticket` object |
| `400 Bad Request` | Validation failure; body contains `ErrorResponse` with `details` |

**cURL — basic create**
```bash
curl -X POST http://localhost:8080/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "CUST-001",
    "customer_email": "alice@example.com",
    "customer_name": "Alice Smith",
    "subject": "Cannot login to account",
    "description": "Getting password error every time I try to sign in. Account is locked.",
    "category": "account_access",
    "priority": "high"
  }'
```

**cURL — create with auto-classification**
```bash
curl -X POST "http://localhost:8080/tickets?auto_classify=true" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "CUST-002",
    "customer_email": "bob@example.com",
    "customer_name": "Bob Jones",
    "subject": "Charged twice this month",
    "description": "I was charged twice for my subscription. Need a refund for the duplicate billing.",
    "category": "other",
    "priority": "medium"
  }'
```

---

### POST /tickets/import

Bulk import tickets from a file. Accepts CSV, JSON, or XML detected by filename extension.

**Content-Type:** `multipart/form-data`

**Form field:** `file` — the file to import

**Supported formats:**

| Extension | Expected structure |
|-----------|--------------------|
| `.csv` | Header row + data rows; column names match Ticket fields |
| `.json` | JSON array of ticket objects |
| `.xml` | `<tickets><ticket>…</ticket></tickets>` |

**Response (200 OK)**
```json
{
  "Total": 50,
  "Successful": 48,
  "Failed": 2,
  "Errors": [
    { "row": 3, "error": "invalid email format" },
    { "row": 17, "error": "description too short (min 10 chars)" }
  ]
}
```

The endpoint always returns `200 OK` and a summary — even when some rows fail. A `400` is returned only if the file itself cannot be read or parsed at all.

**cURL — CSV import**
```bash
curl -X POST http://localhost:8080/tickets/import \
  -F "file=@sample_data/sample_tickets.csv"
```

**cURL — JSON import**
```bash
curl -X POST http://localhost:8080/tickets/import \
  -F "file=@sample_data/sample_tickets.json"
```

**cURL — XML import**
```bash
curl -X POST http://localhost:8080/tickets/import \
  -F "file=@sample_data/sample_tickets.xml"
```

---

### GET /tickets

List all tickets with optional filtering.

**Query parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | string | Filter by `Status` enum value |
| `category` | string | Filter by `Category` enum value |
| `priority` | string | Filter by `Priority` enum value |
| `customer_id` | string | Filter by exact customer ID |

**Response (200 OK):** array of `Ticket` objects (empty array `[]` when none match)

**cURL — all tickets**
```bash
curl http://localhost:8080/tickets
```

**cURL — filter by category and priority**
```bash
curl "http://localhost:8080/tickets?category=technical_issue&priority=high"
```

**cURL — filter by status and customer**
```bash
curl "http://localhost:8080/tickets?status=new&customer_id=CUST-001"
```

---

### GET /tickets/:id

Get a single ticket by its UUID.

**Path parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID string | Ticket ID |

**Responses**

| Status | Description |
|--------|-------------|
| `200 OK` | Body contains `Ticket` |
| `404 Not Found` | No ticket with that ID |

**cURL**
```bash
curl http://localhost:8080/tickets/3fa85f64-5717-4562-b3fc-2c963f66afa6
```

---

### PUT /tickets/:id

Update one or more fields on an existing ticket. All fields are optional; only provided fields are updated.

**Request body (all fields optional)**
```json
{
  "status": "in_progress",
  "assigned_to": "agent@company.com",
  "priority": "urgent",
  "tags": ["vip", "escalated"]
}
```

**Responses**

| Status | Description |
|--------|-------------|
| `200 OK` | Body contains updated `Ticket` |
| `400 Bad Request` | Malformed body |
| `404 Not Found` | Ticket not found |

**cURL — assign and escalate**
```bash
curl -X PUT http://localhost:8080/tickets/3fa85f64-5717-4562-b3fc-2c963f66afa6 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "assigned_to": "agent@company.com",
    "priority": "urgent"
  }'
```

**cURL — resolve a ticket**
```bash
curl -X PUT http://localhost:8080/tickets/3fa85f64-5717-4562-b3fc-2c963f66afa6 \
  -H "Content-Type: application/json" \
  -d '{ "status": "resolved" }'
```

---

### DELETE /tickets/:id

Permanently delete a ticket.

**Responses**

| Status | Description |
|--------|-------------|
| `204 No Content` | Deleted successfully; empty body |
| `404 Not Found` | Ticket not found |

**cURL**
```bash
curl -X DELETE http://localhost:8080/tickets/3fa85f64-5717-4562-b3fc-2c963f66afa6
```

---

### POST /tickets/:id/auto-classify

Run the keyword-based classifier on the ticket's subject and description. Updates and persists the `classification` field on the ticket, then returns the `Classification` object.

**Responses**

| Status | Description |
|--------|-------------|
| `200 OK` | Body contains `Classification` |
| `404 Not Found` | Ticket not found |

**cURL**
```bash
curl -X POST \
  http://localhost:8080/tickets/3fa85f64-5717-4562-b3fc-2c963f66afa6/auto-classify
```

**Example response**
```json
{
  "category": "account_access",
  "priority": "high",
  "confidence": 0.57,
  "reasoning": "Matched 4 of 7 account_access keywords",
  "keywords_found": ["login", "password", "account", "access"]
}
```

---

## Error Reference

| HTTP Status | Cause |
|-------------|-------|
| `400 Bad Request` | Validation failure (missing required field, invalid enum, wrong length) |
| `404 Not Found` | Ticket ID does not exist |
| `500 Internal Server Error` | Unexpected server error |

All error bodies follow the `ErrorResponse` schema. The `request_id` field maps to the `X-Request-ID` response header and can be used to correlate logs.

---

## Swagger UI

An interactive Swagger UI is available when the service is running:

```
http://localhost:8081
```

(Requires `docker compose --profile docs up`)

The raw OpenAPI spec is at `docs/swagger.yaml`.
