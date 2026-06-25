# Implementation Plan: bug001

## Overview
Fix `timeHandler` in `src/main.go` to call `.UTC()` on the result of `time.Now()` so the JSON response always carries a true UTC timestamp instead of the server's local-timezone offset.

## Status
READY

## Changes

### File: src/main.go

**Before**
```go
	now := time.Now()
```

**After**
```go
	now := time.Now().UTC()
```

**Reason**: `time.Now()` returns the current time in the server's local timezone; appending `.UTC()` converts it to UTC before RFC3339 formatting, ensuring the `"utc"` JSON field always reflects coordinated universal time regardless of server locale.

## Test Command
```
docker compose run --rm app go test ./...
```

## Verification Steps
1. `docker compose up --build`
2. `curl -H "X-Api-Key: changeme" http://localhost:8080/time`
3. Expected: `{"utc": "<timestamp ending in Z>"}` — the RFC3339 value must end with `Z` (e.g. `2026-06-25T14:32:01Z`), not a numeric UTC offset such as `+02:00`
