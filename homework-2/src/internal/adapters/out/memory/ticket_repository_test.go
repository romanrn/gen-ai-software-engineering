package memory

import (
	"sync"
	"testing"
	"support-tickets/internal/domain"
	"support-tickets/internal/ports/out"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_Save(t *testing.T) {
	repo := New()
	ticket := &domain.Ticket{
		ID:            uuid.New(),
		CustomerID:    "CUST001",
		CustomerEmail: "test@example.com",
		CustomerName:  "John Doe",
		Subject:       "Test",
		Description:   "Test description",
		Category:      domain.CategoryAccountAccess,
		Priority:      domain.PriorityHigh,
	}

	err := repo.Save(ticket)
	assert.Nil(t, err)
}

func TestMemoryRepository_FindByID(t *testing.T) {
	repo := New()
	ticket := &domain.Ticket{
		ID:            uuid.New(),
		CustomerID:    "CUST001",
		CustomerEmail: "test@example.com",
		CustomerName:  "John Doe",
		Subject:       "Test",
		Description:   "Test description",
		Category:      domain.CategoryAccountAccess,
		Priority:      domain.PriorityHigh,
	}

	repo.Save(ticket)

	found, err := repo.FindByID(ticket.ID.String())
	require.Nil(t, err)
	assert.Equal(t, ticket.ID, found.ID)
	assert.Equal(t, ticket.CustomerID, found.CustomerID)
}

func TestMemoryRepository_FindByID_NotFound(t *testing.T) {
	repo := New()
	found, err := repo.FindByID("nonexistent")

	require.NotNil(t, err)
	assert.Nil(t, found)
	var notFoundErr *domain.ErrNotFound
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestMemoryRepository_FindAll(t *testing.T) {
	repo := New()

	for i := 0; i < 3; i++ {
		ticket := &domain.Ticket{
			ID:            uuid.New(),
			CustomerID:    "CUST001",
			CustomerEmail: "test@example.com",
			CustomerName:  "John Doe",
			Subject:       "Test",
			Description:   "Test description",
			Category:      domain.CategoryAccountAccess,
			Priority:      domain.PriorityHigh,
		}
		repo.Save(ticket)
	}

	tickets, err := repo.FindAll(nil)
	require.Nil(t, err)
	assert.Equal(t, 3, len(tickets))
}

func TestMemoryRepository_FindAll_WithFilter(t *testing.T) {
	repo := New()

	ticket1 := &domain.Ticket{
		ID:            uuid.New(),
		CustomerID:    "CUST001",
		CustomerEmail: "test@example.com",
		CustomerName:  "John Doe",
		Subject:       "Test",
		Description:   "Test description",
		Category:      domain.CategoryAccountAccess,
		Priority:      domain.PriorityHigh,
		Status:        domain.StatusNew,
	}
	repo.Save(ticket1)

	ticket2 := &domain.Ticket{
		ID:            uuid.New(),
		CustomerID:    "CUST002",
		CustomerEmail: "test2@example.com",
		CustomerName:  "Jane Doe",
		Subject:       "Test 2",
		Description:   "Test description 2",
		Category:      domain.CategoryBillingQuestion,
		Priority:      domain.PriorityLow,
		Status:        domain.StatusResolved,
	}
	repo.Save(ticket2)

	status := string(domain.StatusNew)
	filter := &out.TicketFilter{Status: &status}
	tickets, err := repo.FindAll(filter)
	require.Nil(t, err)
	assert.Equal(t, 1, len(tickets))
	assert.Equal(t, ticket1.ID, tickets[0].ID)
}

func TestMemoryRepository_Update(t *testing.T) {
	repo := New()
	ticket := &domain.Ticket{
		ID:            uuid.New(),
		CustomerID:    "CUST001",
		CustomerEmail: "test@example.com",
		CustomerName:  "John Doe",
		Subject:       "Test",
		Description:   "Test description",
		Category:      domain.CategoryAccountAccess,
		Priority:      domain.PriorityHigh,
		Status:        domain.StatusNew,
	}
	repo.Save(ticket)

	updates := map[string]interface{}{
		"status":   domain.StatusResolved,
		"priority": domain.PriorityLow,
	}
	updated, err := repo.Update(ticket.ID.String(), updates)
	require.Nil(t, err)
	assert.Equal(t, domain.StatusResolved, updated.Status)
	assert.Equal(t, domain.PriorityLow, updated.Priority)
}

func TestMemoryRepository_Delete(t *testing.T) {
	repo := New()
	ticket := &domain.Ticket{
		ID:            uuid.New(),
		CustomerID:    "CUST001",
		CustomerEmail: "test@example.com",
		CustomerName:  "John Doe",
		Subject:       "Test",
		Description:   "Test description",
		Category:      domain.CategoryAccountAccess,
		Priority:      domain.PriorityHigh,
	}
	repo.Save(ticket)

	err := repo.Delete(ticket.ID.String())
	assert.Nil(t, err)

	found, err := repo.FindByID(ticket.ID.String())
	require.NotNil(t, err)
	assert.Nil(t, found)
}

func TestMemoryRepository_ConcurrentAccess(t *testing.T) {
	repo := New()
	var wg sync.WaitGroup

	// Concurrent saves
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ticket := &domain.Ticket{
				ID:            uuid.New(),
				CustomerID:    "CUST001",
				CustomerEmail: "test@example.com",
				CustomerName:  "John Doe",
				Subject:       "Test",
				Description:   "Test description",
				Category:      domain.CategoryAccountAccess,
				Priority:      domain.PriorityHigh,
			}
			repo.Save(ticket)
		}(i)
	}

	wg.Wait()

	tickets, err := repo.FindAll(nil)
	require.Nil(t, err)
	assert.Equal(t, 10, len(tickets))
}
