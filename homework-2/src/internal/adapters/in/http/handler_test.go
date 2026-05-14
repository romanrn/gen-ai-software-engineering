package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"support-tickets/internal/adapters/out/memory"
	"support-tickets/internal/domain"
	"support-tickets/internal/service"
	"support-tickets/pkg/errorhandler"
	"support-tickets/pkg/middleware"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupApp wires up the full handler stack for testing.
func setupApp() *fiber.App {
	repo := memory.New()
	svc := service.NewTicketService(repo)
	hdlr := New(svc)

	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.Handler})
	app.Use(middleware.RequestID)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	tickets := app.Group("/tickets")
	tickets.Post("", hdlr.CreateTicket)
	tickets.Post("/import", hdlr.BulkImport)
	tickets.Get("", hdlr.ListTickets)
	tickets.Get("/:id", hdlr.GetTicket)
	tickets.Put("/:id", hdlr.UpdateTicket)
	tickets.Delete("/:id", hdlr.DeleteTicket)
	tickets.Post("/:id/auto-classify", hdlr.AutoClassify)

	return app
}

func jsonBody(v any) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

var validCreateBody = map[string]any{
	"customer_id":    "CUST001",
	"customer_email": "test@example.com",
	"customer_name":  "John Doe",
	"subject":        "Cannot login to account",
	"description":    "I cannot access my account. Getting password error.",
	"category":       "account_access",
	"priority":       "high",
	"metadata": map[string]any{
		"source":      "web_form",
		"browser":     "Chrome",
		"device_type": "desktop",
	},
}

// createTicket is a test helper that POSTs a valid ticket and returns its ID.
func createTicket(t *testing.T, app *fiber.App) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/tickets", jsonBody(validCreateBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var ticket domain.Ticket
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ticket))
	return ticket.ID.String()
}

// --- test_ticket_api (11 tests) ---

func TestCreateTicket_201(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest(http.MethodPost, "/tickets", jsonBody(validCreateBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var ticket domain.Ticket
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ticket))
	assert.NotEmpty(t, ticket.ID)
	assert.Equal(t, "CUST001", ticket.CustomerID)
	assert.Equal(t, domain.StatusNew, ticket.Status)
}

func TestCreateTicket_400_ValidationError(t *testing.T) {
	app := setupApp()
	body := map[string]any{
		"customer_id":    "",
		"customer_email": "not-an-email",
		"customer_name":  "",
		"subject":        "",
		"description":    "short",
		"category":       "invalid",
		"priority":       "invalid",
	}
	req := httptest.NewRequest(http.MethodPost, "/tickets", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestListTickets_200_Empty(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest(http.MethodGet, "/tickets", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tickets []domain.Ticket
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&tickets))
	assert.Equal(t, 0, len(tickets))
}

func TestListTickets_200_WithFilter(t *testing.T) {
	app := setupApp()
	createTicket(t, app)

	req := httptest.NewRequest(http.MethodGet, "/tickets?category=account_access", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tickets []domain.Ticket
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&tickets))
	assert.Equal(t, 1, len(tickets))
}

func TestGetTicket_200(t *testing.T) {
	app := setupApp()
	id := createTicket(t, app)

	req := httptest.NewRequest(http.MethodGet, "/tickets/"+id, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ticket domain.Ticket
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ticket))
	assert.Equal(t, id, ticket.ID.String())
}

func TestGetTicket_404(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest(http.MethodGet, "/tickets/nonexistent-id", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateTicket_200(t *testing.T) {
	app := setupApp()
	id := createTicket(t, app)

	status := "resolved"
	body := map[string]any{"status": status}
	req := httptest.NewRequest(http.MethodPut, "/tickets/"+id, jsonBody(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ticket domain.Ticket
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ticket))
	assert.Equal(t, domain.StatusResolved, ticket.Status)
}

func TestUpdateTicket_404(t *testing.T) {
	app := setupApp()
	body := map[string]any{"status": "resolved"}
	req := httptest.NewRequest(http.MethodPut, "/tickets/nonexistent-id", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteTicket_204(t *testing.T) {
	app := setupApp()
	id := createTicket(t, app)

	req := httptest.NewRequest(http.MethodDelete, "/tickets/"+id, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Confirm gone
	req2 := httptest.NewRequest(http.MethodGet, "/tickets/"+id, nil)
	resp2, _ := app.Test(req2)
	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
}

func TestDeleteTicket_404(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest(http.MethodDelete, "/tickets/nonexistent-id", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAutoClassify_200(t *testing.T) {
	app := setupApp()
	id := createTicket(t, app)

	req := httptest.NewRequest(http.MethodPost, "/tickets/"+id+"/auto-classify", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var c domain.Classification
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&c))
	assert.NotEmpty(t, c.Category)
	assert.True(t, c.Confidence >= 0 && c.Confidence <= 1)
}

func TestBulkImport_200_CSV(t *testing.T) {
	app := setupApp()

	csvData := "customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type\n" +
		"C1,a@example.com,Alice,Subject,Valid description with enough chars,account_access,high,web_form,Chrome,desktop\n" +
		"C2,b@example.com,Bob,Subject 2,Another valid description here too,billing_question,low,email,Firefox,mobile"

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "tickets.csv")
	io.Copy(fw, strings.NewReader(csvData))
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/tickets/import", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(t, float64(2), result["Total"])
	assert.Equal(t, float64(2), result["Successful"])
	assert.Equal(t, float64(0), result["Failed"])
}

func TestBulkImport_400_NoFile(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest(http.MethodPost, "/tickets/import", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAutoClassify_404(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest(http.MethodPost, "/tickets/nonexistent-id/auto-classify", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
