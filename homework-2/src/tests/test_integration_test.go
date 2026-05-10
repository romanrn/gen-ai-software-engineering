package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"support-tickets/internal/domain"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Full ticket lifecycle — create → get → update → delete
func TestIntegration_TicketLifecycle(t *testing.T) {
	app := newApp()

	// Create
	r1, err := app.Test(jsonRequest(http.MethodPost, "/tickets", validTicket))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, r1.StatusCode)
	var created domain.Ticket
	json.NewDecoder(r1.Body).Decode(&created)
	id := created.ID.String()

	// Get
	r2, _ := app.Test(jsonRequest(http.MethodGet, "/tickets/"+id, nil))
	assert.Equal(t, http.StatusOK, r2.StatusCode)

	// Update
	r3, _ := app.Test(jsonRequest(http.MethodPut, "/tickets/"+id,
		map[string]any{"status": "in_progress", "assigned_to": "agent@company.com"}))
	assert.Equal(t, http.StatusOK, r3.StatusCode)
	var updated domain.Ticket
	json.NewDecoder(r3.Body).Decode(&updated)
	assert.Equal(t, domain.StatusInProgress, updated.Status)

	// Delete
	r4, _ := app.Test(jsonRequest(http.MethodDelete, "/tickets/"+id, nil))
	assert.Equal(t, http.StatusNoContent, r4.StatusCode)

	// Confirm gone
	r5, _ := app.Test(jsonRequest(http.MethodGet, "/tickets/"+id, nil))
	assert.Equal(t, http.StatusNotFound, r5.StatusCode)
}

// Test 2: Create ticket → auto-classify → verify classification stored on ticket
func TestIntegration_CreateAndClassify(t *testing.T) {
	app := newApp()

	body := map[string]any{
		"customer_id": "C1", "customer_email": "u@e.com", "customer_name": "User",
		"subject":     "Billing charge error",
		"description": "I was charged twice. Need a refund for the duplicate payment.",
		"category":    "other", "priority": "medium",
		"metadata": map[string]any{"source": "api", "browser": "Chrome", "device_type": "desktop"},
	}
	r1, _ := app.Test(jsonRequest(http.MethodPost, "/tickets", body))
	var created domain.Ticket
	json.NewDecoder(r1.Body).Decode(&created)

	r2, err := app.Test(jsonRequest(http.MethodPost, "/tickets/"+created.ID.String()+"/auto-classify", nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, r2.StatusCode)
	var c domain.Classification
	json.NewDecoder(r2.Body).Decode(&c)
	assert.Equal(t, domain.CategoryBillingQuestion, c.Category)
}

// Test 3: Bulk import CSV fixture → all 50 records imported
func TestIntegration_BulkImportCSV(t *testing.T) {
	app := newApp()
	data, err := os.ReadFile(filepath.Join(fixturesDir(), "sample_tickets.csv"))
	require.NoError(t, err)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "sample_tickets.csv")
	io.Copy(fw, bytes.NewReader(data))
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/tickets/import", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, float64(50), result["Total"])
	assert.Equal(t, float64(50), result["Successful"])
}

// Test 4: Bulk import JSON fixture → all 20 records imported
func TestIntegration_BulkImportJSON(t *testing.T) {
	app := newApp()
	data, err := os.ReadFile(filepath.Join(fixturesDir(), "sample_tickets.json"))
	require.NoError(t, err)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "sample_tickets.json")
	io.Copy(fw, bytes.NewReader(data))
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/tickets/import", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, float64(20), result["Total"])
}

// Test 5: Import invalid data → partial success, Failed > 0
func TestIntegration_BulkImportPartialFailure(t *testing.T) {
	app := newApp()

	csvData := "customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type\n" +
		"C1,bad-email,User,Subject,Valid description with enough characters,account_access,high,web_form,Chrome,desktop\n" +
		"C2,good@example.com,User,Subject,Valid description with enough characters,billing_question,low,email,Firefox,mobile"

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
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, float64(2), result["Total"])
	assert.Equal(t, float64(1), result["Successful"])
	assert.Equal(t, float64(1), result["Failed"])
}

// Test 6: Bulk import CSV then auto-classify an imported ticket
func TestIntegration_BulkImportAndAutoClassify(t *testing.T) {
	app := newApp()

	csvData := "customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type\n" +
		"C1,user1@example.com,User One,Cannot login,Password not working. Account access denied.,account_access,high,web_form,Chrome,desktop\n" +
		"C2,user2@example.com,User Two,Billing charge,I was charged twice. Need refund for duplicate payment.,billing_question,medium,email,Firefox,mobile"

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "tickets.csv")
	io.Copy(fw, strings.NewReader(csvData))
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/tickets/import", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	importResp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, importResp.StatusCode)

	var summary map[string]any
	json.NewDecoder(importResp.Body).Decode(&summary)
	require.Equal(t, float64(2), summary["Successful"])

	// Retrieve imported tickets and auto-classify the billing one
	listResp, _ := app.Test(jsonRequest(http.MethodGet, "/tickets", nil))
	var tickets []domain.Ticket
	json.NewDecoder(listResp.Body).Decode(&tickets)
	require.Len(t, tickets, 2)

	var billingID string
	for _, tk := range tickets {
		if tk.CustomerID == "C2" {
			billingID = tk.ID.String()
		}
	}
	require.NotEmpty(t, billingID)

	classResp, err := app.Test(jsonRequest(http.MethodPost, "/tickets/"+billingID+"/auto-classify", nil))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, classResp.StatusCode)

	var c domain.Classification
	json.NewDecoder(classResp.Body).Decode(&c)
	assert.Equal(t, domain.CategoryBillingQuestion, c.Category)
	assert.Greater(t, c.Confidence, 0.0)

	// Verify classification is persisted on the ticket
	getResp, _ := app.Test(jsonRequest(http.MethodGet, "/tickets/"+billingID, nil))
	var updated domain.Ticket
	json.NewDecoder(getResp.Body).Decode(&updated)
	require.NotNil(t, updated.Classification)
	assert.Equal(t, domain.CategoryBillingQuestion, updated.Classification.Category)
}

// Test 7: 25 concurrent POSTs — all must succeed, total stored == 25
func TestIntegration_ConcurrentCreate(t *testing.T) {
	const n = 25
	app := newApp()

	var wg sync.WaitGroup
	errCh := make(chan error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			body := map[string]any{
				"customer_id":    fmt.Sprintf("C%d", i),
				"customer_email": fmt.Sprintf("user%d@concurrent.example.com", i),
				"customer_name":  fmt.Sprintf("User %d", i),
				"subject":        "Concurrent test ticket",
				"description":    "Testing concurrent ticket creation with enough characters.",
				"category":       string(domain.CategoryTechnicalIssue),
				"priority":       string(domain.PriorityMedium),
			}
			resp, err := app.Test(jsonRequest(http.MethodPost, "/tickets", body))
			if err != nil {
				errCh <- fmt.Errorf("goroutine %d request error: %w", i, err)
				return
			}
			if resp.StatusCode != http.StatusCreated {
				errCh <- fmt.Errorf("goroutine %d got status %d", i, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Error(err)
	}

	// All 25 tickets must be stored
	listResp, err := app.Test(jsonRequest(http.MethodGet, "/tickets", nil))
	require.NoError(t, err)
	var tickets []domain.Ticket
	json.NewDecoder(listResp.Body).Decode(&tickets)
	assert.Len(t, tickets, n)
}

// Test 8: Combined filtering by category + priority
func TestIntegration_FilterByCategoryAndPriority(t *testing.T) {
	app := newApp()

	type combo struct{ category, priority string }
	seeds := []combo{
		{"technical_issue", "high"},
		{"technical_issue", "high"},
		{"technical_issue", "medium"},
		{"billing_question", "high"},
		{"billing_question", "low"},
	}

	for idx, s := range seeds {
		body := map[string]any{
			"customer_id":    fmt.Sprintf("FILTER-C%d", idx),
			"customer_email": fmt.Sprintf("filter%d@example.com", idx),
			"customer_name":  "Filter User",
			"subject":        "Filter test ticket",
			"description":    "Description used for combined filter testing.",
			"category":       s.category,
			"priority":       s.priority,
		}
		r, _ := app.Test(jsonRequest(http.MethodPost, "/tickets", body))
		require.Equal(t, http.StatusCreated, r.StatusCode)
	}

	// category only
	r, _ := app.Test(jsonRequest(http.MethodGet, "/tickets?category=technical_issue", nil))
	var byCategory []domain.Ticket
	json.NewDecoder(r.Body).Decode(&byCategory)
	assert.Len(t, byCategory, 3)

	// priority only
	r, _ = app.Test(jsonRequest(http.MethodGet, "/tickets?priority=high", nil))
	var byPriority []domain.Ticket
	json.NewDecoder(r.Body).Decode(&byPriority)
	assert.Len(t, byPriority, 3)

	// both: technical_issue + high → 2
	r, _ = app.Test(jsonRequest(http.MethodGet, "/tickets?category=technical_issue&priority=high", nil))
	var combined []domain.Ticket
	json.NewDecoder(r.Body).Decode(&combined)
	assert.Len(t, combined, 2)

	// both: billing_question + low → 1
	r, _ = app.Test(jsonRequest(http.MethodGet, "/tickets?category=billing_question&priority=low", nil))
	var single []domain.Ticket
	json.NewDecoder(r.Body).Decode(&single)
	assert.Len(t, single, 1)
}
