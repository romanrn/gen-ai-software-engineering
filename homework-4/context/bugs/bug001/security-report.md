# Security Report: bug001

## Summary
The fix for bug001 was a single, security-neutral change: `time.Now()` →
`time.Now().UTC()` in `timeHandler` (line 43). This change introduces no new
vulnerabilities — it only normalizes timezone handling for the formatted
timestamp. However, the surrounding code that the fix touches and depends on
(`timeHandler` and its `validate` authentication path) contains pre-existing
security weaknesses that remain present after the fix. The most serious is a
hardcoded API key checked into source control.

## Findings

### [CRITICAL] Hardcoded API key (secret in source)
- File: `src/main.go:14`
- Description: The API key used to authorize `/time` requests is hardcoded as
  `const apiKey = "supersecret-hardcoded-key-12345"`. The code's own comment
  (`SEC#1`) flags that it should come from the environment. A secret committed
  to source control is exposed to anyone with repo access, cannot be rotated
  without a code change/redeploy, and is identical across all environments. The
  fix-summary even documents calling the endpoint with `X-Api-Key: changeme`,
  indicating inconsistent/placeholder secret handling.
- Remediation: Load the key from an environment variable or secrets manager at
  startup (e.g. `apiKey := os.Getenv("API_KEY")`) and fail fast if it is empty.
  Remove the literal from source and rotate the exposed value. Keep secrets out
  of version control (use `.env`/secret store, add to `.gitignore`).

### [MEDIUM] Non-constant-time API key comparison (timing side channel)
- File: `src/main.go:54`
- Description: `validate` compares the supplied key with `==`
  (`r.Header.Get("X-Api-Key") == apiKey`). Go's string `==` short-circuits on
  the first differing byte, leaking timing information that can, in principle,
  be used to recover the secret byte-by-byte.
- Remediation: Use a constant-time comparison:
  `crypto/subtle.ConstantTimeCompare([]byte(provided), []byte(apiKey)) == 1`.
  Consider hashing both sides to equal length first to avoid leaking length.

### [LOW] Missing request size / method restrictions on handlers
- File: `src/main.go:38`, `src/main.go:48`
- Description: `timeHandler` and `healthHandler` accept any HTTP method and do
  not constrain request bodies. There is no read timeout configured on the
  server (`http.ListenAndServe` uses default `http.Server` with no
  `ReadTimeout`/`WriteTimeout`), leaving the service exposed to slowloris-style
  resource exhaustion.
- Remediation: Use an explicit `http.Server` with `ReadHeaderTimeout`,
  `ReadTimeout`, and `WriteTimeout` set. Optionally reject unexpected methods
  (return 405 for non-GET).

### [INFO] Manual JSON string construction
- File: `src/main.go:45`, `src/main.go:50`
- Description: JSON responses are built with `fmt.Fprintf` and raw string
  interpolation. The interpolated values here (an RFC3339 timestamp and an
  integer) are not attacker-controlled, so there is no injection risk today.
  However, the pattern is fragile and would become an output-encoding/injection
  risk if any user-controlled string were ever inserted.
- Remediation: Marshal responses with `encoding/json` (e.g. `json.NewEncoder(w)`
  over a struct) to guarantee correct escaping.

### [INFO] Fix change (line 43) reviewed — no security impact
- File: `src/main.go:43`
- Description: The bug001 fix (`time.Now().UTC()`) is safe. It performs no I/O,
  takes no external input, and only affects timezone normalization of an
  internally-generated timestamp. No injection, secret, or validation concerns
  introduced.
- Remediation: None required.

## Conclusion
ISSUES FOUND

The bug001 fix itself is **secure** and introduces no vulnerabilities. However,
the file contains pre-existing security issues that still need attention:

1. CRITICAL — Remove the hardcoded API key (`src/main.go:14`) and source it from
   the environment/secrets manager; rotate the leaked value.
2. MEDIUM — Replace the `==` key comparison with a constant-time comparison
   (`src/main.go:54`).
3. LOW — Add server timeouts and method restrictions.

These are outside the scope of the bug001 one-line fix but should be tracked and
remediated, as they affect the same authentication path the fix relies on.

## References
- `src/main.go` (full source after fix, reviewed)
- `context/bugs/bug001/fix-summary.md` (change description: line 43)
