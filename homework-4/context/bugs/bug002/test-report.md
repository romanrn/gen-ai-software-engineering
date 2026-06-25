# Test Report: bug002

## Summary
Tests generated: 3
Tests passed: 3
Tests failed: 0

## Generated Tests

| Test name | F | I | R | S | T | Result |
|-----------|---|---|---|---|---|--------|
| TestUptimeCounterDoesNotOverflowPastInt8Max | ✓ | ✓ | ✓ | ✓ | ✓ | PASS |
| TestHealthHandlerReportsUptimeBeyondInt8Max | ✓ | ✓ | ✓ | ✓ | ✓ | PASS |
| TestHealthHandlerUptimeNonNegativeAtInt8Boundary | ✓ | ✓ | ✓ | ✓ | ✓ | PASS |

### FIRST Compliance Notes

| Test name | F (< 100ms) | I (no shared state) | R (no ext deps) | S (has assertions) | T (covers fix only) |
|-----------|-------------|----------------------|------------------|--------------------|----------------------|
| TestUptimeCounterDoesNotOverflowPastInt8Max | ✓ No sleep/IO | ✓ save/restore uptimeSeconds via defer | ✓ pure arithmetic, no network | ✓ checks value == 128 and >= 0 | ✓ targets int64 type change |
| TestHealthHandlerReportsUptimeBeyondInt8Max | ✓ httptest, no real HTTP | ✓ save/restore uptimeSeconds via defer | ✓ httptest recorder, no external calls | ✓ checks JSON value == 200 and >= 0 | ✓ targets int64 type change |
| TestHealthHandlerUptimeNonNegativeAtInt8Boundary | ✓ httptest, no real HTTP | ✓ save/restore uptimeSeconds via defer | ✓ httptest recorder, no external calls | ✓ checks JSON value == 128 | ✓ targets int64 type change |

## Test Output

```
=== RUN   TestTimeHandler
--- PASS: TestTimeHandler (0.00s)
=== RUN   TestHealthHandler
--- PASS: TestHealthHandler (0.00s)
=== RUN   TestValidate
--- PASS: TestValidate (0.00s)
=== RUN   TestTimeHandlerResponseBodyIsUTC
--- PASS: TestTimeHandlerResponseBodyIsUTC (0.00s)
=== RUN   TestTimeHandlerTimestampEndsWithZ
--- PASS: TestTimeHandlerTimestampEndsWithZ (0.00s)
=== RUN   TestUptimeCounterDoesNotOverflowPastInt8Max
--- PASS: TestUptimeCounterDoesNotOverflowPastInt8Max (0.00s)
=== RUN   TestHealthHandlerReportsUptimeBeyondInt8Max
--- PASS: TestHealthHandlerReportsUptimeBeyondInt8Max (0.00s)
=== RUN   TestHealthHandlerUptimeNonNegativeAtInt8Boundary
--- PASS: TestHealthHandlerUptimeNonNegativeAtInt8Boundary (0.00s)
PASS
ok  	homework4	0.001s
```

## References
- src/main_test.go
- context/bugs/bug002/fix-summary.md
