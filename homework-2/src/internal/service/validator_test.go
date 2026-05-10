package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateTicket_ValidData(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreateTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Valid Subject",
		"This is a valid description with enough characters",
		"account_access",
		"high",
	)

	assert.Nil(t, err)
}

func TestValidateCreateTicket_EmptyCustomerID(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreateTicket(
		"",
		"test@example.com",
		"John Doe",
		"Subject",
		"Valid description with enough characters",
		"account_access",
		"high",
	)

	require.NotNil(t, err)
	assert.Equal(t, 1, len(err.Details))
	assert.Equal(t, "customer_id", err.Details[0].Field)
}

func TestValidateCreateTicket_InvalidEmail(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"not-an-email",
		"missing@domain",
		"@example.com",
		"user@",
	}

	for _, email := range tests {
		err := v.ValidateCreateTicket(
			"CUST001",
			email,
			"John Doe",
			"Subject",
			"Valid description with enough characters",
			"account_access",
			"high",
		)
		require.NotNil(t, err, "should fail for email: %s", email)
		assert.True(t, len(err.Details) > 0)
	}
}

func TestValidateCreateTicket_SubjectLength(t *testing.T) {
	v := NewValidator()

	// Empty subject
	err := v.ValidateCreateTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		"",
		"Valid description with enough characters",
		"account_access",
		"high",
	)
	require.NotNil(t, err)

	// Too long subject (> 200 chars)
	longSubject := ""
	for i := 0; i < 201; i++ {
		longSubject += "a"
	}
	err = v.ValidateCreateTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		longSubject,
		"Valid description with enough characters",
		"account_access",
		"high",
	)
	require.NotNil(t, err)
}

func TestValidateCreateTicket_DescriptionLength(t *testing.T) {
	v := NewValidator()

	// Too short description (< 10 chars)
	err := v.ValidateCreateTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Subject",
		"Short",
		"account_access",
		"high",
	)
	require.NotNil(t, err)

	// Too long description (> 2000 chars)
	longDesc := ""
	for i := 0; i < 2001; i++ {
		longDesc += "a"
	}
	err = v.ValidateCreateTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Subject",
		longDesc,
		"account_access",
		"high",
	)
	require.NotNil(t, err)
}

func TestValidateCreateTicket_InvalidCategory(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreateTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Subject",
		"Valid description with enough characters",
		"invalid_category",
		"high",
	)

	require.NotNil(t, err)
	assert.True(t, len(err.Details) > 0)
}

func TestValidateCreateTicket_InvalidPriority(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreateTicket(
		"CUST001",
		"test@example.com",
		"John Doe",
		"Subject",
		"Valid description with enough characters",
		"account_access",
		"invalid_priority",
	)

	require.NotNil(t, err)
	assert.True(t, len(err.Details) > 0)
}

func TestValidateCreateTicket_MultipleErrors(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreateTicket(
		"",
		"bad-email",
		"",
		"",
		"Short",
		"invalid",
		"bad",
	)

	require.NotNil(t, err)
	assert.True(t, len(err.Details) > 1, "should have multiple validation errors")
}

func TestValidCategories(t *testing.T) {
	tests := []struct {
		name      string
		category  string
		shouldErr bool
	}{
		{"account_access", "account_access", false},
		{"technical_issue", "technical_issue", false},
		{"billing_question", "billing_question", false},
		{"feature_request", "feature_request", false},
		{"bug_report", "bug_report", false},
		{"other", "other", false},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			err := v.ValidateCreateTicket(
				"CUST001",
				"test@example.com",
				"Name",
				"Subject",
				"Valid description with enough characters",
				tt.category,
				"high",
			)

			if tt.shouldErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestValidPriorities(t *testing.T) {
	tests := []string{"urgent", "high", "medium", "low"}

	for _, priority := range tests {
		t.Run(priority, func(t *testing.T) {
			v := NewValidator()
			err := v.ValidateCreateTicket(
				"CUST001",
				"test@example.com",
				"Name",
				"Subject",
				"Valid description with enough characters",
				"account_access",
				priority,
			)

			require.Nil(t, err)
		})
	}
}
