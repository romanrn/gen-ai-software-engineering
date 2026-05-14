package in

import "support-tickets/internal/domain"

type TicketListRequest struct {
	Status     *string
	Category   *string
	Priority   *string
	CustomerID *string
}

type TicketService interface {
	Create(customerID, email, name, subject, description string, category domain.Category, priority domain.Priority, metadata domain.TicketMetadata) (*domain.Ticket, error)
	BulkImport(format string, data []byte) (*BulkImportResult, error)
	List(req *TicketListRequest) ([]*domain.Ticket, error)
	GetByID(id string) (*domain.Ticket, error)
	Update(id string, updates map[string]interface{}) (*domain.Ticket, error)
	Delete(id string) error
	AutoClassify(id string) (*domain.Classification, error)
}

type BulkImportResult struct {
	Total      int
	Successful int
	Failed     int
	Errors     []ImportError
}

type ImportError struct {
	Row     int
	Message string
}
