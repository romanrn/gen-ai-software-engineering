package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestTimeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "supersecret-hardcoded-key-12345")
	w := httptest.NewRecorder()

	timeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestValidate(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "supersecret-hardcoded-key-12345")
	if !validate(req) {
		t.Error("expected valid key to pass")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/time", nil)
	req2.Header.Set("X-Api-Key", "wrong-key")
	if validate(req2) {
		t.Error("expected wrong key to fail")
	}
}

// --- Tests for bug001 fix: time.Now() → time.Now().UTC() in timeHandler ---

// TestTimeHandlerResponseBodyIsUTC verifies that timeHandler returns a JSON
// body whose "utc" field parses to a time.Time in the UTC location.
// Covers: bug001 fix — time.Now().UTC() ensures the returned Location is UTC.
func TestTimeHandlerResponseBodyIsUTC(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "supersecret-hardcoded-key-12345")
	w := httptest.NewRecorder()

	timeHandler(w, req)

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	ts, ok := result["utc"]
	if !ok {
		t.Fatal("response JSON is missing the 'utc' key")
	}

	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Fatalf("'utc' value %q is not a valid RFC3339 timestamp: %v", ts, err)
	}

	if parsed.Location() != time.UTC {
		t.Errorf("expected UTC location, got %q — fix may be missing .UTC() call", parsed.Location().String())
	}
}

// TestTimeHandlerTimestampEndsWithZ verifies that the UTC timestamp string
// ends with "Z" (the RFC3339 UTC designator) and not with a numeric offset
// such as "+02:00", which would indicate local server time was used instead.
// Covers: bug001 fix — without .UTC(), Format(RFC3339) emits the local offset.
func TestTimeHandlerTimestampEndsWithZ(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "supersecret-hardcoded-key-12345")
	w := httptest.NewRecorder()

	timeHandler(w, req)

	body := strings.TrimSpace(w.Body.String())
	// Expected body shape: {"utc": "2006-01-02T15:04:05Z"}
	// The UTC RFC3339 designator is "Z"; a local timezone would appear as "+HH:MM".
	if !strings.Contains(body, `Z"}`) {
		t.Errorf("expected timestamp to end with UTC designator 'Z', got body: %s", body)
	}
}

// --- Tests for bug002 fix: int8 → int64 for uptimeSeconds variable ---

// TestUptimeCounterDoesNotOverflowPastInt8Max verifies that the uptimeSeconds
// counter does not exhibit int8 overflow behaviour at the 127→128 boundary.
// With the old int8 type, incrementing past 127 would silently wrap to -128.
// Covers: bug002 fix — var uptimeSeconds int64.
func TestUptimeCounterDoesNotOverflowPastInt8Max(t *testing.T) {
	original := uptimeSeconds
	defer func() { uptimeSeconds = original }()

	uptimeSeconds = 127
	uptimeSeconds++ // would overflow to -128 with int8

	if uptimeSeconds < 0 {
		t.Errorf("uptimeSeconds is negative (%d) after crossing 127 — indicates int8 overflow", uptimeSeconds)
	}
	if uptimeSeconds != 128 {
		t.Errorf("expected uptimeSeconds to be 128 after increment from 127, got %d", uptimeSeconds)
	}
}

// TestHealthHandlerReportsUptimeBeyondInt8Max verifies that healthHandler
// correctly serialises a large uptime value (> 127) that would be unrepresentable
// by int8. Before the fix the JSON field would contain a negative number.
// Covers: bug002 fix — var uptimeSeconds int64.
func TestHealthHandlerReportsUptimeBeyondInt8Max(t *testing.T) {
	original := uptimeSeconds
	defer func() { uptimeSeconds = original }()

	uptimeSeconds = 200 // exceeds int8 max of 127

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	rawUptime, ok := result["uptime_seconds"]
	if !ok {
		t.Fatal("response JSON is missing the 'uptime_seconds' key")
	}

	uptime, ok := rawUptime.(float64)
	if !ok {
		t.Fatalf("expected 'uptime_seconds' to be a number, got %T", rawUptime)
	}

	if uptime < 0 {
		t.Errorf("uptime_seconds is negative (%v) — indicates integer overflow from int8", uptime)
	}
	if int64(uptime) != 200 {
		t.Errorf("expected uptime_seconds to be 200, got %v", uptime)
	}
}

// TestHealthHandlerUptimeNonNegativeAtInt8Boundary verifies that the uptime
// value reported by healthHandler is exactly 128 — the first value that would
// overflow a signed int8 (wrapping to -128). With int64 it must remain positive.
// Covers: bug002 fix — var uptimeSeconds int64.
func TestHealthHandlerUptimeNonNegativeAtInt8Boundary(t *testing.T) {
	original := uptimeSeconds
	defer func() { uptimeSeconds = original }()

	uptimeSeconds = 128 // first value that overflows int8 (becomes -128 with int8)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	uptime, ok := result["uptime_seconds"].(float64)
	if !ok {
		t.Fatal("expected 'uptime_seconds' to be a number in JSON response")
	}

	if int64(uptime) != 128 {
		t.Errorf("expected uptime_seconds to be 128, got %v — int8 overflow may be present", uptime)
	}
}

// --- Tests for sec001 fix: hardcoded API key → os.Getenv("API_KEY") ---

// TestApiKeyEqualsEnvVar verifies that the package-level apiKey variable equals
// the current value of the API_KEY environment variable. Before the fix, apiKey
// was a compile-time const; now it is a var initialised via os.Getenv("API_KEY").
// If apiKey were ever re-hardcoded, this assertion would fail when the env var
// is absent or set to a different value.
// Covers: sec001 fix — var apiKey = os.Getenv("API_KEY")
func TestApiKeyEqualsEnvVar(t *testing.T) {
	want := os.Getenv("API_KEY")
	if apiKey != want {
		t.Errorf("apiKey %q does not match os.Getenv(\"API_KEY\") %q — secret may be hardcoded instead of loaded from environment",
			apiKey, want)
	}
}

// TestValidateUsesEnvSourcedKey verifies that validate() accepts a request
// whose X-Api-Key header matches a key that was set at runtime (not at compile
// time). The test controls apiKey directly to ensure independence from whatever
// API_KEY is set to in the test environment.
// Covers: sec001 fix — validate() enforces the env-sourced apiKey.
func TestValidateUsesEnvSourcedKey(t *testing.T) {
	original := apiKey
	defer func() { apiKey = original }()

	apiKey = "test-runtime-key-sec001"

	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "test-runtime-key-sec001")

	if !validate(req) {
		t.Error("validate() returned false for a key matching the runtime-configured apiKey")
	}
}

// TestValidateRejectsOldHardcodedKey verifies that the old hardcoded secret
// "supersecret-hardcoded-key-12345" is rejected when apiKey is set to a
// different runtime value. This guards against the removed constant being
// silently reintroduced as a fallback or default.
// Covers: sec001 fix — removal of const "supersecret-hardcoded-key-12345".
func TestValidateRejectsOldHardcodedKey(t *testing.T) {
	original := apiKey
	defer func() { apiKey = original }()

	const oldHardcoded = "supersecret-hardcoded-key-12345"
	apiKey = "different-runtime-key-sec001"

	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", oldHardcoded)

	if validate(req) {
		t.Errorf("validate() accepted the old hardcoded key %q — the secret must not be embedded in source code", oldHardcoded)
	}
}

// TestTimeHandlerReturns401WithWrongEnvKey verifies that timeHandler returns
// HTTP 401 Unauthorized when the X-Api-Key header does not match the
// runtime-configured key. This exercises the full request path using the
// env-sourced apiKey, not a compile-time constant.
// Covers: sec001 fix — timeHandler rejects mismatched runtime key.
func TestTimeHandlerReturns401WithWrongEnvKey(t *testing.T) {
	original := apiKey
	defer func() { apiKey = original }()

	apiKey = "correct-runtime-key-sec001"

	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "wrong-key")
	w := httptest.NewRecorder()

	timeHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for mismatched key, got %d", w.Code)
	}
}

// TestTimeHandlerReturns200WithCorrectEnvKey verifies that timeHandler returns
// HTTP 200 OK when the X-Api-Key header matches the runtime-configured key.
// This confirms the positive authentication path works correctly after the
// sec001 fix moved the key out of source code and into the environment.
// Covers: sec001 fix — timeHandler accepts correct runtime key.
func TestTimeHandlerReturns200WithCorrectEnvKey(t *testing.T) {
	original := apiKey
	defer func() { apiKey = original }()

	apiKey = "correct-runtime-key-sec001"

	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "correct-runtime-key-sec001")
	w := httptest.NewRecorder()

	timeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK for correct runtime key, got %d", w.Code)
	}
}