# Plan: homework-4 — 4-Agent Pipeline Implementation

## Context

homework-4 requires building a generic 4-agent pipeline that runs against any bug directory and produces all artifacts inside it. The 3 issues (bug001, bug002, sec001) are demonstration examples. Currently Phase 0 is complete (src/, docker-compose.yml, context/bugs/ seed files).

---

## Pipeline Model

```
./run-pipeline.sh <bug-dir>

  Bug Researcher (manual)
        ↓  reads: <bug-dir>/research/codebase-research.md
  Research Verifier
        ↓  writes: <bug-dir>/research/verified-research.md
  Bug Planner (manual)
        ↓  reads: <bug-dir>/implementation-plan.md
  Bug Fixer
        ↓  writes: <bug-dir>/fix-summary.md
  Security Verifier
        ↓  writes: <bug-dir>/security-report.md
  Unit Test Generator
        ↓  writes: <bug-dir>/test-report.md
```

---

## Phase 0 — Go Application ✅

All files created and verified (`go build`, `go test` pass).

> All Phase 0 files are **manually authored seed files** — the starting state the pipeline operates on, not agent outputs. `context/bugs/` files document the observable wrong behavior so agents have input context to work from.

- `src/main.go` — 3 seeded issues (bug#1, bug#2, sec#1)
- `src/go.mod`, `src/Dockerfile`
- `docker-compose.yml` — no local Go needed
- `src/main_test.go` — stub tests
- `context/bugs/bug001/bug-context.md` — observable symptom only
- `context/bugs/bug002/bug-context.md` — observable symptom only
- `context/bugs/sec001/bug-context.md` — observable symptom only

---

## Phase 1 — Skills

**`skills/research-quality-measurement.md`**
- GOLD: all file:line refs verified, snippets match source, no discrepancies
- SILVER: ≥90% refs verified, minor discrepancies documented
- BRONZE: ≥70% refs verified, some snippets differ
- FAIL: <70% verified or critical discrepancies

**`skills/unit-tests-FIRST.md`**
- Fast: each test < 100ms
- Independent: no shared state between tests
- Repeatable: deterministic, no external deps
- Self-validating: explicit pass/fail assertions
- Timely: tests added in same PR as code change

---

## Phase 2 — Research Artifacts (manual, pre-pipeline)

Per each issue — written by developer, read by agents as input:

- `context/bugs/bug001/research/codebase-research.md` + `implementation-plan.md`
- `context/bugs/bug002/research/codebase-research.md` + `implementation-plan.md`
- `context/bugs/sec001/research/codebase-research.md` + `implementation-plan.md`

`codebase-research.md` — exact file:line references, code snippets, root cause analysis.
`implementation-plan.md` — before/after code per file, test command.

---

## Phase 3 — Agent Files

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

**`agents/research-verifier.agent.md`** — model: `claude-opus-4-8`
- Input: `$BUG_DIR/research/codebase-research.md`, source files
- Uses skill: `skills/research-quality-measurement.md`
- Output: `$BUG_DIR/research/verified-research.md`

**`agents/bug-fixer.agent.md`** — model: `claude-haiku-4-5-20251001`
- Input: `$BUG_DIR/implementation-plan.md`, `$BUG_DIR/research/verified-research.md`
- Applies fixes to `src/main.go`
- Runs: `docker compose run --rm app go test ./...`
- Output: `$BUG_DIR/fix-summary.md`

**`agents/security-verifier.agent.md`** — model: `claude-opus-4-8`
- Input: `$BUG_DIR/fix-summary.md`, `src/main.go`
- Rates findings: CRITICAL/HIGH/MEDIUM/LOW/INFO
- Output: `$BUG_DIR/security-report.md` (no code edits)

**`agents/unit-test-generator.agent.md`** — model: `claude-sonnet-4-6`
- Input: `$BUG_DIR/fix-summary.md`, `src/main.go`
- Uses skill: `skills/unit-tests-FIRST.md`
- Generates/updates tests in `src/main_test.go`
- Runs: `docker compose run --rm app go test ./...`
- Output: `$BUG_DIR/test-report.md`

---

## Phase 4 — Pipeline Runner

**`run-pipeline.sh`**
```bash
#!/bin/bash
set -e
BUG_DIR=${1:?Usage: ./run-pipeline.sh <bug-dir>}

claude --agent agents/research-verifier.agent.md --context "$BUG_DIR"
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

## Phase 5 — Documentation

**`README.md`** — author (Roman Reznik), overview, pipeline usage, app run command
**`HOWTORUN.md`** — prerequisites (Docker only), run app, run pipeline per issue

---

## Execution Order

```
1.1  skills/research-quality-measurement.md
1.2  skills/unit-tests-FIRST.md
2.1  context/bugs/bug001/research/codebase-research.md
2.2  context/bugs/bug001/implementation-plan.md
2.3  context/bugs/bug002/research/codebase-research.md
2.4  context/bugs/bug002/implementation-plan.md
2.5  context/bugs/sec001/research/codebase-research.md
2.6  context/bugs/sec001/implementation-plan.md
3.1  agents/research-verifier.agent.md
3.2  agents/bug-fixer.agent.md
3.3  agents/security-verifier.agent.md
3.4  agents/unit-test-generator.agent.md
4.1  run-pipeline.sh
5.1  README.md
5.2  HOWTORUN.md
```

---

## Verification

1. `docker compose up --build` — server starts on :8080
2. `curl http://localhost:8080/time` — returns wrong timezone (demonstrates bug#1)
3. `./run-pipeline.sh context/bugs/bug001` — runs end-to-end, all artifacts created
4. `./run-pipeline.sh context/bugs/bug002` — same
5. `./run-pipeline.sh context/bugs/sec001` — same
6. `curl http://localhost:8080/time` — returns correct UTC after fixes applied
7. All 6 expected artifact files exist in each of the 3 bug dirs