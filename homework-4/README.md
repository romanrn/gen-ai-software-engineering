# Homework-4 — 6-Agent Bug-Fix Pipeline


> **Student Name**: Roman Reznik
> **Date Submitted**: 25-06-2026
> **AI Tools Used**: [List tools, e.g., Claude Code]

---

---

## Overview

This homework implements a **6-agent pipeline** that automatically researches, verifies, plans, fixes, security-reviews, and tests bugs in a Go HTTP service — all driven by a single shell command.

The pipeline is **generic**: point it at any bug directory containing a `bug-context.md` and it produces all six artifacts end-to-end with no manual steps between agents.

```
./run-pipeline.sh <bug-dir>

  Bug Researcher        → research/codebase-research.md
        ↓
  Research Verifier     → research/verified-research.md
        ↓
  Bug Planner           → implementation-plan.md
        ↓
  Bug Fixer             → fix-summary.md   (+ edits src/)
        ↓
  Security Verifier     → security-report.md
        ↓
  Unit Test Generator   → test-report.md   (+ extends src/main_test.go)
```

---

## The Application

A minimal Go HTTP service (`src/main.go`) seeded with three intentional issues:

| ID | Type | Description |
|----|------|-------------|
| bug001 | Bug | `GET /time` returns local server time instead of UTC |
| bug002 | Bug | Uptime counter uses `int8` — overflows to negative after ~2 min |
| sec001 | Security | API key hardcoded in source (`supersecret-hardcoded-key-12345`) |

### Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/time` | `X-Api-Key` header | Returns current UTC time |
| GET | `/health` | `X-Api-Key` header | Returns uptime in seconds |

---

## Agents

| # | Agent | File | Model | Role |
|---|-------|------|-------|------|
| 1 | Bug Researcher | `agents/bug-researcher.agent.md` | claude-opus-4-8 | Traces symptom to root cause in source |
| 2 | Research Verifier | `agents/research-verifier.agent.md` | claude-opus-4-8 | Fact-checks every file:line reference and rates quality |
| 3 | Bug Planner | `agents/bug-planner.agent.md` | claude-sonnet-4-6 | Turns verified research into a precise before/after fix plan |
| 4 | Bug Fixer | `agents/bug-fixer.agent.md` | claude-haiku-4-5-20251001 | Applies the plan, runs tests, documents changes |
| 5 | Security Verifier | `agents/security-verifier.agent.md` | claude-opus-4-8 | Scans changed code for vulnerabilities, rates CRITICAL→INFO |
| 6 | Unit Test Generator | `agents/unit-test-generator.agent.md` | claude-sonnet-4-6 | Generates FIRST-compliant tests for changed code, runs them |

### Model Justification

| Agent | Model | Reason |
|-------|-------|--------|
| Bug Researcher | claude-opus-4-8 | Root-cause analysis requires deep code reasoning |
| Research Verifier | claude-opus-4-8 | Fact-checking file:line references requires precision |
| Bug Planner | claude-sonnet-4-6 | Planning from verified facts — balanced capability |
| Bug Fixer | claude-haiku-4-5-20251001 | Routine mechanical edits; speed and cost matter |
| Security Verifier | claude-opus-4-8 | Security analysis requires nuanced judgment |
| Unit Test Generator | claude-sonnet-4-6 | Code understanding without max reasoning overhead |

---

## Skills

| Skill | Used by | Purpose |
|-------|---------|---------|
| `skills/codebase-research-format.md` | Bug Researcher | Structures research output |
| `skills/research-quality-measurement.md` | Research Verifier | GOLD/SILVER/BRONZE/FAIL rating levels |
| `skills/implementation-plan-format.md` | Bug Planner | Structures fix plan output |
| `skills/unit-tests-FIRST.md` | Unit Test Generator | FIRST checklist (Fast, Independent, Repeatable, Self-validating, Timely) |

---

## Running the Pipeline

See [HOWTORUN.md](HOWTORUN.md) for full prerequisites and step-by-step instructions.

```bash
# Run pipeline for each issue
./run-pipeline.sh context/bugs/bug001
./run-pipeline.sh context/bugs/bug002
./run-pipeline.sh context/bugs/sec001
```

Each run produces 6 artifact files inside the bug directory.

---

## Project Structure

```
homework-4/
├── README.md
├── HOWTORUN.md
├── run-pipeline.sh                      ← single-command pipeline runner
├── docker-compose.yml
├── src/
│   ├── main.go                          ← Go HTTP service (seeded issues)
│   ├── main_test.go                     ← test file (extended by pipeline)
│   ├── go.mod
│   └── Dockerfile
├── agents/                              ← 6 agent definitions
├── skills/                              ← 4 skill definitions
├── context/bugs/
│   ├── bug001/                          ← wrong timezone
│   ├── bug002/                          ← int8 overflow
│   └── sec001/                          ← hardcoded API key
└── docs/
    └── screenshots/                     ← pipeline run, fixes, security scan, tests
```