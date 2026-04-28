package memory

import (
	"context"
	"sync"

	"github.com/reznikrn/banking-api/internal/domain"
	portout "github.com/reznikrn/banking-api/internal/ports/out"
)

type transactionRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Transaction
}

func NewTransactionRepository() portout.TransactionRepository {
	return &transactionRepository{
		data: make(map[string]*domain.Transaction),
	}
}

func (r *transactionRepository) Save(_ context.Context, tx *domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[tx.ID] = tx
	return nil
}

func (r *transactionRepository) FindByID(_ context.Context, id string) (*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, ok := r.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return tx, nil
}

func (r *transactionRepository) FindAll(_ context.Context) ([]*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Transaction, 0, len(r.data))
	for _, tx := range r.data {
		result = append(result, tx)
	}
	return result, nil
}
