# Fix Summary: bug001

## Overall Status
PASSED

## Changes Made

### src/main.go
- Location: `timeHandler` function, line 43
- Before: `now := time.Now()`
- After: `now := time.Now().UTC()`
- Test Result: PASSED

## Test Output
```
ok  	homework4	0.708s
```

Tests executed successfully with: `go test ./...` from the src directory.

## Manual Verification

To verify the fix at runtime:

```bash
# Build and start the server
docker compose up --build

# In another terminal, call the /time endpoint with the required API key
curl -H "X-Api-Key: changeme" http://localhost:8080/time
```

Expected output: 
```json
{"utc": "<RFC3339 timestamp ending in Z>"}
```

Example: `{"utc": "2026-06-25T14:32:01Z"}`

The timestamp must end with `Z` (indicating UTC) rather than a numeric UTC offset like `+02:00`.

## Root Cause
`time.Now()` returns the current time in the server's local timezone. Without calling `.UTC()` on the result, the subsequent RFC3339 formatting will preserve the server's local timezone offset rather than converting to UTC. This means the JSON response under the `"utc"` key would incorrectly reflect local server time, not Coordinated Universal Time.

## Fix Applied
Added `.UTC()` method call to `time.Now()` at line 43 of `src/main.go`. This converts the returned time to UTC before formatting, ensuring the `"utc"` JSON field always contains a true UTC timestamp as intended.

## References
- Modified: `src/main.go` (1 line changed)
- Implementation plan: `context/bugs/bug001/implementation-plan.md`
- Research: `context/bugs/bug001/research/verified-research.md`
