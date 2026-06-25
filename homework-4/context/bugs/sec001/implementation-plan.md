# Implementation Plan: sec#1

## Overview
Replace the hardcoded `apiKey` constant with a package-level variable loaded from the `API_KEY` environment variable, so the secret is never stored in source code.

## Status
READY

## Changes

### File: src/main.go

**Before**
```go
// SEC#1: hardcoded API key — should come from environment
const apiKey = "supersecret-hardcoded-key-12345"
```

**After**
```go
// SEC#1: API key loaded from environment variable API_KEY
var apiKey = os.Getenv("API_KEY")
```

**Reason**: Replacing the hardcoded constant with `os.Getenv("API_KEY")` removes the secret from source code; the `validate` function at line 53–55 already references `apiKey` by name and requires no further changes, and `os` is already imported.

## Test Command
docker compose run --rm app go test ./...

## Verification Steps
1. `API_KEY=mysecretkey docker compose up --build`
2. `curl -H "X-Api-Key: mysecretkey" http://localhost:8080/time`
3. Expected: `{"utc": "<current UTC timestamp>"}` with HTTP 200
4. `curl http://localhost:8080/time` (no key)
5. Expected: `Unauthorized` with HTTP 401
6. `curl -H "X-Api-Key: wrongkey" http://localhost:8080/time`
7. Expected: `Unauthorized` with HTTP 401
