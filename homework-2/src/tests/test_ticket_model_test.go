package tests

import (
	"support-tickets/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: NewTicket sets all defaults correctly
func TestModel_NewTicket_Defaults(t *testing.T) {
	ticket := domain.NewTicket("C1", "u@e.com", "User", "Subject", "Description",
		domain.CategoryAccountAccess, domain.PriorityHigh, domain.TicketMetadata{})

	require.NotNil(t, ticket)
	assert.NotEmpty(t, ticket.ID)
	assert.Equal(t, domain.StatusNew, ticket.Status)
	assert.WithinDuration(t, time.Now(), ticket.CreatedAt, 2*time.Second)
	assert.Nil(t, ticket.ResolvedAt)
	assert.Nil(t, ticket.AssignedTo)
	assert.Empty(t, ticket.Tags)
}

// Test 2–7: Valid enum values accepted
func TestModel_ValidCategories(t *testing.T) {
	valid := []string{"account_access", "technical_issue", "billing_question",
		"feature_request", "bug_report", "other"}
	for _, v := range valid {
		assert.True(t, domain.IsValidCategory(v), v)
	}
}

func TestModel_InvalidCategory(t *testing.T) {
	assert.False(t, domain.IsValidCategory("unknown"))
	assert.False(t, domain.IsValidCategory(""))
}

func TestModel_ValidPriorities(t *testing.T) {
	for _, v := range []string{"urgent", "high", "medium", "low"} {
		assert.True(t, domain.IsValidPriority(v), v)
	}
}

func TestModel_InvalidPriority(t *testing.T) {
	assert.False(t, domain.IsValidPriority("critical"))
	assert.False(t, domain.IsValidPriority(""))
}

func TestModel_ValidStatuses(t *testing.T) {
	for _, v := range []string{"new", "in_progress", "waiting_customer", "resolved", "closed"} {
		assert.True(t, domain.IsValidStatus(v), v)
	}
}

func TestModel_ValidSources(t *testing.T) {
	for _, v := range []string{"web_form", "email", "api", "chat", "phone"} {
		assert.True(t, domain.IsValidSource(v), v)
	}
}

func TestModel_ValidDeviceTypes(t *testing.T) {
	for _, v := range []string{"desktop", "mobile", "tablet"} {
		assert.True(t, domain.IsValidDeviceType(v), v)
	}
}

// Test 8: ErrValidation carries field-level details
func TestModel_ErrValidation_Details(t *testing.T) {
	err := &domain.ErrValidation{Details: []domain.ValidationDetail{
		{Field: "email", Message: "invalid"},
	}}
	assert.Contains(t, err.Error(), "1")
}

// Test 9: ErrNotFound and ErrInvalidInput error messages
func TestModel_ErrorMessages(t *testing.T) {
	notFound := &domain.ErrNotFound{ID: "abc-123"}
	assert.Contains(t, notFound.Error(), "abc-123")

	invalid := &domain.ErrInvalidInput{Message: "bad value"}
	assert.Contains(t, invalid.Error(), "bad value")
}
