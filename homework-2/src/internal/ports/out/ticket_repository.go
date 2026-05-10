package out

import "support-tickets/internal/domain"

type TicketFilter struct {
	Status     *string
	Category   *string
	Priority   *string
	CustomerID *string
}

type TicketRepository interface {
	Save(ticket *domain.Ticket) error
	FindByID(id string) (*domain.Ticket, error)
	FindAll(filter *TicketFilter) ([]*domain.Ticket, error)
	Update(id string, updates map[string]interface{}) (*domain.Ticket, error)
	Delete(id string) error
}
