package tests

import (
	"encoding/json"
	"net/http"
	"support-tickets/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validTicket = map[string]any{
	"customer_id":    "CUST-API-01",
	"customer_email": "api@example.com",
	"customer_name":  "API Tester",
	"subject":        "Cannot login to account",
	"description":    "I cannot access my account. Getting password error every time.",
	"category":       "account_access",
	"priority":       "high",
	"metadata":       map[string]any{"source": "web_form", "browser": "Chrome", "device_type": "desktop"},
}

func postTicket(t *testing.T) string {
	t.Helper()
	app := newApp()
	resp, err := app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var ticket domain.Ticket
	require.NoError(t, decodeJSON(resp.Body, &ticket))
	return ticket.ID.String()
}

// Test 1: POST /tickets → 201 with correct fields
func TestAPI_CreateTicket_201(t *testing.T) {
	app := newApp()
	resp, err := app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var ticket domain.Ticket
	require.NoError(t, decodeJSON(resp.Body, &ticket))
	assert.NotEmpty(t, ticket.ID)
	assert.Equal(t, "CUST-API-01", ticket.CustomerID)
	assert.Equal(t, domain.StatusNew, ticket.Status)
}

// Test 2: POST /tickets → 400 on validation failure
func TestAPI_CreateTicket_400(t *testing.T) {
	app := newApp()
	bad := map[string]any{"customer_email": "not-an-email", "subject": "", "priority": "bad"}
	resp, err := app.Test(jsonRequest(http.MethodPost, "/tickets", bad))
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Test 3: GET /tickets → 200, empty list
func TestAPI_ListTickets_Empty(t *testing.T) {
	app := newApp()
	resp, err := app.Test(jsonRequest(http.MethodGet, "/tickets", nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var list []domain.Ticket
	require.NoError(t, decodeJSON(resp.Body, &list))
	assert.Len(t, list, 0)
}

// Test 4: GET /tickets?category=account_access → 200, filtered
func TestAPI_ListTickets_Filter(t *testing.T) {
	app := newApp()
	// seed one ticket
	app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))

	resp, err := app.Test(jsonRequest(http.MethodGet, "/tickets?category=account_access", nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var list []domain.Ticket
	require.NoError(t, decodeJSON(resp.Body, &list))
	assert.Len(t, list, 1)
}

// Test 5: GET /tickets/:id → 200
func TestAPI_GetTicket_200(t *testing.T) {
	app := newApp()
	// seed + retrieve ID
	r1, _ := app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))
	var created domain.Ticket
	json.NewDecoder(r1.Body).Decode(&created)

	resp, err := app.Test(jsonRequest(http.MethodGet, "/tickets/"+created.ID.String(), nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test 6: GET /tickets/:id → 404 for unknown ID
func TestAPI_GetTicket_404(t *testing.T) {
	app := newApp()
	resp, err := app.Test(jsonRequest(http.MethodGet, "/tickets/unknown-id", nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Test 7: PUT /tickets/:id → 200, status updated
func TestAPI_UpdateTicket_200(t *testing.T) {
	app := newApp()
	r1, _ := app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))
	var created domain.Ticket
	json.NewDecoder(r1.Body).Decode(&created)

	resp, err := app.Test(jsonRequest(http.MethodPut, "/tickets/"+created.ID.String(),
		map[string]any{"status": "resolved"}))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updated domain.Ticket
	require.NoError(t, decodeJSON(resp.Body, &updated))
	assert.Equal(t, domain.StatusResolved, updated.Status)
}

// Test 8: PUT /tickets/:id → 404 for unknown ID
func TestAPI_UpdateTicket_404(t *testing.T) {
	app := newApp()
	resp, err := app.Test(jsonRequest(http.MethodPut, "/tickets/ghost-id",
		map[string]any{"status": "resolved"}))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Test 9: DELETE /tickets/:id → 204
func TestAPI_DeleteTicket_204(t *testing.T) {
	app := newApp()
	r1, _ := app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))
	var created domain.Ticket
	json.NewDecoder(r1.Body).Decode(&created)

	resp, err := app.Test(jsonRequest(http.MethodDelete, "/tickets/"+created.ID.String(), nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// Test 10: DELETE /tickets/:id → 404
func TestAPI_DeleteTicket_404(t *testing.T) {
	app := newApp()
	resp, err := app.Test(jsonRequest(http.MethodDelete, "/tickets/ghost-id", nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Test 11: POST /tickets/:id/auto-classify → 200
func TestAPI_AutoClassify_200(t *testing.T) {
	app := newApp()
	r1, _ := app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))
	var created domain.Ticket
	json.NewDecoder(r1.Body).Decode(&created)

	resp, err := app.Test(jsonRequest(http.MethodPost, "/tickets/"+created.ID.String()+"/auto-classify", nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var c domain.Classification
	require.NoError(t, decodeJSON(resp.Body, &c))
	assert.Equal(t, domain.CategoryAccountAccess, c.Category)
	assert.True(t, c.Confidence > 0)
}
