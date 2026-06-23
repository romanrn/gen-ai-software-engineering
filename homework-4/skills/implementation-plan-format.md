# Skill: Implementation Plan Format

Use this skill to structure the implementation plan that the Bug Fixer will execute.

## Output File

Write to: `$BUG_DIR/implementation-plan.md`

## Rules

- Before/After blocks must be minimal — change only what is necessary
- Do not introduce refactoring beyond the fix
- One change block per logical fix — do not bundle unrelated changes
- The test command must be runnable without modification

## Output Template

Copy this template exactly:

```markdown
# Implementation Plan: <issue-id>

## Overview
<one sentence: what is being fixed and why>

## Status
READY | BLOCKED: <reason if blocked>

## Changes

### File: src/main.go

**Before**
<exact current code>

**After**
<exact replacement code>

**Reason**: <one sentence explaining why this change fixes the bug>

## Test Command
docker compose run --rm app go test ./...

## Verification Steps
<manual steps to confirm the fix works at runtime, e.g.:>
1. `docker compose up --build`
2. `curl -H "X-Api-Key: changeme" http://localhost:8080/time`
3. Expected: <expected output>
```