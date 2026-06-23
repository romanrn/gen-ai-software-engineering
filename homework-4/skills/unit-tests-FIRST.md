# Skill: Unit Tests FIRST

Use this skill to evaluate and write unit tests that satisfy the FIRST principles.

## Principles

### F — Fast
Each test must complete in under 100ms.
No real network calls, no sleep, no file I/O outside the package under test.

### I — Independent
Tests must not depend on each other or on shared mutable state.
Each test sets up its own input and can run in any order or in isolation.

### R — Repeatable
Tests must produce the same result on every run regardless of environment, time, or external systems.
Use `httptest` package for HTTP handlers; mock time if needed.

### S — Self-validating
Every test must have an explicit assertion (`t.Errorf`, `t.Fatalf`, etc.).
A test that passes without asserting anything is not a test.

### T — Timely
Tests must be written for the changed code only — no tests for unchanged logic.
Tests are added in the same change as the code fix.

## Checklist

Before submitting generated tests, verify each test satisfies:

- [ ] Runs in < 100ms
- [ ] Does not read from or write to the filesystem or network
- [ ] Has at least one assertion
- [ ] Passes when run in isolation (`go test -run TestName`)
- [ ] Covers only changed/new code from `fix-summary.md`

## Output

Record FIRST compliance per test in `test-report.md`:

```
| Test name | F | I | R | S | T | Notes |
|-----------|---|---|---|---|---|-------|
| TestXxx   | ✓ | ✓ | ✓ | ✓ | ✓ |       |
```