package in

import (
	"context"
	"time"

	"github.com/reznikrn/banking-api/internal/domain"
)

type TransactionFilter struct {
	AccountID string
	Type      string
	From      *time.Time
	To        *time.Time
}

type CreateTransactionInput struct {
	FromAccount string
	ToAccount   string
	Amount      float64
	Currency    string
	Type        string
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (*domain.Transaction, error)
	GetTransaction(ctx context.Context, id string) (*domain.Transaction, error)
	ListTransactions(ctx context.Context, filter TransactionFilter) ([]domain.Transaction, error)
	GetAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error)
	GetAccountSummary(ctx context.Context, accountID string) (*domain.AccountSummary, error)
}
