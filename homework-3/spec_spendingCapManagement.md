# Spending Cap Management — Specification

| Field | Value |
|-------|-------|
| Feature | Spending Cap Management |
| Version | 1.0 |
| Status | Draft for engineering review |
| Regulatory level | Basic (audit log + PII/PAN masking) |
| Primary stakeholders | End-user, Support agent, Fraud team, Ops/Compliance |

---

## High-Level Objective

- Let an account holder set, view, change, and remove **spending caps** (limits on how much can be spent on a card or account over a defined window) so that overspend is prevented at authorization time, while every change is auditable and visible to support, fraud, and compliance.
- **Scope boundary**: This feature manages cap *configuration and enforcement signals* only. It does **not** own the payment authorization engine, ledger, or notification delivery — it exposes a decision (`ALLOW`/`DENY`) and emits events that those systems consume.

---

## Mid-Level Objectives

Each objective is observable — there is a measurable change in the world when it succeeds.

1. **MO-1 — Cap lifecycle**: A user can create, read, update, and delete caps scoped to an account or a specific card. After a successful write, the new cap is returned and persisted.
2. **MO-2 — Multi-dimensional caps**: The system supports caps by **period** (daily / weekly / monthly / per-transaction) and by **category** (e.g. all, ATM, online, merchant-category). Multiple caps can apply to one card simultaneously.
3. **MO-3 — Real-time enforcement decision**: Given a proposed transaction amount and context, the system returns an `ALLOW` or `DENY` decision reflecting all applicable caps and current cumulative spend within each window.
4. **MO-4 — Spend accumulation & windows**: Approved spend is accumulated per cap window; windows roll over correctly at period boundaries in the account's configured timezone.
5. **MO-5 — Auditability**: Every cap mutation and every `DENY` decision produces an immutable audit record attributing actor, before/after state, timestamp, and reason.
6. **MO-6 — Role-scoped visibility**: End-users see and edit only their own caps; support agents can view (not edit) with masked identifiers; fraud and compliance can view all caps and full audit history.
7. **MO-7 — Fraud guardrails**: The fraud team can place a **temporary hard cap** (or zero-cap freeze) on a card that the end-user cannot override or remove until released.

---

## Non-Functional & Policy Requirements

### Security & Privacy
- **NFR-SEC-1**: PAN (card number) is never stored or logged in this service; caps reference cards by an internal `card_id` / tokenized reference only.
- **NFR-SEC-2**: Any identifier surfaced to a support agent is **masked** (e.g. card shown as `•••• 4321`, user email as `u***@gmail.com`).
- **NFR-SEC-3**: All write endpoints require authenticated, authorized actors; authorization is enforced server-side per-resource (deny by default).
- **NFR-SEC-4**: All monetary values use a fixed-precision decimal type with explicit currency; floating-point money is prohibited.

### Auditability (Basic regulatory level)
- **NFR-AUD-1**: Audit records are **append-only** — no update or delete path exists in code or schema.
- **NFR-AUD-2**: Each audit record carries: `event_id`, `actor_id`, `actor_role`, `action`, `target_card_id`, `before`, `after`, `reason`, `request_id`, `created_at` (UTC).
- **NFR-AUD-3**: Audit writes are part of the same transaction as the state change; a failed audit write rolls back the change.

### Reliability
- **NFR-REL-1**: The enforcement decision path must **fail safe**. On internal error or datastore unavailability the configured policy applies a deterministic fallback (default: `DENY` for fraud freezes, `ALLOW` with flag for ordinary caps — see Implementation Notes IN-7).
- **NFR-REL-2**: Spend accumulation is idempotent per `transaction_id`; replaying the same approved transaction must not double-count.

### Performance (assumed targets — see README for justification)
- **NFR-PERF-1**: Enforcement decision (`evaluate`) **p95 ≤ 50 ms**, **p99 ≤ 120 ms** server-side. *Rationale: it sits inline in the card-authorization path, which typically has a sub-second total budget; cap evaluation must be a small slice of it.*
- **NFR-PERF-2**: Cap CRUD endpoints **p95 ≤ 300 ms**.
- **NFR-PERF-3**: Listing caps is paginated, **max page size 100**, default 25.
- **NFR-PERF-4**: Read-after-write consistency for a user's own caps ≤ 1 s (a user who saves a cap and refreshes sees it).
- **NFR-PERF-5**: `evaluate` sustains **≥ 500 req/s** per instance at the latency targets above.

---

## Implementation Notes

- **IN-1 — Money**: Represent money as a `Money` value object — `shopspring/decimal` amount + ISO-4217 currency, constructed from string (never `float64`). Reject mixed-currency comparisons; a cap and the transaction it gates must share currency or be converted upstream (out of scope here — assume same currency).
- **IN-2 — IDs**: All resource IDs are opaque UUIDv4 strings. Never expose database primary keys or sequential integers.
- **IN-3 — Idempotency**: All mutating endpoints accept an `Idempotency-Key` header; repeated keys return the original result without re-applying.
- **IN-4 — Time & windows**: Window boundaries are computed in the **account's timezone** but persisted in UTC. "Daily" = local calendar day; "monthly" = local calendar month. Document the boundary rule explicitly in code.
- **IN-5 — Cap precedence**: When multiple caps apply, **the most restrictive applicable cap wins**. A fraud hard-cap always overrides user caps.
- **IN-6 — Error semantics**: Use typed domain errors → stable HTTP codes: validation `422`, auth `401/403`, not-found `404`, conflict/idempotency replay `409`, cap-blocked decision is `200` with body `{decision: "DENY"}` (a deny is a *valid answer*, not an HTTP error).
- **IN-7 — Fail-safe matrix**: Persist enough cap state in a fast cache that `evaluate` can still apply fraud freezes during a primary-DB outage. Fallback policy is config-driven, defaults in NFR-REL-1.
- **IN-8 — No PII in logs**: Application logs reference `card_id`/`user_id` only. A central masking helper is the single allowed path for rendering any identifier to a human.
- **IN-9 — Validation rules**: Cap amount must be `> 0`; period must be in the allowed enum; per-transaction cap cannot exceed a daily cap on the same card/category (warn, don't silently accept). Removing a cap requires the cap not be a fraud hard-cap.

---

## Context

### Beginning context
- `card_id` and `user_id` exist in an upstream account service (assumed; accessed through an outbound port, not owned here).
- Empty Go module (`go.mod`) for the spending-cap service: no domain, ports, adapters, or tests yet.
- A shared structured logger exists but has **no** masking helper yet.

### Ending context

The service follows the hexagonal layout defined in `CLAUDE.md` (Architecture) — that file is the single source of truth for the directory structure; it is not repeated here. After the Low-Level Tasks are complete, the following must exist:

- A buildable Go module exposing the REST API (cap CRUD, `evaluate`, `fraud hold`, `audit read`).
- Domain, ports, service, and adapter layers wired together in `cmd/api`, with all business logic behind ports (no logic in HTTP handlers).
- `pkg/masking` as the single masking helper, used by every human-facing response.
- An append-only audit store plus `docs/audit-fields.md` listing every audited field (NFR-AUD-2).
- A test suite (unit, integration, decision-matrix, concurrency, failure-injection) and a benchmark, with fixtures covering edge cases EC-1…EC-12.

---

## Edge Cases & Failure Modes

| # | Scenario | Expected behavior | Audit/compliance implication |
|---|----------|-------------------|------------------------------|
| EC-1 | User sets cap below current period spend | Accept cap; remaining = 0; subsequent spend denied | Audit the new cap; not an error |
| EC-2 | Transaction exactly equals remaining cap | `ALLOW` (≤ is allowed; `>` denies) | Normal accumulation |
| EC-3 | Two concurrent transactions race the last remaining amount | Serialize per card window; only one wins, other `DENY` | No double-spend; both attempts evaluable |
| EC-4 | Replay of an already-counted `transaction_id` | No double count; return prior result | Idempotency (NFR-REL-2) |
| EC-5 | User tries to remove/override a fraud hard-cap | `403`; cap unchanged | Audit the *attempt* with reason `FRAUD_HOLD_LOCKED` |
| EC-6 | Cap amount ≤ 0 or unknown period | `422` validation error | Rejected before persistence |
| EC-7 | Window boundary crossed mid-evaluation | Use window for the transaction's effective time; old window's spend not carried over | Correct rollover |
| EC-8 | Primary datastore unavailable during `evaluate` | Apply fail-safe matrix (IN-7); flag decision `degraded: true` | Degraded decisions are audited for later review |
| EC-9 | Support agent opens a card they shouldn't access | `403`; nothing rendered | Access attempt logged |
| EC-10 | Stale cached cap after a recent update | Read-after-write ≤ 1 s (NFR-PERF-4); evaluate prefers authoritative store for fraud holds | Avoid enforcing a removed/changed cap |
| EC-11 | Fraud places zero-cap freeze while a transaction is mid-flight | Freeze applies to all *new* evaluations immediately | Audit freeze with actor = fraud agent |
| EC-12 | Empty state: card has no caps | `evaluate` returns `ALLOW` (nothing restricts it) | No audit needed for allow-by-absence |

---

## Verification

How we know each Mid-Level Objective is met.

| Objective | Verification method |
|-----------|---------------------|
| MO-1 | Integration tests: create→read→update→delete returns expected states; persistence survives restart fixture |
| MO-2 | Unit tests asserting multiple simultaneous caps stored and retrieved per scope/period/category |
| MO-3 | Decision tests across a fixture matrix of (amount, applicable caps, prior spend) → expected `ALLOW`/`DENY` |
| MO-4 | Window tests at period boundaries across timezones; rollover fixtures (DST day included) |
| MO-5 | Audit tests: every mutation and every deny yields exactly one append-only record with required fields (NFR-AUD-2) |
| MO-6 | Authorization tests per role: user/support/fraud/compliance matrix of allowed actions; masking asserted in support responses |
| MO-7 | Tests proving user cannot remove/override a fraud hold; only fraud/compliance can release |

**Review checkpoints**:

- [ ] Compliance review of `docs/audit-fields.md` against NFR-AUD-2.
- [ ] Security review confirming no PAN/PII path into logs.
- [ ] Manual masking spot-check of every support-facing response shape.

**Test categories required**:

- [ ] Unit (models, precedence, masking, windows).
- [ ] Integration (CRUD + audit transaction atomicity).
- [ ] e2e/decision (evaluate matrix).
- [ ] Concurrency (EC-3).
- [ ] Failure injection (EC-8).

---

## Low-Level Tasks

Each task names the objective it serves and ends with acceptance criteria.

### 1. Define domain models (MO-1, MO-2)
**Prompt**: "Create domain types for Cap, CapPeriod (enum daily/weekly/monthly/per_transaction), CapScope (account/card + optional category), and FraudHold, using a Money value object backed by shopspring/decimal with currency. Use constructors that validate, keep fields unexported for immutability."


**File**: `internal/domain/cap.go`, `internal/domain/money.go`

**Function/class**: `Cap`, `CapPeriod`, `CapScope`, `FraudHold`, `Money`, `NewMoney`, `NewCap`

**Details**: Money is `decimal.Decimal` + ISO-4217 currency, constructed from string (never float64); reject non-positive amounts in `NewCap`. UUIDv4 IDs. Domain package has zero external imports beyond `decimal` and `uuid`.

**Acceptance**: Constructors reject amount ≤ 0 and unknown period; `NewMoney` rejects float64 input and mixed-currency comparison; table-driven tests cover 100% of constructor/enum branches.

### 2. Spend window accumulation (MO-4, MO-5)
**Prompt**: "Implement per-cap spend windows that compute current-window boundaries in the account timezone and accumulate approved spend idempotently by transaction_id."

**File**: `internal/domain/spend_window.go`, `internal/service/accumulator.go`

**Function/class**: `SpendWindow`, `Accumulator.Record(tx)`

**Details**: Boundaries per IN-4; idempotency per NFR-REL-2; correct rollover incl. DST.

**Acceptance**: Replaying a transaction_id does not double-count; boundary tests pass for daily/weekly/monthly across a DST transition.

### 3. Cap CRUD service (MO-1, MO-6)
**Prompt**: "Implement create/read/update/delete for caps with deny-by-default per-resource authorization and idempotency-key support."

**File**: `internal/ports/in/cap_service.go` (use-case interface), `internal/service/cap_service.go` (impl), `internal/ports/out/cap_repository.go` (repository interface), `internal/adapters/out/memory/cap_repository.go` (repository impl), `internal/adapters/in/http/cap_handler.go` (Fiber HTTP handlers)

**Function/class**: `CapService` interface — `Create`, `Get`, `List`, `Update`, `Delete`; `CapRepository` interface + in-memory impl; `CapHandler` HTTP handlers delegating to `CapService` (no business logic)

**Details**: IN-3 idempotency; IN-9 validation; users edit only own caps; delete forbidden on fraud holds (EC-5).

**Acceptance**: Unauthorized actor → `403`; invalid cap → `422`; repeated idempotency key returns original; deleting a fraud hold → `403` and audited attempt.

### 4. Precedence resolver (MO-2, MO-3)
**Prompt**: "Given all caps applicable to a card+category+period, return the binding (most restrictive) cap, with fraud holds overriding everything."

**File**: `internal/service/precedence.go`

**Function/class**: `ResolveBindingCap(cardID, ctx)`

**Details**: IN-5 most-restrictive-wins; fraud hard-cap overrides.

**Acceptance**: Fixture matrix proves the correct binding cap is chosen; fraud hold always wins.

### 5. Enforcement decision engine (MO-3, MO-4)
**Prompt**: "Implement evaluate(transaction) returning ALLOW/DENY with the binding cap and remaining amount, applying the fail-safe matrix on datastore errors."

**File**: `internal/ports/in/enforcement.go` (use-case interface), `internal/service/enforcement.go` (impl), `internal/ports/out/spend_store.go` (spend-store interface), `internal/adapters/out/memory/spend_store.go` (impl + fail-safe cache), `internal/adapters/in/http/enforcement_handler.go` (Fiber handler for `evaluate`)

**Function/class**: `Enforcement.Evaluate(ctx, tx) (Decision, error)`; `SpendStore` interface + in-memory impl; `EnforcementHandler` HTTP handler (delegates to `Enforcement`, no business logic)

**Details**: `≤ remaining` allows, `>` denies (EC-2); concurrency-safe per card window (EC-3); fail-safe (IN-7, EC-8) sets `degraded: true`.

**Acceptance**: Decision matrix tests pass; concurrent race test never double-spends; injected DB outage yields fail-safe decision with `degraded` flag and an audit record.

### 6. Append-only audit writer (MO-5)

**Prompt**: "Implement an append-only audit writer enrolled in the caller's transaction, recording all NFR-AUD-2 fields; no update/delete methods."

**File**: `internal/ports/out/audit.go` (writer interface), `internal/service/audit.go` (impl), `internal/adapters/out/memory/audit_store.go` (append-only store impl), `internal/adapters/in/http/audit_handler.go` (Fiber handler for audit read), `docs/audit-fields.md`

**Function/class**: `AuditLog.Record(ctx, event)`; append-only `AuditStore`; `AuditHandler` HTTP handler (read, restricted to fraud/compliance)

**Details**: NFR-AUD-1/2/3; reason codes incl. `FRAUD_HOLD_LOCKED`, `CAP_EXCEEDED`, `DEGRADED_DECISION`.

**Acceptance**: A failed audit write rolls back the state change; schema/code expose no mutation path; `docs/audit-fields.md` lists every field and matches the writer.

### 7. Masking utility (NFR-SEC-2, MO-6)

**Prompt**: "Create the single masking helper for card references, emails, and money where required, used by all human-facing/support responses."

**File**: `pkg/masking/masking.go`

**Function/class**: `MaskCard`, `MaskEmail`

**Details**: IN-8 — the only allowed path to render an identifier to a human; card → `•••• 4321`, email → `u***@gmail.com`.

**Acceptance**: Support-facing responses contain only masked identifiers; unit tests cover short/edge inputs; a lint/grep check finds no raw-identifier rendering outside this helper.

### 8. Fraud hold workflow (MO-7)

**Prompt**: "Implement placing and releasing a fraud hold (incl. zero-cap freeze) that users cannot override; only fraud/compliance roles can release."

**File**: `internal/service/fraud_hold.go` (impl), `internal/adapters/in/http/fraud_handler.go` (Fiber handlers for place/release)

**Function/class**: `FraudHoldService.Place`, `FraudHoldService.Release`; `FraudHandler` HTTP handlers (delegate to `FraudHoldService`, role-restricted)

**Details**: EC-5, EC-11; release restricted to fraud/compliance; both place and release audited.

**Acceptance**: User override attempt → `403` + audited; freeze applies to all new evaluations immediately; only authorized roles release.

### 9. REST API wiring & cross-cutting (MO-1, MO-3, MO-6)

**Prompt**: "Register all routes and wire the per-entity handlers (cap, enforcement, fraud, audit) into the Fiber app with role-scoped middleware, the global error mapper, and pagination — no business logic here."

**File**: `internal/adapters/in/http/router.go` (route registration), `cmd/api/main.go` (DI wiring), `pkg/errorhandler/errorhandler.go` (typed errors → HTTP codes), `pkg/middleware/*.go` (RequestID, logger, auth/role)

**Function/class**: `RegisterRoutes(app, handlers...)` + global error handler + middleware (the entity handlers themselves are built in Tasks 3, 5, 6, 8)

**Details**: IN-6 error mapping; pagination NFR-PERF-3; role-scoped access (audit read restricted to fraud/compliance); DENY is `200` not an error.

**Acceptance**: Endpoint tests assert status codes per IN-6; pagination caps at 100; audit endpoint `403` for users/support; every entity handler is reachable through a registered route.

### 10. Test suite & fixtures (Verification)

**Prompt**: "Build unit, integration, decision-matrix, concurrency, and failure-injection tests with shared fixtures covering every edge case EC-1..EC-12."

**File**: `internal/**/*_test.go`, `tests/`, `tests/fixtures/`

**Function/class**: table-driven `Test*` functions + shared fixtures

**Details**: cover all edge cases and the verification matrix; include a degraded-mode and DST-rollover fixture; use `httptest`/Fiber test app for black-box integration.

**Acceptance**: Every EC-row has a named test; `go test ./...` runs all categories in CI; coverage ≥ 90% on `internal/service/` and `internal/domain/`.

### 11. Performance check (NFR-PERF)

**Prompt**: "Add a benchmark/load test asserting evaluate p95 ≤ 50 ms and ≥ 500 req/s per instance against seeded data."

**File**: `tests/perf/evaluate_bench_test.go`

**Function/class**: `BenchmarkEvaluate`, `TestEvaluateLatencyP95`

**Details**: NFR-PERF-1/5; use `go test -bench` + a latency-percentile harness; report p95/p99; run against a representative cap fixture set.

**Acceptance**: Benchmark records p95/p99 and throughput; CI fails if p95 > 50 ms or throughput < 500 req/s.