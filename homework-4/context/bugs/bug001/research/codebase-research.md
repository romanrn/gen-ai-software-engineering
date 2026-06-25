# Codebase Research: bug001 — Wrong Timezone

## Summary
`GET /time` is reported to return the server's local time instead of UTC, so on a server in UTC+3 the response is 3 hours ahead of actual UTC. The root cause is in `timeHandler`, which builds the timestamp with `time.Now()` (local time) and never converts it to UTC, while still labeling the value as `"utc"` in the JSON response.

## Root Cause
At `src/main.go:43` the handler calls `now := time.Now()`. In Go, `time.Now()` returns the current time in the machine's **local** time zone. At `src/main.go:45` this value is formatted with `now.Format(time.RFC3339)` and emitted under the JSON key `"utc"`. Because `now` carries the local zone offset, `time.RFC3339` renders that local offset (e.g. `+03:00`) rather than `Z`/UTC. No `.UTC()` conversion is applied anywhere between obtaining the time and formatting it, so the response reflects the server's local timezone instead of UTC. The mismatch between the `"utc"` label and the local-zoned value is exactly the reported symptom.

## File References
- `src/main.go:43` — `now := time.Now()` obtains the current time in the server's local timezone (not UTC).
- `src/main.go:45` — formats `now` as RFC3339 and outputs it under the `"utc"` JSON key, propagating the local offset to clients.
- `src/main.go:37` — existing comment flags this as `BUG#2: time.Now() without .UTC()`.
- `src/main.go:38` — `timeHandler` function definition where the affected logic lives.

## Code Snippets
```go
// BUG#2: time.Now() without .UTC() — returns local server time, not UTC
func timeHandler(w http.ResponseWriter, r *http.Request) {
	if !validate(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	now := time.Now()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"utc": "%s"}`, now.Format(time.RFC3339))
}
```

## Impact
Any client calling `GET /time` (with a valid `X-Api-Key`) receives a timestamp labeled `"utc"` that actually carries the server's local timezone offset. On any server not configured to UTC, the value is wrong by the local offset (e.g. +3 hours on a UTC+3 host). This affects all consumers relying on the field for accurate UTC, and the defect is hard to detect when the server itself runs in UTC, surfacing only under cross-timezone testing or deployment.

## References
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/context/bugs/bug001/bug-context.md`
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/skills/codebase-research-format.md`
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/src/main.go`
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/src/main_test.go` (located during search; not relevant to root cause)
