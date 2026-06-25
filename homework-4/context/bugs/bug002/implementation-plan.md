# Implementation Plan: bug002

## Overview
Fix `uptimeSeconds` overflow by changing its type from `int8` (max 127) to `int64`, preventing the counter from wrapping to negative values after ~127 seconds of server uptime.

## Status
READY

## Changes

### File: src/main.go

**Before**
```go
// BUG#1: uptimeSeconds uses int8 — overflows after 127 seconds
var uptimeSeconds int8 = 0
```

**After**
```go
var uptimeSeconds int64 = 0
```

**Reason**: `int8` overflows to -128 after 127 increments; replacing it with `int64` gives the counter a range of ~292 billion years, eliminating the overflow entirely without requiring any changes to `tick()` or `healthHandler`.

## Test Command
```
docker compose run --rm app go test ./...
```

## Verification Steps
1. `docker compose up --build`
2. Wait at least 130 seconds (or temporarily shorten the sleep in `tick()` to `1 * time.Millisecond` and revert after)
3. `curl http://localhost:8080/health`
4. Expected: `{"status": "ok", "uptime_seconds": <positive integer ≥ 130}` — the value must remain positive and increasing, never negative
