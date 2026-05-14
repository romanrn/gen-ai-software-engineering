package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXMLImporter_ValidData(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<tickets>
  <ticket>
    <customer_id>CUST001</customer_id>
    <customer_email>test@example.com</customer_email>
    <customer_name>John Doe</customer_name>
    <subject>Test Subject</subject>
    <description>This is a valid description with enough characters</description>
    <category>account_access</category>
    <priority>high</priority>
    <metadata>
      <source>web_form</source>
      <browser>Chrome</browser>
      <device_type>desktop</device_type>
    </metadata>
  </ticket>
</tickets>`)

	imp := NewXMLImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "CUST001", records[0].CustomerID)
	assert.Equal(t, "account_access", records[0].Category)
}

func TestXMLImporter_MultipleRecords(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<tickets>
  <ticket>
    <customer_id>CUST001</customer_id>
    <customer_email>test1@example.com</customer_email>
    <customer_name>User One</customer_name>
    <subject>Subject 1</subject>
    <description>Valid description with enough characters</description>
    <category>billing_question</category>
    <priority>low</priority>
    <metadata>
      <source>email</source>
      <browser>Chrome</browser>
      <device_type>mobile</device_type>
    </metadata>
  </ticket>
  <ticket>
    <customer_id>CUST002</customer_id>
    <customer_email>test2@example.com</customer_email>
    <customer_name>User Two</customer_name>
    <subject>Subject 2</subject>
    <description>Valid description with enough characters</description>
    <category>technical_issue</category>
    <priority>high</priority>
    <metadata>
      <source>api</source>
      <browser>Firefox</browser>
      <device_type>desktop</device_type>
    </metadata>
  </ticket>
</tickets>`)

	imp := NewXMLImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 2, len(records))
}

func TestXMLImporter_EmptyTickets(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<tickets>
</tickets>`)

	imp := NewXMLImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 0, len(records))
}

func TestXMLImporter_InvalidXML(t *testing.T) {
	data := []byte(`<tickets><ticket>unclosed`)

	imp := NewXMLImporter()
	records, err := imp.Parse(data)

	require.NotNil(t, err)
	assert.Nil(t, records)
}

func TestXMLImporter_MissingFields(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<tickets>
  <ticket>
    <customer_id>CUST001</customer_id>
    <customer_email>test@example.com</customer_email>
  </ticket>
</tickets>`)

	imp := NewXMLImporter()
	records, err := imp.Parse(data)

	require.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "CUST001", records[0].CustomerID)
	assert.Equal(t, "", records[0].Subject)
}
