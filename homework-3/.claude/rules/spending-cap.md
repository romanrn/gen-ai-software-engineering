# spending-cap.md — Spending Cap Management rules

> Rules specific to the Spending Cap Management feature. General regulated-fintech defaults are in [`rules.md`](./rules.md). For the full rationale see [spec_spendingCapManagement.md](../../spec_spendingCapManagement.md).

## Enforcement & caps

- **`evaluate` decision is a value, not an error.** `spend ≤ remaining` → `ALLOW`; `>` → `DENY`. Both return HTTP `200` with a `{decision: ...}` body — never a `4xx`/`5xx`. (spec IN-6, EC-2)
- **Most-restrictive applicable cap wins**, and a fraud hold overrides all user caps. (spec IN-5)
- **Fraud holds are locked for users.** End-users and support can never modify or remove a fraud hold; only fraud/compliance release it. Log override attempts with reason `FRAUD_HOLD_LOCKED`. (spec MO-7, EC-5, EC-11)
- **Fail safe in `evaluate`.** Degrade deterministically per the fail-safe matrix; fraud freezes must still apply during a primary-DB outage; set `degraded: true` and audit the degraded decision. (spec IN-7, NFR-REL-1, EC-8)
- **Spend accumulation is idempotent per `transaction_id`** — replays return the prior result; never double-count a replayed approved transaction. (spec NFR-REL-2, EC-4)
- **Serialize evaluation + accumulation per card window.** Concurrent transactions racing the last remaining amount must not double-spend — one wins, the other gets `DENY`. (spec EC-3)
- **A card with no applicable caps** → `evaluate` returns `ALLOW` by absence, with no audit record. (spec EC-12)
- **Windows are computed in the account timezone, persisted in UTC**; handle period rollover including DST. (spec IN-4, EC-7)

## Audit (feature)

- Every cap mutation **and every `DENY`** writes one append-only audit record in the same transaction, with a reason code: `CAP_EXCEEDED`, `FRAUD_HOLD_LOCKED`, or `DEGRADED_DECISION`. (spec NFR-AUD-2)

## Naming & limits

- Service types: `Enforcement`, `CapService`, `AuditLog`, `FraudHoldService` in `internal/service/`; their interfaces in `internal/ports/{in,out}`.
- Pagination: default page size 25, max 100. (spec NFR-PERF-3)

## Review Gates (feature — PR rejected if violated)

1. Any edge case EC-1..EC-12 touched without a corresponding test → reject.
2. An `evaluate` change without a perf check against NFR-PERF-1 (p95 ≤ 50 ms) → reject.
3. A cap `DENY` (including a degraded decision) that emits no audit record → reject.