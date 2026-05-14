package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONImporter_ValidData(t *testing.T) {
	data := []byte(`[
  {
    "customer_id": "CUST001",
    "customer_email": "test@example.com",
    "customer_name": "John Doe",
    "subject": "Test Subject",
    "description": "This is a valid description with enough characters",
    "category": "account_access",
    "priority": "high",
    "metadata": {
      "source": "web_form",
      "browser": "Chrome",
      "device_type": "desktop"
    }
  }
]`)

	imp := NewJSONImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "CUST001", records[0].CustomerID)
	assert.Equal(t, "account_access", records[0].Category)
}

func TestJSONImporter_MultipleRecords(t *testing.T) {
	data := []byte(`[
  {
    "customer_id": "CUST001",
    "customer_email": "test1@example.com",
    "customer_name": "User One",
    "subject": "Subject 1",
    "description": "Valid description with enough characters",
    "category": "billing_question",
    "priority": "low",
    "metadata": {"source": "email", "browser": "Chrome", "device_type": "mobile"}
  },
  {
    "customer_id": "CUST002",
    "customer_email": "test2@example.com",
    "customer_name": "User Two",
    "subject": "Subject 2",
    "description": "Valid description with enough characters",
    "category": "technical_issue",
    "priority": "high",
    "metadata": {"source": "api", "browser": "Firefox", "device_type": "desktop"}
  }
]`)

	imp := NewJSONImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 2, len(records))
}

func TestJSONImporter_EmptyArray(t *testing.T) {
	data := []byte(`[]`)

	imp := NewJSONImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 0, len(records))
}

func TestJSONImporter_InvalidJSON(t *testing.T) {
	data := []byte(`[{invalid json}]`)

	imp := NewJSONImporter()
	records, err := imp.Parse(data)

	require.NotNil(t, err)
	assert.Nil(t, records)
}

func TestJSONImporter_MissingFields(t *testing.T) {
	data := []byte(`[
  {
    "customer_id": "CUST001",
    "customer_email": "test@example.com"
  }
]`)

	imp := NewJSONImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "CUST001", records[0].CustomerID)
	assert.Equal(t, "", records[0].Subject)
}
