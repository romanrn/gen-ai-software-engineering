# rules.md — General FinTech AI Rules

> Reusable, regulated-fintech defaults that apply to **any** feature in this project. Rules specific to the Spending Cap Management feature live in [`spending-cap.md`](./spending-cap.md). For the full rationale see [spec_spendingCapManagement.md](../../spec_spendingCapManagement.md) and [CLAUDE.md](../../CLAUDE.md).

## Domain Rules (banking) — non-negotiable

1. **Never store or log a PAN.** Reference cards by `card_id` / token only. (spec NFR-SEC-1)
2. **Mask before display.** Any identifier shown to a human goes through `pkg/masking`. There is no other allowed path. (spec IN-8, NFR-SEC-2)
3. **Money is `decimal.Decimal` + currency.** No `float64` math, no mixed-currency comparison; construct from string. (spec IN-1, NFR-SEC-4)
4. **Audit is append-only and transactional.** Every state mutation writes one audit record in the same transaction as the change. No update/delete on audit. (spec NFR-AUD-1..3)
5. **Deny by default.** Authorization is enforced server-side, per resource. Absence of an explicit allow = `403`. (spec NFR-SEC-3)
6. **Idempotent writes.** Honor `Idempotency-Key`; repeated keys return the original result without re-applying. (spec IN-3)
7. **Expose only opaque UUIDv4 IDs** — never database primary keys or sequential integers. (spec IN-2)

## Never

- Never widen authorization scope to make a test pass.

## Review Gates (PR will be rejected if violated)

1. Any monetary `float64` → reject.
2. Any raw identifier rendered outside `pkg/masking` → reject.
3. Missing audit record on a state mutation → reject.
4. New audit mutation/delete path → reject.