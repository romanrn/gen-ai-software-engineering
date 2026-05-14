package importer

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

type csvImporter struct{}

func (c *csvImporter) Parse(data []byte) ([]ImportRecord, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, nil
	}

	headers := records[0]
	var results []ImportRecord

	// Map headers to column indices
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}

	// Parse records
	for i := 1; i < len(records); i++ {
		row := records[i]

		// Ensure row has enough columns
		maxIdx := 0
		for _, idx := range headerMap {
			if idx > maxIdx {
				maxIdx = idx
			}
		}
		if len(row) <= maxIdx {
			continue
		}

		record := ImportRecord{
			CustomerID:    getField(row, headerMap, "customer_id"),
			CustomerEmail: getField(row, headerMap, "customer_email"),
			CustomerName:  getField(row, headerMap, "customer_name"),
			Subject:       getField(row, headerMap, "subject"),
			Description:   getField(row, headerMap, "description"),
			Category:      getField(row, headerMap, "category"),
			Priority:      getField(row, headerMap, "priority"),
			Metadata: ImportMetadata{
				Source:     getField(row, headerMap, "source"),
				Browser:    getField(row, headerMap, "browser"),
				DeviceType: getField(row, headerMap, "device_type"),
			},
		}

		results = append(results, record)
	}

	return results, nil
}

func getField(row []string, headerMap map[string]int, fieldName string) string {
	if idx, ok := headerMap[fieldName]; ok && idx < len(row) {
		return row[idx]
	}
	return ""
}
