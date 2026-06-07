# 🏦 Homework 3: Specification-Driven Design — Spending Cap Management

> **Student Name**: Roman Reznik
> **Date Submitted**: 2026-06-07
> **AI Tools Used**: Claude Code (Opus 4.8) - tasks; CoPilot Chat (GPT-5.4) - verification;

---

## 📋 Project Overview

A complete **specification package** (no implementation) for a **Spending Cap Management** feature in a regulated banking context. Users set, view, change, and remove spending caps (by period and category) that are enforced in real time at authorization, with every change auditable and visible to support, fraud, and compliance.

**Stakeholders**: End-user, Support agent, Fraud team, Ops/Compliance.
**Regulatory level**: Basic — audit log + PII/PAN masking.

### Deliverables in this folder

| File | Purpose |
|------|---------|
| [`spec_spendingCapManagement.md`](./spec_spendingCapManagement.md) | Layered specification: high-level objective → mid-level objectives → non-functional/policy → implementation notes → context → edge cases → verification → low-level tasks |
| [`CLAUDE.md`](./CLAUDE.md) | AI coding-partner configuration for Claude Code: stack, domain rules, testing/verification expectations, edge-case handling |
| [`.claude/rules/rules.md`](./.claude/rules/rules.md) | General regulated-fintech AI rules (reusable across features) + PR review gates |
| [`.claude/rules/spending-cap.md`](./.claude/rules/spending-cap.md) | Spending-Cap-specific rules (enforcement, cap precedence, fraud holds, fail-safe, windows) + feature review gates |
| [`.claude/rules/code_style.md`](./.claude/rules/code_style.md) | Code style & conventions (structure, error→HTTP mapping, doc comments) |
| [`.claude/rules/testing_rules.md`](./.claude/rules/testing_rules.md) | Testing & verification rules (test categories, coverage, verification mechanisms) |
| `README.md` | This file — rationale and industry best practices |

> All rule files are auto-loaded by Claude Code at session start — it discovers every `.md` under `.claude/rules/` (same priority as `CLAUDE.md`).

---

## 🎯 Rationale

**Why this structure.** The spec is layered so requirements are traceable top to bottom: each Mid-Level Objective (MO-1…MO-7) is observable, each Low-Level Task names the objective it serves and ends with acceptance criteria, and the Verification table maps every objective back to how we'd prove it. An engineer or AI agent can execute it without guessing.

**Why these performance targets.** Spending-cap enforcement sits *inline* in the card-authorization path, which usually has a sub-second end-to-end budget shared across many checks. So `evaluate` gets a small slice: **p95 ≤ 50 ms, p99 ≤ 120 ms** (NFR-PERF-1). CRUD operations are user-interactive but not latency-critical, hence a looser **p95 ≤ 300 ms** (NFR-PERF-2). Pagination is capped at 100 (NFR-PERF-3) to bound query cost, and read-after-write ≤ 1 s (NFR-PERF-4) matches the expectation that a user who saves a cap and refreshes sees it. Throughput of ≥ 500 req/s per instance (NFR-PERF-5) reflects a high-volume authorization stream. All numbers are labeled **assumed targets** and justified rather than asserted as "fast."

**Why this verification depth.** Because the feature blocks people's money, the costly failure modes are *double-spend*, *leaked PII*, and *missing audit*. The spec therefore makes concurrency (EC-3), idempotency (EC-4, NFR-REL-2), fail-safe degradation (EC-8), and audit atomicity (NFR-AUD-3) first-class — each has a dedicated edge-case row, a verification entry, and an acceptance criterion in the task list, not a vague mention.

**Why the fraud-hold mechanic.** Realistic FinTech requires that an investigator can freeze a card and the cardholder cannot undo it. MO-7 + EC-5/EC-11 + Task 8 encode that as an override-everything hard cap with role-restricted release, which also exercises the authorization model meaningfully.

**Why this architecture & decomposition.** The package assumes a Go hexagonal (ports & adapters) layout, defined once in `CLAUDE.md` (Architecture) so the spec never repeats it. Each feature task (3 Cap CRUD, 5 Enforcement, 6 Audit, 8 Fraud hold) is decomposed through all layers it touches — use-case interface → service → outbound port → in-memory adapter → per-entity HTTP handler — so a task is independently buildable and testable. Task 9 only wires routes, middleware, and the error mapper; it owns no handler files, so tasks never collide on the same file. Domain stays import-free, services depend on ports (never adapters), and concrete adapters are bound only in `cmd/api`.

**Why rules are modular.** Project rules live under `.claude/rules/` split by concern — `rules.md` (reusable regulated-fintech defaults), `spending-cap.md` (feature behavior), `code_style.md`, `testing_rules.md` — and are auto-loaded by Claude Code each session. `CLAUDE.md` holds only non-rule context (mission, stack, architecture, DoD) and points to the rule files, so nothing is duplicated across files.

---

## 🏛️ Industry Best Practices — and where they appear

| Best practice | Where in the package |
|---------------|----------------------|
| **No PAN in storage/logs**; tokenized card references | spec NFR-SEC-1, IN-8; rules.md Domain Rule 1 |
| **PII/identifier masking** through a single audited path | spec NFR-SEC-2 + Task 7; rules.md Domain Rule 2 + Review Gate 2 |
| **Decimal money, explicit currency**, no float | spec IN-1, NFR-SEC-4; rules.md Domain Rule 3 + Review Gate 1 |
| **Append-only, transactional audit trail** with defined fields | spec NFR-AUD-1..3 + Task 6 + `docs/audit-fields.md`; rules.md Domain Rule 4; spending-cap.md "Audit" |
| **Deny-by-default, per-resource authorization** | spec NFR-SEC-3, MO-6; rules.md Domain Rule 5 |
| **Idempotency** (keys + per-transaction de-dup) | spec IN-3, NFR-REL-2, EC-4; Task 2; rules.md Domain Rule 6; spending-cap.md "Enforcement" |
| **Fail-safe / graceful degradation** in the inline decision path | spec NFR-REL-1, IN-7, EC-8; Task 5; spending-cap.md "Enforcement" |
| **Deny is a valid business answer, not an error** (clean error semantics) | spec IN-6; spending-cap.md "Enforcement"; code_style.md "Errors" |
| **Concurrency safety / no double-spend** | spec EC-3, Task 5; verification matrix MO-3 |
| **Timezone-correct windows incl. DST** | spec IN-4, EC-7, Task 2; spending-cap.md "Enforcement" |
| **Opaque UUID IDs, no leaked PKs** | spec IN-2; rules.md Domain Rule 7 |
| **Compliance review checkpoints + field documentation** | spec Verification "Review checkpoints" + `docs/audit-fields.md` |
| **Measurable SLOs with stated rationale** | spec NFR-PERF-1..5; this README "Rationale" |
| **Traceability goals→tasks** (every task tags its objective) | spec Low-Level Tasks; CLAUDE.md Definition of Done |
| **Hexagonal architecture** (ports & adapters), layered task decomposition | CLAUDE.md Architecture; spec Tasks 3/5/6/8 (domain→port→adapter→handler), Task 9 (wiring) |
| **Modular, auto-loaded AI rules** with no cross-file duplication | `.claude/rules/*.md`; CLAUDE.md "Domain Rules" pointer |

---

<div align="center">

*This project was completed as part of the AI-Assisted Development course. No coding required — the specification is the graded artifact.*

</div>