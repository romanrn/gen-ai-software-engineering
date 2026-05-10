package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTicket(t *testing.T) {
	metadata := TicketMetadata{
		Source:     SourceWebForm,
		Browser:    "Chrome",
		DeviceType: DeviceTypeDesktop,
	}

	ticket := NewTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Test Subject",
		"This is a test description",
		CategoryAccountAccess,
		PriorityHigh,
		metadata,
	)

	assert.NotNil(t, ticket)
	assert.NotEmpty(t, ticket.ID)
	assert.Equal(t, "CUST001", ticket.CustomerID)
	assert.Equal(t, "test@example.com", ticket.CustomerEmail)
	assert.Equal(t, "John Doe", ticket.CustomerName)
	assert.Equal(t, "Test Subject", ticket.Subject)
	assert.Equal(t, "This is a test description", ticket.Description)
	assert.Equal(t, CategoryAccountAccess, ticket.Category)
	assert.Equal(t, PriorityHigh, ticket.Priority)
	assert.Equal(t, StatusNew, ticket.Status)
	assert.NotNil(t, ticket.CreatedAt)
	assert.NotNil(t, ticket.UpdatedAt)
	assert.Nil(t, ticket.ResolvedAt)
	assert.Nil(t, ticket.AssignedTo)
	assert.Equal(t, 0, len(ticket.Tags))
	assert.Equal(t, metadata, ticket.Metadata)
}

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		category string
		valid    bool
	}{
		{"account_access", true},
		{"technical_issue", true},
		{"billing_question", true},
		{"feature_request", true},
		{"bug_report", true},
		{"other", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		result := IsValidCategory(tt.category)
		assert.Equal(t, tt.valid, result, "category %s", tt.category)
	}
}

func TestIsValidPriority(t *testing.T) {
	tests := []struct {
		priority string
		valid    bool
	}{
		{"urgent", true},
		{"high", true},
		{"medium", true},
		{"low", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		result := IsValidPriority(tt.priority)
		assert.Equal(t, tt.valid, result, "priority %s", tt.priority)
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		status string
		valid  bool
	}{
		{"new", true},
		{"in_progress", true},
		{"waiting_customer", true},
		{"resolved", true},
		{"closed", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		result := IsValidStatus(tt.status)
		assert.Equal(t, tt.valid, result, "status %s", tt.status)
	}
}

func TestIsValidSource(t *testing.T) {
	tests := []struct {
		source string
		valid  bool
	}{
		{"web_form", true},
		{"email", true},
		{"api", true},
		{"chat", true},
		{"phone", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		result := IsValidSource(tt.source)
		assert.Equal(t, tt.valid, result, "source %s", tt.source)
	}
}

func TestIsValidDeviceType(t *testing.T) {
	tests := []struct {
		deviceType string
		valid      bool
	}{
		{"desktop", true},
		{"mobile", true},
		{"tablet", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		result := IsValidDeviceType(tt.deviceType)
		assert.Equal(t, tt.valid, result, "deviceType %s", tt.deviceType)
	}
}

func TestErrValidation(t *testing.T) {
	details := []ValidationDetail{
		{Field: "email", Message: "invalid email"},
		{Field: "subject", Message: "subject required"},
	}

	err := &ErrValidation{Details: details}
	assert.NotEmpty(t, err.Error())
	assert.Contains(t, err.Error(), "2")
}

func TestErrValidation_EmptyDetails(t *testing.T) {
	err := &ErrValidation{}
	assert.Equal(t, "validation error", err.Error())
}

func TestErrNotFound(t *testing.T) {
	err := &ErrNotFound{ID: "test-id"}
	assert.NotEmpty(t, err.Error())
	assert.Contains(t, err.Error(), "test-id")
}

func TestErrInvalidInput(t *testing.T) {
	err := &ErrInvalidInput{Message: "bad input"}
	assert.NotEmpty(t, err.Error())
	assert.Contains(t, err.Error(), "bad input")
}
