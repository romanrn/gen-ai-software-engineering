# Codebase Research: bug002 — Incorrect Uptime Reporting

## Summary
The `/health` endpoint reports a `uptime_seconds` value that becomes negative and decreases after roughly 2 minutes of runtime. The root cause is that the uptime counter `uptimeSeconds` is declared as an `int8`, whose maximum value is 127. A background goroutine increments it once per second, so after 127 seconds the counter overflows past the `int8` maximum and wraps around to -128, then continues incrementing toward 0 — producing the observed negative, "decreasing then wrapping" values.

## Root Cause
`uptimeSeconds` is declared with type `int8` (`src/main.go:11`). A Go `int8` is a signed 8-bit integer with a value range of -128 to 127. The `tick()` goroutine sleeps one second and executes `uptimeSeconds++` in an infinite loop (`src/main.go:57`–`src/main.go:61`). After 127 increments (~127 seconds, i.e. just over 2 minutes), incrementing the counter overflows the `int8`: the value wraps from 127 to -128. The `healthHandler` reads this same variable and serializes it with the `%d` verb (`src/main.go:50`), so the endpoint reports the wrapped negative value. Because the goroutine keeps incrementing, the value climbs from -128 back up, then overflows again at 127 — appearing as uptime that "decreases" relative to a normal monotonically increasing counter.

## File References
- `src/main.go:11` — `var uptimeSeconds int8 = 0` declares the uptime counter as a signed 8-bit integer (max 127).
- `src/main.go:48`–`src/main.go:51` — `healthHandler` writes `uptimeSeconds` into the JSON response using `%d`.
- `src/main.go:50` — the line that serializes `uptime_seconds` from the overflow-prone variable.
- `src/main.go:57`–`src/main.go:61` — `tick()` goroutine increments `uptimeSeconds` every second in an unbounded loop.
- `src/main.go:60` — `uptimeSeconds++` is the increment that overflows once the value passes 127.
- `src/main.go:20` — `go tick()` starts the incrementing goroutine at startup.

## Code Snippets

Declaration (`src/main.go:10`–`src/main.go:11`):
```go
// BUG#1: uptimeSeconds uses int8 — overflows after 127 seconds
var uptimeSeconds int8 = 0
```

Health handler (`src/main.go:48`–`src/main.go:51`):
```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "uptime_seconds": %d}`, uptimeSeconds)
}
```

Tick goroutine (`src/main.go:57`–`src/main.go:62`):
```go
func tick() {
	for {
		time.Sleep(1 * time.Second)
		uptimeSeconds++
	}
}
```

Goroutine launch (`src/main.go:18`–`src/main.go:20`):
```go
func main() {
	startTime = time.Now()
	go tick()
```

## Impact
At runtime, the `/health` endpoint reports correct increasing uptime only for the first 127 seconds. After that, `uptime_seconds` wraps to -128 and oscillates between -128 and 127 (overflowing every 256 seconds). Any monitoring system, load balancer health check, dashboard, or alerting rule that consumes `uptime_seconds` will receive nonsensical negative or non-monotonic values, undermining uptime tracking and any logic that assumes uptime only grows. The service itself remains functional (requests are still served), so the defect is silent and only surfaces after the process has been alive for more than ~2 minutes, making it easy to miss in short-lived tests. Note that the `startTime` variable (`src/main.go:16`, set at `src/main.go:19`) is captured but never used to compute uptime, so the only source of the reported value is the overflow-prone counter.

## References
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/context/bugs/bug002/bug-context.md`
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/skills/codebase-research-format.md`
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/src/main.go`
