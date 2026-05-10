package service

import (
	"testing"
	"support-tickets/internal/adapters/out/memory"
	"support-tickets/internal/domain"
	"support-tickets/internal/ports/in"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTicketService_Create(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	ticket, err := svc.Create(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Test Subject",
		"This is a valid description with enough characters",
		domain.CategoryAccountAccess,
		domain.PriorityHigh,
		domain.TicketMetadata{
			Source:     domain.SourceWebForm,
			Browser:    "Chrome",
			DeviceType: domain.DeviceTypeDesktop,
		},
	)

	require.Nil(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "CUST001", ticket.CustomerID)
	assert.Equal(t, domain.CategoryAccountAccess, ticket.Category)
}

func TestTicketService_CreateWithInvalidData(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	ticket, err := svc.Create(
		"",
		"invalid",
		"",
		"",
		"short",
		domain.CategoryAccountAccess,
		domain.PriorityHigh,
		domain.TicketMetadata{},
	)

	require.NotNil(t, err)
	assert.Nil(t, ticket)
}

func TestTicketService_GetByID(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	created, _ := svc.Create(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Test Subject",
		"This is a valid description with enough characters",
		domain.CategoryAccountAccess,
		domain.PriorityHigh,
		domain.TicketMetadata{},
	)

	retrieved, err := svc.GetByID(created.ID.String())
	require.Nil(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
}

func TestTicketService_List(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	for i := 0; i < 3; i++ {
		svc.Create(
			"CUST001",
			"test@example.com",
			"John Doe",
			"Subject",
			"This is a valid description with enough characters",
			domain.CategoryAccountAccess,
			domain.PriorityHigh,
			domain.TicketMetadata{},
		)
	}

	tickets, err := svc.List(&in.TicketListRequest{})
	require.Nil(t, err)
	assert.Equal(t, 3, len(tickets))
}

func TestTicketService_Update(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	created, _ := svc.Create(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Test Subject",
		"This is a valid description with enough characters",
		domain.CategoryAccountAccess,
		domain.PriorityHigh,
		domain.TicketMetadata{},
	)

	updates := map[string]interface{}{
		"status":   domain.StatusResolved,
		"priority": domain.PriorityLow,
	}

	updated, err := svc.Update(created.ID.String(), updates)
	require.Nil(t, err)
	assert.Equal(t, domain.StatusResolved, updated.Status)
	assert.Equal(t, domain.PriorityLow, updated.Priority)
}

func TestTicketService_Delete(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	created, _ := svc.Create(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Test Subject",
		"This is a valid description with enough characters",
		domain.CategoryAccountAccess,
		domain.PriorityHigh,
		domain.TicketMetadata{},
	)

	err := svc.Delete(created.ID.String())
	require.Nil(t, err)

	_, err = svc.GetByID(created.ID.String())
	require.NotNil(t, err)
}

func TestTicketService_AutoClassify(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	created, _ := svc.Create(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Cannot login to account",
		"I cannot access my account. Getting password error.",
		domain.CategoryOther,
		domain.PriorityMedium,
		domain.TicketMetadata{},
	)

	classification, err := svc.AutoClassify(created.ID.String())
	require.Nil(t, err)
	assert.NotNil(t, classification)
	assert.Equal(t, domain.CategoryAccountAccess, classification.Category)
}

func TestTicketService_BulkImport_CSV(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	csvData := []byte(`customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type
CUST001,test@example.com,John Doe,Subject,This is a valid description with enough characters,account_access,high,web_form,Chrome,desktop
CUST002,test2@example.com,Jane Doe,Subject 2,This is a valid description with enough characters,billing_question,low,email,Firefox,mobile`)

	result, err := svc.BulkImport("csv", csvData)
	require.Nil(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, result.Successful)
	assert.Equal(t, 0, result.Failed)
}

func TestTicketService_BulkImport_JSON(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	jsonData := []byte(`[
  {
    "customer_id": "CUST001",
    "customer_email": "test@example.com",
    "customer_name": "John Doe",
    "subject": "Subject",
    "description": "This is a valid description with enough characters",
    "category": "account_access",
    "priority": "high",
    "metadata": {"source": "web_form", "browser": "Chrome", "device_type": "desktop"}
  }
]`)

	result, err := svc.BulkImport("json", jsonData)
	require.Nil(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, 1, result.Successful)
}

func TestTicketService_BulkImport_WithErrors(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	csvData := []byte(`customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type
CUST001,invalid-email,John Doe,Subject,This is a valid description with enough characters,account_access,high,web_form,Chrome,desktop
CUST002,test2@example.com,Jane Doe,Subject 2,This is a valid description with enough characters,billing_question,low,email,Firefox,mobile`)

	result, err := svc.BulkImport("csv", csvData)
	require.Nil(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 1, result.Successful)
	assert.Equal(t, 1, result.Failed)
	assert.True(t, len(result.Errors) > 0)
}

func TestTicketService_BulkImport_UnsupportedFormat(t *testing.T) {
	repo := memory.New()
	svc := NewTicketService(repo)

	result, err := svc.BulkImport("unsupported", []byte{})
	require.NotNil(t, err)
	assert.Nil(t, result)
}
