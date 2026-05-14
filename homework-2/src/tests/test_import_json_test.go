package tests

import (
	"os"
	"path/filepath"
	"support-tickets/pkg/importer"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Parse fixture file (20 records)
func TestJSON_FixtureFile(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(fixturesDir(), "sample_tickets.json"))
	require.NoError(t, err)

	imp := importer.NewJSONImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 20, len(records))
}

// Test 2: All fields parsed correctly
func TestJSON_FieldMapping(t *testing.T) {
	data := []byte(`[{"customer_id":"C1","customer_email":"u@e.com","customer_name":"Alice",
		"subject":"Sub","description":"Desc","category":"technical_issue","priority":"urgent",
		"metadata":{"source":"api","browser":"Firefox","device_type":"mobile"}}]`)

	imp := importer.NewJSONImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "C1", records[0].CustomerID)
	assert.Equal(t, "technical_issue", records[0].Category)
	assert.Equal(t, "api", records[0].Metadata.Source)
}

// Test 3: Empty JSON array returns zero records
func TestJSON_EmptyArray(t *testing.T) {
	imp := importer.NewJSONImporter()
	records, err := imp.Parse([]byte(`[]`))
	require.NoError(t, err)
	assert.Len(t, records, 0)
}

// Test 4: Malformed JSON returns error
func TestJSON_MalformedInput(t *testing.T) {
	imp := importer.NewJSONImporter()
	_, err := imp.Parse([]byte(`[{not valid json}]`))
	assert.NotNil(t, err)
}

// Test 5: Missing optional fields default to zero values
func TestJSON_MissingOptionalFields(t *testing.T) {
	data := []byte(`[{"customer_id":"C1","customer_email":"u@e.com"}]`)
	imp := importer.NewJSONImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "", records[0].Subject)
	assert.Equal(t, "", records[0].Category)
}
