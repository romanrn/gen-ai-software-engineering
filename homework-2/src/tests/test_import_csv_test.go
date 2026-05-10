package tests

import (
	"os"
	"path/filepath"
	"support-tickets/pkg/importer"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Parse fixture file (50 records)
func TestCSV_FixtureFile(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(fixturesDir(), "sample_tickets.csv"))
	require.NoError(t, err)

	imp := importer.NewCSVImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 50, len(records))
}

// Test 2: All required fields mapped from headers
func TestCSV_FieldMapping(t *testing.T) {
	data := []byte("customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type\n" +
		"C1,user@example.com,Alice,Subject Line,Valid description text here,account_access,high,web_form,Chrome,desktop")

	imp := importer.NewCSVImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "C1", records[0].CustomerID)
	assert.Equal(t, "user@example.com", records[0].CustomerEmail)
	assert.Equal(t, "account_access", records[0].Category)
	assert.Equal(t, "web_form", records[0].Metadata.Source)
	assert.Equal(t, "Chrome", records[0].Metadata.Browser)
	assert.Equal(t, "desktop", records[0].Metadata.DeviceType)
}

// Test 3: Header-only file returns zero records
func TestCSV_HeaderOnly(t *testing.T) {
	data := []byte("customer_id,customer_email,customer_name,subject,description,category,priority,source,browser,device_type\n")
	imp := importer.NewCSVImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	assert.Len(t, records, 0)
}

// Test 4: Empty file returns nil, no error
func TestCSV_EmptyFile(t *testing.T) {
	imp := importer.NewCSVImporter()
	records, err := imp.Parse([]byte(""))
	require.NoError(t, err)
	assert.Nil(t, records)
}

// Test 5: Missing columns default to empty string
func TestCSV_MissingColumns(t *testing.T) {
	data := []byte("customer_id,customer_email\nC1,u@e.com")
	imp := importer.NewCSVImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "C1", records[0].CustomerID)
	assert.Equal(t, "", records[0].Subject)
}

// Test 6: Invalid CSV (unclosed quote) returns error
func TestCSV_MalformedFile(t *testing.T) {
	data := []byte("customer_id,customer_email\n\"unclosed,value")
	imp := importer.NewCSVImporter()
	_, err := imp.Parse(data)
	assert.NotNil(t, err)
}
