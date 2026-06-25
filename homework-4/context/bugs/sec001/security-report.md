# Security Report: sec#1

## Summary
The original CWE-798 (Use of Hard-Coded Credentials) vulnerability has been
remediated. The hardcoded API key `"supersecret-hardcoded-key-12345"` was
removed and the key is now sourced from the `API_KEY` environment variable at
startup (`var apiKey = os.Getenv("API_KEY")`), consistent with 12-factor
secret management. This is a genuine improvement.

However, the fix introduces and leaves exposed two authentication weaknesses in
the surrounding `validate()` logic that should be addressed before this is
considered fully secure: a fail-open authentication bypass when `API_KEY` is
unset/empty, and a non-constant-time credential comparison.

## Findings

### [HIGH] Authentication bypass when API_KEY is unset or empty (fail-open)
- File: `src/main.go:14`, `src/main.go:53-55`
- Description: `apiKey` is initialized to `os.Getenv("API_KEY")`, which returns
  an empty string if the variable is not set. `validate()` compares the
  incoming header against `apiKey` with `r.Header.Get("X-Api-Key") == apiKey`.
  `http.Request.Header.Get` also returns `""` when the header is absent.
  Therefore, if the server is started without `API_KEY` (or with it set to an
  empty value), every request — including those with no `X-Api-Key` header at
  all — passes validation. Authentication silently fails open instead of failing
  closed. This is the most serious residual risk: a deployment/config mistake
  results in a fully unauthenticated endpoint with no error or warning.
- Remediation: Validate the secret at startup and refuse to run (or reject all
  requests) when it is empty. For example, in `main()`:
  ```go
  apiKey = os.Getenv("API_KEY")
  if apiKey == "" {
      log.Fatal("API_KEY must be set")
  }
  ```
  Additionally, harden `validate()` to explicitly reject empty/missing keys:
  ```go
  func validate(r *http.Request) bool {
      key := r.Header.Get("X-Api-Key")
      if apiKey == "" || key == "" {
          return false
      }
      return subtle.ConstantTimeCompare([]byte(key), []byte(apiKey)) == 1
  }
  ```

### [LOW] Non-constant-time API key comparison (timing side channel)
- File: `src/main.go:54`
- Description: `validate()` compares the provided API key with the expected key
  using Go's `==` string operator. `==` short-circuits on the first differing
  byte, so its execution time leaks information about how many leading bytes
  matched. A network attacker who can measure response timing could, in theory,
  recover the key byte-by-byte. Practical exploitability over a network is
  limited by jitter, so this is rated LOW, but it is trivial to fix.
- Remediation: Use a constant-time comparison from the standard library:
  ```go
  import "crypto/subtle"
  // ...
  return subtle.ConstantTimeCompare(
      []byte(r.Header.Get("X-Api-Key")), []byte(apiKey)) == 1
  ```

### [INFO] Secret may be exposed via process environment / logs
- File: `src/main.go:14`
- Description: Moving the secret to an environment variable is the correct
  pattern, but environment variables can still leak through process listings
  (`/proc/<pid>/environ`), crash dumps, child processes, or accidental logging
  of the environment. No such logging occurs in the current code, so this is
  informational only.
- Remediation: Ensure the deployment platform stores `API_KEY` in a secrets
  manager, restricts access to the container/process environment, and never logs
  the full environment. Avoid passing the secret as a CLI argument.

### [INFO] /health endpoint is unauthenticated
- File: `src/main.go:48-51`
- Description: `healthHandler` does not call `validate()` and exposes
  `uptime_seconds`. This is normal and generally acceptable for a liveness probe
  and is not introduced by this fix; flagged only for completeness. The
  disclosed information (status + uptime) is low sensitivity.
- Remediation: No action required unless uptime/health is considered sensitive,
  in which case restrict the endpoint at the network/ingress layer.

## Conclusion
ISSUES FOUND

The original hardcoded-credential vulnerability (CWE-798) is correctly fixed.
Remaining work before this can be marked SECURE:
1. **HIGH** — Prevent the fail-open authentication bypass by requiring a
   non-empty `API_KEY` at startup and rejecting empty/missing keys in
   `validate()`. This is the key item to fix.
2. **LOW** — Switch the key comparison to `crypto/subtle.ConstantTimeCompare`
   to eliminate the timing side channel.

No injection, XSS/CSRF, unsafe dependency, or other hardcoded-secret issues were
identified in the reviewed scope. (Note: JSON responses are built with
`fmt.Fprintf` rather than `encoding/json`; current values are server-controlled
so there is no injection risk, but `encoding/json` is recommended for safety if
user-influenced data is ever added.)

## References
- `src/main.go` — full source reviewed (after fix)
- `context/bugs/sec001/fix-summary.md` — change description (lines 13-14)
