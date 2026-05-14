package importer

import (
	"encoding/xml"
	"fmt"
)

type xmlImporter struct{}

type xmlTickets struct {
	Tickets []xmlRecord `xml:"ticket"`
}

type xmlRecord struct {
	CustomerID    string `xml:"customer_id"`
	CustomerEmail string `xml:"customer_email"`
	CustomerName  string `xml:"customer_name"`
	Subject       string `xml:"subject"`
	Description   string `xml:"description"`
	Category      string `xml:"category"`
	Priority      string `xml:"priority"`
	Metadata      struct {
		Source     string `xml:"source"`
		Browser    string `xml:"browser"`
		DeviceType string `xml:"device_type"`
	} `xml:"metadata"`
}

func (x *xmlImporter) Parse(data []byte) ([]ImportRecord, error) {
	var xmlData xmlTickets
	if err := xml.Unmarshal(data, &xmlData); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	var results []ImportRecord
	for _, xrec := range xmlData.Tickets {
		record := ImportRecord{
			CustomerID:    xrec.CustomerID,
			CustomerEmail: xrec.CustomerEmail,
			CustomerName:  xrec.CustomerName,
			Subject:       xrec.Subject,
			Description:   xrec.Description,
			Category:      xrec.Category,
			Priority:      xrec.Priority,
			Metadata: ImportMetadata{
				Source:     xrec.Metadata.Source,
				Browser:    xrec.Metadata.Browser,
				DeviceType: xrec.Metadata.DeviceType,
			},
		}
		results = append(results, record)
	}

	return results, nil
}
