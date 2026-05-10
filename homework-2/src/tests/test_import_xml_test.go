package tests

import (
	"os"
	"path/filepath"
	"support-tickets/pkg/importer"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Parse fixture file (30 records)
func TestXML_FixtureFile(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(fixturesDir(), "sample_tickets.xml"))
	require.NoError(t, err)

	imp := importer.NewXMLImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 30, len(records))
}

// Test 2: All fields parsed correctly
func TestXML_FieldMapping(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<tickets>
  <ticket>
    <customer_id>C1</customer_id>
    <customer_email>u@e.com</customer_email>
    <customer_name>Alice</customer_name>
    <subject>Subject</subject>
    <description>Description</description>
    <category>bug_report</category>
    <priority>low</priority>
    <metadata><source>chat</source><browser>Safari</browser><device_type>tablet</device_type></metadata>
  </ticket>
</tickets>`)

	imp := importer.NewXMLImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "C1", records[0].CustomerID)
	assert.Equal(t, "bug_report", records[0].Category)
	assert.Equal(t, "chat", records[0].Metadata.Source)
	assert.Equal(t, "tablet", records[0].Metadata.DeviceType)
}

// Test 3: Empty <tickets> root returns zero records
func TestXML_EmptyRoot(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?><tickets></tickets>`)
	imp := importer.NewXMLImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	assert.Len(t, records, 0)
}

// Test 4: Malformed XML returns error
func TestXML_MalformedInput(t *testing.T) {
	imp := importer.NewXMLImporter()
	_, err := imp.Parse([]byte(`<tickets><ticket>unclosed`))
	assert.NotNil(t, err)
}

// Test 5: Missing child elements default to empty strings
func TestXML_MissingChildElements(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<tickets>
  <ticket>
    <customer_id>C1</customer_id>
    <customer_email>u@e.com</customer_email>
  </ticket>
</tickets>`)
	imp := importer.NewXMLImporter()
	records, err := imp.Parse(data)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "", records[0].Subject)
	assert.Equal(t, "", records[0].Category)
}
