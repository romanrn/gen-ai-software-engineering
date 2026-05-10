package importer

import (
	"encoding/json"
	"fmt"
)

type jsonImporter struct{}

type jsonRecord struct {
	CustomerID    string `json:"customer_id"`
	CustomerEmail string `json:"customer_email"`
	CustomerName  string `json:"customer_name"`
	Subject       string `json:"subject"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	Priority      string `json:"priority"`
	Metadata      struct {
		Source     string `json:"source"`
		Browser    string `json:"browser"`
		DeviceType string `json:"device_type"`
	} `json:"metadata"`
}

func (j *jsonImporter) Parse(data []byte) ([]ImportRecord, error) {
	var jsonRecords []jsonRecord
	if err := json.Unmarshal(data, &jsonRecords); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var results []ImportRecord
	for _, jrec := range jsonRecords {
		record := ImportRecord{
			CustomerID:    jrec.CustomerID,
			CustomerEmail: jrec.CustomerEmail,
			CustomerName:  jrec.CustomerName,
			Subject:       jrec.Subject,
			Description:   jrec.Description,
			Category:      jrec.Category,
			Priority:      jrec.Priority,
			Metadata: ImportMetadata{
				Source:     jrec.Metadata.Source,
				Browser:    jrec.Metadata.Browser,
				DeviceType: jrec.Metadata.DeviceType,
			},
		}
		results = append(results, record)
	}

	return results, nil
}
