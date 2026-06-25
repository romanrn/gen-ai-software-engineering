# CLAUDE.md — AI Coding Partner Guidelines

> This file configures how Claude Code (and any AI coding agent) must behave when implementing the [Spending Cap Management spec](./spec_spendingCapManagement.md). Read it before generating any code.
>
> Enforceable rules and PR review gates live under `.claude/rules/` — `rules.md` (general regulated-fintech defaults), `spending-cap.md` (this feature's rules), `code_style.md` (code style & conventions), and `testing_rules.md` (testing & verification). Claude Code auto-loads every `.md` there at session start (same priority as this file), so no import is needed.

## Mission

You are implementing a **regulated FinTech** feature. Correctness, auditability, and not leaking sensitive data outrank cleverness, brevity, and speed of delivery. When in doubt, choose the safer, more explicit, more auditable option and state the assumption.

## Tech Stack Assumptions

- **Language**: Go (1.22+).
- **Money**: `shopspring/decimal` only — never `float64` for monetary amounts; construct from string.
- **API**: REST via **Fiber**; handlers hold no business logic.
- **Persistence**: a repository behind an outbound port (in-memory `sync.RWMutex` map for now) + an append-only audit store; a fast cap cache for fail-safe enforcement reads.
- **Testing**: standard `go test` with table-driven tests; black-box integration via `httptest`/Fiber test app; benchmarks via `go test -bench`.
- **IDs**: UUIDv4, opaque (`google/uuid`).

If a dependency is unstated, pick the conventional choice for this stack and note it — do not invent exotic libraries.

## Architecture

Hexagonal architecture — domain at the center, everything else plugged in through interfaces.

```
Domain Layer  →  Ports (interfaces)  →  Adapters (implementations)
```

```
src/
├── cmd/api/                   # Entry point — DI wiring, Fiber server setup
├── internal/
│   ├── domain/                # Cap, Money, SpendWindow, enums, typed errors (zero external imports beyond decimal/uuid)
│   ├── ports/
│   │   ├── in/                # CapService, Enforcement interfaces (use cases)
│   │   └── out/               # CapRepository, AuditLog, SpendStore interfaces (storage)
│   ├── service/               #  Precedence · Accumulator · Enforcement · CapService impl · FraudHoldService
│   └── adapters/
│       ├── in/http/           # Fiber handlers per entity: cap · enforcement · fraud · audit (no business logic)
│       └── out/memory/        # in-memory cap repository · spend store · append-only audit store (+ fail-safe cap cache)
├── pkg/
│   ├── masking/               # the single masking helper for identifiers/money
│   ├── errorhandler/          # global Fiber error mapper (typed errors → HTTP codes)
│   └── middleware/            # RequestID, logger, auth/role extraction
└── tests/                     # black-box integration + performance tests
    └── fixtures/              # sample data for tests
```

**Rules of the layout**: domain imports nothing inward-facing; services depend on ports, never on adapters; adapters depend on ports; HTTP handlers translate transport↔domain and never embed business logic. Wire concrete adapters to ports only in `cmd/api/`.

## Domain Rules (banking)

The rules are auto-loaded every session and not duplicated here — treat them as binding:
- General regulated-fintech defaults (PAN/PII, money, audit, authorization, idempotency, UUIDs) → `.claude/rules/rules.md`.
- Spending-Cap-specific rules (enforcement decisions, cap precedence, fraud holds, fail-safe, spend windows) → `.claude/rules/spending-cap.md`.
- Code style & conventions (structure, error→HTTP mapping, doc comments) → `.claude/rules/code_style.md`.
- Testing & verification (test categories, coverage, verification mechanisms) → `.claude/rules/testing_rules.md`.

> Prohibitions ("never") are not repeated here — see the files under `.claude/rules/` (Domain Rules + Never + Review Gates). One extra working norm: if the spec is ambiguous, do **not** silently widen scope — state the assumption in the PR description.

## Definition of Done for any task

A task is done only when: code matches the named spec objective/NFR, its acceptance criteria pass, edge-case tests are green, no PII/PAN appears in logs, audit records are emitted where required, and a one-line traceability note (`implements MO-x / NFR-y`) is in the PR.
