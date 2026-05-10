package importer

type ImportMetadata struct {
	Source     string
	Browser    string
	DeviceType string
}

type ImportRecord struct {
	CustomerID    string
	CustomerEmail string
	CustomerName  string
	Subject       string
	Description   string
	Category      string
	Priority      string
	Metadata      ImportMetadata
}

type TicketImporter interface {
	Parse(data []byte) ([]ImportRecord, error)
}

func NewCSVImporter() TicketImporter {
	return &csvImporter{}
}

func NewJSONImporter() TicketImporter {
	return &jsonImporter{}
}

func NewXMLImporter() TicketImporter {
	return &xmlImporter{}
}
