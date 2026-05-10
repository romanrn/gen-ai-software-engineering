package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVImporter_ValidData(t *testing.T) {
	data := []byte(`customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type
CUST001,test@example.com,John Doe,Subject,This is a valid description with enough characters,account_access,high,web_form,Chrome,desktop`)

	imp := NewCSVImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "CUST001", records[0].CustomerID)
	assert.Equal(t, "test@example.com", records[0].CustomerEmail)
	assert.Equal(t, "account_access", records[0].Category)
}

func TestCSVImporter_MultipleRecords(t *testing.T) {
	data := []byte(`customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type
CUST001,test1@example.com,User One,Subject 1,Valid description with enough characters,account_access,high,web_form,Chrome,desktop
CUST002,test2@example.com,User Two,Subject 2,Valid description with enough characters,billing_question,low,email,Firefox,mobile`)

	imp := NewCSVImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 2, len(records))
}

func TestCSVImporter_EmptyFile(t *testing.T) {
	data := []byte("")

	imp := NewCSVImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Nil(t, records)
}

func TestCSVImporter_HeaderOnly(t *testing.T) {
	data := []byte(`customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type`)

	imp := NewCSVImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 0, len(records))
}

func TestCSVImporter_MissingColumns(t *testing.T) {
	data := []byte(`customer_id,customer_email
CUST001,test@example.com`)

	imp := NewCSVImporter()
	records, err := imp.Parse(data)

	// Should still parse but with missing fields
	require.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "CUST001", records[0].CustomerID)
	assert.Equal(t, "", records[0].Subject)
}
