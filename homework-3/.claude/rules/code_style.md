# code_style.md — Code Style & Conventions

> General code-style and structural conventions for this Go project. Domain/feature rules live in [`rules.md`](./rules.md) and [`spending-cap.md`](./spending-cap.md). For the full rationale see [CLAUDE.md](../../CLAUDE.md) and [spec_spendingCapManagement.md](../../spec_spendingCapManagement.md).

## Structure

- Small, single-responsibility functions; keep pure domain logic separated from I/O.
- Isolate persistence and external systems behind interfaces (ports) so use cases are unit-testable.
- No magic numbers — limits, page sizes, and latency budgets come from named constants/config.

## Errors

- Typed domain errors mapped to stable HTTP codes: validation `422`, auth `401/403`, not-found `404`, idempotency replay/conflict `409`.

## Documentation

- Go doc comments on exported types and methods, stating which spec objective/NFR they satisfy (`// implements MO-3 / NFR-PERF-1`).