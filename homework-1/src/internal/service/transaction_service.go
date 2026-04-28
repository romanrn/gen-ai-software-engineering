package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reznikrn/banking-api/internal/domain"
	portin "github.com/reznikrn/banking-api/internal/ports/in"
	portout "github.com/reznikrn/banking-api/internal/ports/out"
)

type transactionService struct {
	repo portout.TransactionRepository
}

func NewTransactionService(repo portout.TransactionRepository) portin.TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(ctx context.Context, input portin.CreateTransactionInput) (*domain.Transaction, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}

	tx := &domain.Transaction{
		ID:          uuid.New().String(),
		FromAccount: input.FromAccount,
		ToAccount:   input.ToAccount,
		Amount:      input.Amount,
		Currency:    strings.ToUpper(input.Currency),
		Type:        domain.TransactionType(strings.ToLower(input.Type)),
		Timestamp:   time.Now().UTC(),
		Status:      domain.StatusCompleted,
	}

	if err := s.repo.Save(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *transactionService) ListTransactions(ctx context.Context, filter portin.TransactionFilter) ([]domain.Transaction, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Transaction, 0)
	for _, tx := range all {
		if matchesFilter(tx, filter) {
			result = append(result, *tx)
		}
	}
	return result, nil
}

func (s *transactionService) GetAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var balance float64
	var currency string
	found := false

	for _, tx := range all {
		if tx.Status != domain.StatusCompleted {
			continue
		}
		if tx.FromAccount != accountID && tx.ToAccount != accountID {
			continue
		}
		found = true
		currency = tx.Currency

		switch tx.Type {
		case domain.TypeDeposit:
			if tx.ToAccount == accountID {
				balance += tx.Amount
			}
		case domain.TypeWithdrawal:
			if tx.FromAccount == accountID {
				balance -= tx.Amount
			}
		case domain.TypeTransfer:
			if tx.ToAccount == accountID {
				balance += tx.Amount
			}
			if tx.FromAccount == accountID {
				balance -= tx.Amount
			}
		}
	}

	if !found {
		return nil, domain.ErrNotFound
	}

	return &domain.AccountBalance{
		AccountID: accountID,
		Balance:   balance,
		Currency:  currency,
	}, nil
}

func (s *transactionService) GetAccountSummary(ctx context.Context, accountID string) (*domain.AccountSummary, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	summary := &domain.AccountSummary{AccountID: accountID}
	found := false

	for _, tx := range all {
		if tx.FromAccount != accountID && tx.ToAccount != accountID {
			continue
		}
		found = true
		summary.TransactionCount++

		if tx.Timestamp.After(summary.LastTransactionDate) {
			summary.LastTransactionDate = tx.Timestamp
		}

		switch tx.Type {
		case domain.TypeDeposit:
			summary.TotalDeposits += tx.Amount
		case domain.TypeWithdrawal:
			summary.TotalWithdrawals += tx.Amount
		case domain.TypeTransfer:
			if tx.ToAccount == accountID {
				summary.TotalDeposits += tx.Amount
			}
			if tx.FromAccount == accountID {
				summary.TotalWithdrawals += tx.Amount
			}
		}
	}

	if !found {
		return nil, domain.ErrNotFound
	}

	return summary, nil
}

func matchesFilter(tx *domain.Transaction, f portin.TransactionFilter) bool {
	if f.AccountID != "" && tx.FromAccount != f.AccountID && tx.ToAccount != f.AccountID {
		return false
	}
	if f.Type != "" && string(tx.Type) != strings.ToLower(f.Type) {
		return false
	}
	if f.From != nil && tx.Timestamp.Before(*f.From) {
		return false
	}
	if f.To != nil && tx.Timestamp.After(*f.To) {
		return false
	}
	return true
}
