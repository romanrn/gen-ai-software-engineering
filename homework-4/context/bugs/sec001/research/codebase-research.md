# Codebase Research: sec#1 — Exposed API Secret

## Summary
The reported symptom is that the API secret used to authenticate requests is visible in the source code repository, meaning anyone with read access can obtain it and the key cannot be rotated without modifying and redeploying the app. The suspected root cause is confirmed: the API key is declared as a hardcoded string constant directly in `src/main.go` and is committed to version control, rather than being read from an external source such as an environment variable or secret store.

## Root Cause
The authentication secret is defined at compile time as a package-level constant literal in the source file (`src/main.go:14`). The `validate` function compares the inbound `X-Api-Key` request header against this constant (`src/main.go:53-55`), so the credential lives entirely inside the source code. Because the value is a literal in a tracked Go source file, it is checked into version control and embedded into the compiled binary. There is no code path that loads the key from the environment or any external configuration — note that the program reads `PORT` from the environment (`src/main.go:25`) but does not do the same for the API key. As a result, the secret is exposed to every developer with repository read access, and changing it requires editing the source and redeploying the application.

## File References
- `src/main.go:13` — comment marking the hardcoded secret defect: `// SEC#1: hardcoded API key — should come from environment`
- `src/main.go:14` — the API key declared as a hardcoded constant string literal
- `src/main.go:38-42` — `timeHandler` rejects requests that fail `validate`, gating the endpoint on the hardcoded key
- `src/main.go:53-55` — `validate` compares the request's `X-Api-Key` header against the hardcoded `apiKey` constant
- `src/main.go:25-28` — `PORT` is read from the environment, showing env-based config exists but is not used for the secret

## Code Snippets
From `src/main.go:13-14`:
```go
// SEC#1: hardcoded API key — should come from environment
const apiKey = "supersecret-hardcoded-key-12345"
```

From `src/main.go:53-55`:
```go
func validate(r *http.Request) bool {
	return r.Header.Get("X-Api-Key") == apiKey
}
```

From `src/main.go:38-42`:
```go
func timeHandler(w http.ResponseWriter, r *http.Request) {
	if !validate(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
```

From `src/main.go:25-28` (contrast: env-based config that the secret does not use):
```go
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
```

## Impact
At runtime the application authenticates requests using a secret whose value is permanently recorded in the repository's source and git history, and compiled into the distributed binary. Anyone with read access to the repository (or to the binary) can extract `supersecret-hardcoded-key-12345` and make authenticated requests to `/time`. Rotating or revoking the credential is impossible without a code change and redeploy, and the secret remains recoverable from version-control history even after such a change. This corresponds to OWASP A07:2021 — secret exposure through version control. Affected parties: any developer or actor with repository/binary access, and all consumers relying on the endpoint's authentication.

## References
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/context/bugs/sec001/bug-context.md`
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/skills/codebase-research-format.md`
- `~/Training/gen-ai/gen-ai-software-engineering/homework-4/src/main.go`
