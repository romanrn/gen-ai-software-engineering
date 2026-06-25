# Security Report: bug002

## Summary
The bug002 fix changed `uptimeSeconds` from `int8` to `int64` (line 11), eliminating
the integer-overflow defect. This change is functionally safe and introduces no new
security weaknesses. However, the surrounding file still contains a pre-existing
**hardcoded credential** and an **insecure (non-constant-time) credential comparison**.
These were not introduced by the fix but remain present in the reviewed source and are
reported below per scan scope. Overall posture: the fix itself is clean, but the file
carries unresolved authentication-handling issues.

## Findings

### [HIGH] Hardcoded API key in source
- File: `src/main.go:14`
- Description: The API key `"supersecret-hardcoded-key-12345"` is hardcoded as a
  compile-time constant. Secrets committed to source are exposed to anyone with repo
  access, leak through version-control history, and cannot be rotated without a rebuild
  and redeploy. The inline comment (`SEC#1`) acknowledges it should come from the
  environment, but the code still uses the literal. Not introduced by the bug002 fix.
- Remediation: Load the key from an environment variable or a secrets manager at startup
  (e.g. `apiKey := os.Getenv("API_KEY")`), fail fast if it is unset, and remove the
  literal value from source and from git history.

### [MEDIUM] Non-constant-time credential comparison (timing side channel)
- File: `src/main.go:54`
- Description: `validate` compares the incoming `X-Api-Key` header to the secret using
  Go's `==` string operator, which short-circuits on the first differing byte. This
  leaks per-byte timing information that can let an attacker recover the key incrementally
  via a timing attack.
- Remediation: Use a constant-time comparison such as
  `crypto/subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Api-Key")), []byte(apiKey)) == 1`.

### [LOW] Missing authentication on /health endpoint
- File: `src/main.go:48-51`
- Description: `healthHandler` returns server state (`uptime_seconds`) without any
  authentication. This is common and often acceptable for health checks, but it does
  expose uptime/operational metadata to unauthenticated callers. Note: with the int64
  fix the field can no longer go negative, so no incorrect/garbage data is exposed.
- Remediation: If uptime is considered sensitive, restrict the endpoint to internal
  networks or require auth; otherwise document that exposure is intentional.

### [INFO] bug002 fix is security-safe
- File: `src/main.go:11`
- Description: The `int8` → `int64` change removes the overflow that previously caused
  `uptime_seconds` to wrap to negative values. `int64` will not overflow in any practical
  runtime. No injection, no untrusted input, and no new attack surface result from this
  change.
- Remediation: None required.

## Conclusion
ISSUES FOUND

The bug002 fix itself is **secure** and correctly resolves the integer-overflow defect.
The reviewed file, however, retains two pre-existing authentication weaknesses that
still need attention: the hardcoded API key (HIGH) and the non-constant-time comparison
(MEDIUM). These are outside the bug002 change scope but should be tracked and remediated.

## References
- Reviewed: `src/main.go` (full source after fix)
- Reviewed: `context/bugs/bug002/fix-summary.md`
