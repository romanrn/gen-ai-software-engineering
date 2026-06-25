# Plan: homework-4 — 4-Agent Pipeline Implementation

## Context

homework-4 requires a generic 6-agent pipeline (4 required + 2 additional) that runs against any bug directory and produces all artifacts inside it. The only manual input per issue is `bug-context.md`. The 3 issues (bug001, bug002, sec001) demonstrate the pipeline works for different problem types.

---

## Pipeline

```
./run-pipeline.sh <bug-dir>

  bug-researcher        reads: bug-context.md, src/
                        skill: codebase-research-format
                        writes: research/codebase-research.md

  research-verifier     reads: research/codebase-research.md, src/
                        skill: research-quality-measurement
                        writes: research/verified-research.md

  bug-planner           reads: research/verified-research.md, src/
                        skill: implementation-plan-format
                        writes: implementation-plan.md

  bug-fixer             reads: implementation-plan.md
                        runs:  docker compose run --rm app go test ./...
                        writes: fix-summary.md

  security-verifier     reads: fix-summary.md, src/
                        writes: security-report.md

  unit-test-generator   reads: fix-summary.md, src/
                        skill: unit-tests-FIRST
                        writes: test-report.md
```

---

## Phase 0 — Go Application ✅

- `src/main.go` — 3 seeded issues (bug#1: wrong timezone, bug#2: int8 overflow, sec#1: hardcoded key)
- `src/go.mod`, `src/Dockerfile`, `src/main_test.go`
- `docker-compose.yml` — no local Go needed
- `context/bugs/bug001/bug-context.md` — observable symptom only
- `context/bugs/bug002/bug-context.md` — observable symptom only
- `context/bugs/sec001/bug-context.md` — observable symptom only

---

## Phase 1 — Skills ✅

- `skills/codebase-research-format.md` — structure for research output (bug-researcher)
- `skills/research-quality-measurement.md` — GOLD/SILVER/BRONZE/FAIL levels (research-verifier)
- `skills/implementation-plan-format.md` — structure for implementation plan (bug-planner)
- `skills/unit-tests-FIRST.md` — FIRST checklist (unit-test-generator)

---

## Phase 2 — Agents

Agent frontmatter format:
```yaml
---
name: <agent-name>
description: <one-line role>
model: <model-id>
tools:
  - Read
  - Edit
  - Bash
---
```

| Agent | File | Model | Status |
|-------|------|-------|--------|
| bug-researcher | `agents/bug-researcher.agent.md` | claude-opus-4-8 | done |
| research-verifier | `agents/research-verifier.agent.md` | claude-opus-4-8 | — |
| bug-planner | `agents/bug-planner.agent.md` | claude-sonnet-4-6 | done |
| bug-fixer | `agents/bug-fixer.agent.md` | claude-haiku-4-5-20251001 | — |
| security-verifier | `agents/security-verifier.agent.md` | claude-opus-4-8 | — |
| unit-test-generator | `agents/unit-test-generator.agent.md` | claude-sonnet-4-6 | — |

---

## Phase 3 — Pipeline Runner

**`run-pipeline.sh`**
```bash
#!/bin/bash
set -e
BUG_DIR=${1:?Usage: ./run-pipeline.sh <bug-dir>}

claude --agent agents/bug-researcher.agent.md --context "$BUG_DIR"
claude --agent agents/research-verifier.agent.md --context "$BUG_DIR"
claude --agent agents/bug-planner.agent.md --context "$BUG_DIR"
claude --agent agents/bug-fixer.agent.md --context "$BUG_DIR"
claude --agent agents/security-verifier.agent.md --context "$BUG_DIR"
claude --agent agents/unit-test-generator.agent.md --context "$BUG_DIR"
```

Usage:
```bash
./run-pipeline.sh context/bugs/bug001
./run-pipeline.sh context/bugs/bug002
./run-pipeline.sh context/bugs/sec001
```

---

## Phase 4 — Documentation

- `README.md` — author (Roman Reznik), overview, pipeline usage, app run command
- `HOWTORUN.md` — prerequisites (Docker only), run app, run pipeline per issue

---

## Execution Order

```
Phase 2:
  2.1  agents/research-verifier.agent.md
  2.2  agents/bug-fixer.agent.md
  2.3  agents/security-verifier.agent.md
  2.4  agents/unit-test-generator.agent.md

Phase 3:
  3.1  run-pipeline.sh

Phase 4:
  4.1  README.md
  4.2  HOWTORUN.md
```

---

## Verification

1. `API_KEY=supersecret-hardcoded-key-12345 docker compose up --build` — server starts on :8080
2. `curl -H "X-Api-Key: supersecret-hardcoded-key-12345" http://localhost:8080/time` — returns wrong timezone (bug#1 visible)
3. `./run-pipeline.sh context/bugs/bug001` — all 6 agents run, all artifacts created
4. `./run-pipeline.sh context/bugs/bug002` — same
5. `./run-pipeline.sh context/bugs/sec001` — same
6. `docker compose down` — stop the running container
7. `docker compose run --rm app go test ./...` — all tests green (against fixed code)
8. `API_KEY=your-secret docker compose up --build` — rebuild image with fixed code, start server
9. `curl -H "X-Api-Key: $API_KEY" http://localhost:8080/time` — returns correct UTC after fixes