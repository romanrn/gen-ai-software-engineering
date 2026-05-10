package memory

import (
	"fmt"
	"sync"
	"time"
	"support-tickets/internal/domain"
	"support-tickets/internal/ports/out"
)

type TicketRepository struct {
	mu      sync.RWMutex
	tickets map[string]*domain.Ticket
}

func New() out.TicketRepository {
	return &TicketRepository{
		tickets: make(map[string]*domain.Ticket),
	}
}

func (r *TicketRepository) Save(ticket *domain.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if ticket == nil {
		return fmt.Errorf("ticket cannot be nil")
	}

	r.tickets[ticket.ID.String()] = ticket
	return nil
}

func (r *TicketRepository) FindByID(id string) (*domain.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ticket, ok := r.tickets[id]
	if !ok {
		return nil, &domain.ErrNotFound{ID: id}
	}

	return ticket, nil
}

func (r *TicketRepository) FindAll(filter *out.TicketFilter) ([]*domain.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Ticket

	for _, ticket := range r.tickets {
		if !r.matchesFilter(ticket, filter) {
			continue
		}
		results = append(results, ticket)
	}

	return results, nil
}

func (r *TicketRepository) Update(id string, updates map[string]interface{}) (*domain.Ticket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket, ok := r.tickets[id]
	if !ok {
		return nil, &domain.ErrNotFound{ID: id}
	}

	// Apply updates
	if status, ok := updates["status"]; ok {
		if s, ok := status.(domain.Status); ok {
			ticket.Status = s
		}
	}
	if priority, ok := updates["priority"]; ok {
		if p, ok := priority.(domain.Priority); ok {
			ticket.Priority = p
		}
	}
	if category, ok := updates["category"]; ok {
		if c, ok := category.(domain.Category); ok {
			ticket.Category = c
		}
	}
	if assignedTo, ok := updates["assigned_to"]; ok {
		if s, ok := assignedTo.(string); ok {
			ticket.AssignedTo = &s
		}
	}
	if tags, ok := updates["tags"]; ok {
		if t, ok := tags.([]string); ok {
			ticket.Tags = t
		}
	}
	if classification, ok := updates["classification"]; ok {
		if c, ok := classification.(*domain.Classification); ok {
			ticket.Classification = c
		}
	}

	// Update timestamp
	now := time.Now().UTC()
	ticket.UpdatedAt = now

	return ticket, nil
}

func (r *TicketRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tickets[id]; !ok {
		return &domain.ErrNotFound{ID: id}
	}

	delete(r.tickets, id)
	return nil
}

func (r *TicketRepository) matchesFilter(ticket *domain.Ticket, filter *out.TicketFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Status != nil && ticket.Status != domain.Status(*filter.Status) {
		return false
	}

	if filter.Category != nil && ticket.Category != domain.Category(*filter.Category) {
		return false
	}

	if filter.Priority != nil && ticket.Priority != domain.Priority(*filter.Priority) {
		return false
	}

	if filter.CustomerID != nil && ticket.CustomerID != *filter.CustomerID {
		return false
	}

	return true
}
