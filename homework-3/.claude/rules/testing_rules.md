# testing_rules.md — Testing & Verification Rules

> How code in this project must be tested and verified. The enforcement gates for these live in [`rules.md`](./rules.md) and [`spending-cap.md`](./spending-cap.md); general fintech rules in `rules.md`. For the rationale see [spec_spendingCapManagement.md](../../spec_spendingCapManagement.md) (Verification) and [CLAUDE.md](../../CLAUDE.md).

## Expectations

- Every edge case enumerated in the spec (**EC-1 … EC-12**) has a named test.
- Required test categories: unit (domain, precedence, windows, masking) · integration (CRUD + audit atomicity) · decision matrix (`evaluate`) · concurrency (EC-3 race) · failure injection (EC-8 degraded mode) · performance (NFR-PERF-1/5).
- Tests are table-driven; black-box integration via `httptest` / Fiber test app.
- Coverage target ≥ 90% on `internal/service/` and `internal/domain/`.

## Verification mechanisms

- A test proves a failed audit write **rolls back** the state change (verifies the transactional-audit rule in `rules.md`).
- A grep/lint check proves no raw identifier is rendered outside `pkg/masking` (verifies the masking rule in `rules.md`).