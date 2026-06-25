# Homework-4 Specification: 4-Agent Pipeline

## Application: Go UTC Time Service

A minimal Go HTTP service with seeded bugs and a security vulnerability for the pipeline to process.

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/time` | Returns current UTC time |
| GET | `/health` | Health check |

### Seeded Issues

**bug#1 — Wrong timezone** — returns local server time instead of UTC.
**bug#2 — Integer overflow in uptime counter** — uptime counter overflows after ~2 minutes, producing negative values in `/health`.
**sec#1 — Hardcoded API key** — API key hardcoded in source, exposed in version control.

---

## Pipeline Model

The pipeline is **generic** — accepts any bug directory and runs all 6 agents end-to-end.
The only manual input is `bug-context.md` (seed file describing the observable symptom).

```
./run-pipeline.sh <bug-dir>

  Bug Researcher        → <bug-dir>/research/codebase-research.md
        ↓
  Bug Research Verifier → <bug-dir>/research/verified-research.md
        ↓
  Bug Planner           → <bug-dir>/implementation-plan.md
        ↓
  Bug Fixer             → <bug-dir>/fix-summary.md
        ↓
  Security Verifier     → <bug-dir>/security-report.md
        ↓
  Unit Test Generator   → <bug-dir>/test-report.md
```

---

## Project Structure

```
homework-4/
├── SPEC.md
├── PLAN.md
├── README.md
├── HOWTORUN.md
├── docker-compose.yml
├── run-pipeline.sh                          ← ./run-pipeline.sh <bug-dir>
├── src/
│   ├── main.go
│   ├── go.mod
│   ├── main_test.go
│   └── Dockerfile
├── agents/
│   ├── bug-researcher.agent.md
│   ├── research-verifier.agent.md
│   ├── bug-planner.agent.md
│   ├── bug-fixer.agent.md
│   ├── security-verifier.agent.md
│   └── unit-test-generator.agent.md
├── skills/
│   ├── codebase-research-format.md          ← used by bug-researcher
│   ├── research-quality-measurement.md      ← used by research-verifier
│   ├── implementation-plan-format.md        ← used by bug-planner
│   └── unit-tests-FIRST.md                  ← used by unit-test-generator
├── context/bugs/
│   ├── bug001/
│   │   ├── bug-context.md                   ← manually authored seed
│   │   ├── research/
│   │   │   ├── codebase-research.md         ← agent: bug-researcher
│   │   │   └── verified-research.md         ← agent: research-verifier
│   │   ├── implementation-plan.md           ← agent: bug-planner
│   │   ├── fix-summary.md                   ← agent: bug-fixer
│   │   ├── security-report.md               ← agent: security-verifier
│   │   └── test-report.md                   ← agent: unit-test-generator
│   ├── bug002/
│   │   └── (same structure)
│   └── sec001/
│       └── (same structure)
└── docs/screenshots/
```

---

## Action Items (Ordered)

### Phase 0 — Scaffold ✅

> All Phase 0 files are **manually authored** — the starting state the pipeline operates on.
> `bug-context.md` describes only the observable symptom — no code, no location, no fix hints.

| # | Action | Status |
|---|--------|--------|
| 0.1 | `src/main.go` — Go HTTP service with 3 seeded issues | done |
| 0.2 | `src/go.mod` | done |
| 0.3 | `src/Dockerfile` — multi-stage build | done |
| 0.4 | `docker-compose.yml` — exposes :8080, passes `API_KEY` env | done |
| 0.5 | `src/main_test.go` — stub tests | done |
| 0.6 | `context/bugs/bug001/bug-context.md` | done |
| 0.7 | `context/bugs/bug002/bug-context.md` | done |
| 0.8 | `context/bugs/sec001/bug-context.md` | done |

### Phase 1 — Skills ✅

| # | Action | Status |
|---|--------|--------|
| 1.1 | `skills/codebase-research-format.md` | done |
| 1.2 | `skills/research-quality-measurement.md` | done |
| 1.3 | `skills/implementation-plan-format.md` | done |
| 1.4 | `skills/unit-tests-FIRST.md` | done |

### Phase 2 — Agents

| # | Action | Output |
|---|--------|--------|
| 2.1 | `agents/bug-researcher.agent.md` — model: claude-opus-4-8; reads bug-context + src; uses codebase-research-format skill; writes `research/codebase-research.md` | done |
| 2.2 | `agents/research-verifier.agent.md` — model: claude-opus-4-8; reads codebase-research; uses research-quality-measurement skill; writes `research/verified-research.md` | — |
| 2.3 | `agents/bug-planner.agent.md` — model: claude-sonnet-4-6; reads verified-research; uses implementation-plan-format skill; writes `implementation-plan.md` | done |
| 2.4 | `agents/bug-fixer.agent.md` — model: claude-haiku-4-5-20251001; reads implementation-plan; applies fixes; runs tests via Docker; writes `fix-summary.md` | — |
| 2.5 | `agents/security-verifier.agent.md` — model: claude-opus-4-8; reads fix-summary + src; rates CRITICAL/HIGH/MEDIUM/LOW/INFO; writes `security-report.md` | — |
| 2.6 | `agents/unit-test-generator.agent.md` — model: claude-sonnet-4-6; reads fix-summary + src; uses FIRST skill; writes `test-report.md` | — |

### Phase 3 — Pipeline Runner

| # | Action | Output |
|---|--------|--------|
| 3.1 | Create `run-pipeline.sh <bug-dir>` — runs all 6 agents in order | shell script |
| 3.2 | `chmod +x run-pipeline.sh` | — |
| 3.3 | Test: `./run-pipeline.sh context/bugs/bug001` end-to-end | — |

### Phase 4 — Documentation & Screenshots

| # | Action | Output |
|---|--------|--------|
| 4.1 | `README.md` — author info, overview, pipeline usage, app run command | README |
| 4.2 | `HOWTORUN.md` — prerequisites (Docker only), run app, run pipeline | HOWTORUN |
| 4.3 | Screenshots for each of the 3 issues | `docs/screenshots/` |

### Phase 5 — Verification & PR

| # | Action | Output |
|---|--------|--------|
| 5.1 | `./run-pipeline.sh` runs end-to-end for all 3 issues | — |
| 5.2 | Verify all 6 artifact files exist in each bug dir | — |
| 5.3 | `docker compose down` — stop the running container | — |
| 5.4 | `docker compose run --rm app go test ./...` — all tests green (against fixed code) | — |
| 5.5 | `docker compose up --build` — rebuild image with fixed code, start server | — |
| 5.6 | `curl -H "X-Api-Key: $API_KEY" http://localhost:8080/time` — returns correct UTC after fixes | — |
| 5.7 | Open PR with summary, screenshots, author info | PR |

---

## Agent Model Justification

| Agent | Model | Reason |
|-------|-------|--------|
| bug-researcher | claude-opus-4-8 | Root cause analysis requires strong reasoning over code |
| research-verifier | claude-opus-4-8 | Fact-checking file:line references requires precision |
| bug-planner | claude-sonnet-4-6 | Planning from verified facts — balanced capability |
| bug-fixer | claude-haiku-4-5-20251001 | Routine mechanical edits; speed and cost matter |
| security-verifier | claude-opus-4-8 | Security analysis requires nuanced judgment |
| unit-test-generator | claude-sonnet-4-6 | Code understanding without max reasoning needed |

---

## Bug Details (Spec Reference — not for agents)

### bug#1 — Wrong UTC time

```go
// BEFORE
now := time.Now()
// AFTER
now := time.Now().UTC()
```

### bug#2 — Integer overflow in uptime counter

```go
// BEFORE
var uptimeSeconds int8 = 0
// AFTER
var uptimeSeconds int64 = 0
```

### sec#1 — Hardcoded API key

```go
// BEFORE
const apiKey = "supersecret-hardcoded-key-12345"
func validate(r *http.Request) bool {
    return r.Header.Get("X-Api-Key") == apiKey
}
// AFTER
func validate(r *http.Request) bool {
    key := os.Getenv("API_KEY")
    return key != "" && r.Header.Get("X-Api-Key") == key
}
```