package out

import (
	"context"

	"github.com/reznikrn/banking-api/internal/domain"
)

type TransactionRepository interface {
	Save(ctx context.Context, tx *domain.Transaction) error
	FindByID(ctx context.Context, id string) (*domain.Transaction, error)
	FindAll(ctx context.Context) ([]*domain.Transaction, error)
}
