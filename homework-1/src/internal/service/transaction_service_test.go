package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reznikrn/banking-api/internal/adapters/out/memory"
	"github.com/reznikrn/banking-api/internal/domain"
	portin "github.com/reznikrn/banking-api/internal/ports/in"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction_SetsDerivedFields(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := NewTransactionService(repo)

	input := portin.CreateTransactionInput{
		FromAccount: "ACC-12345",
		ToAccount:   "ACC-67890",
		Amount:      25.5,
		Currency:    "usd",
		Type:        "TRANSFER",
	}

	tx, err := svc.CreateTransaction(context.Background(), input)
	require.NoError(t, err)

	require.NotEmpty(t, tx.ID)
	_, err = uuid.Parse(tx.ID)
	require.NoError(t, err)
	assert.Equal(t, "USD", tx.Currency)
	assert.Equal(t, domain.TypeTransfer, tx.Type)
	assert.Equal(t, domain.StatusCompleted, tx.Status)
	assert.False(t, tx.Timestamp.IsZero())
}

func TestGetTransaction_NotFound(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := NewTransactionService(repo)

	_, err := svc.GetTransaction(context.Background(), "missing")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestListTransactions_FilterByAccountTypeAndDate(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := NewTransactionService(repo)

	tx1 := &domain.Transaction{
		ID:          "tx-1",
		FromAccount: "ACC-11111",
		ToAccount:   "ACC-22222",
		Amount:      10,
		Currency:    "USD",
		Type:        domain.TypeTransfer,
		Timestamp:   time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC),
		Status:      domain.StatusCompleted,
	}
	tx2 := &domain.Transaction{
		ID:          "tx-2",
		FromAccount: "ACC-33333",
		ToAccount:   "ACC-22222",
		Amount:      20,
		Currency:    "USD",
		Type:        domain.TypeDeposit,
		Timestamp:   time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC),
		Status:      domain.StatusCompleted,
	}

	require.NoError(t, repo.Save(context.Background(), tx1))
	require.NoError(t, repo.Save(context.Background(), tx2))

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	result, err := svc.ListTransactions(context.Background(), portin.TransactionFilter{
		AccountID: "ACC-11111",
		Type:      "transfer",
		From:      &from,
		To:        &to,
	})
	require.NoError(t, err)

	require.Len(t, result, 1)
	assert.Equal(t, "tx-1", result[0].ID)
}

func TestGetAccountBalance_ComputesFromTransactionTypes(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := NewTransactionService(repo)
	ctx := context.Background()

	fixtures := []*domain.Transaction{
		{
			ID:          "d1",
			FromAccount: "ACC-12345",
			ToAccount:   "ACC-12345",
			Amount:      100,
			Currency:    "USD",
			Type:        domain.TypeDeposit,
			Timestamp:   time.Now().UTC(),
			Status:      domain.StatusCompleted,
		},
		{
			ID:          "w1",
			FromAccount: "ACC-12345",
			ToAccount:   "ACC-12345",
			Amount:      20,
			Currency:    "USD",
			Type:        domain.TypeWithdrawal,
			Timestamp:   time.Now().UTC(),
			Status:      domain.StatusCompleted,
		},
		{
			ID:          "t1",
			FromAccount: "ACC-12345",
			ToAccount:   "ACC-67890",
			Amount:      15,
			Currency:    "USD",
			Type:        domain.TypeTransfer,
			Timestamp:   time.Now().UTC(),
			Status:      domain.StatusCompleted,
		},
		{
			ID:          "pending",
			FromAccount: "ACC-12345",
			ToAccount:   "ACC-12345",
			Amount:      999,
			Currency:    "USD",
			Type:        domain.TypeDeposit,
			Timestamp:   time.Now().UTC(),
			Status:      domain.StatusPending,
		},
	}

	for _, tx := range fixtures {
		require.NoError(t, repo.Save(ctx, tx))
	}

	balance, err := svc.GetAccountBalance(ctx, "ACC-12345")
	require.NoError(t, err)

	want := 65.0
	assert.Equal(t, want, balance.Balance)
	assert.Equal(t, "USD", balance.Currency)
}

func TestGetAccountBalance_NotFound(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := NewTransactionService(repo)

	_, err := svc.GetAccountBalance(context.Background(), "ACC-00000")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestGetAccountSummary_ComputesTotalsAndLastDate(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := NewTransactionService(repo)
	ctx := context.Background()

	t1 := time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 3, 11, 10, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	fixtures := []*domain.Transaction{
		{
			ID:          "d1",
			FromAccount: "ACC-12345",
			ToAccount:   "ACC-12345",
			Amount:      50,
			Currency:    "USD",
			Type:        domain.TypeDeposit,
			Timestamp:   t1,
			Status:      domain.StatusCompleted,
		},
		{
			ID:          "w1",
			FromAccount: "ACC-12345",
			ToAccount:   "ACC-12345",
			Amount:      20,
			Currency:    "USD",
			Type:        domain.TypeWithdrawal,
			Timestamp:   t2,
			Status:      domain.StatusCompleted,
		},
		{
			ID:          "t1",
			FromAccount: "ACC-12345",
			ToAccount:   "ACC-67890",
			Amount:      10,
			Currency:    "USD",
			Type:        domain.TypeTransfer,
			Timestamp:   t3,
			Status:      domain.StatusCompleted,
		},
	}

	for _, tx := range fixtures {
		require.NoError(t, repo.Save(ctx, tx))
	}

	summary, err := svc.GetAccountSummary(ctx, "ACC-12345")
	require.NoError(t, err)

	assert.Equal(t, 3, summary.TransactionCount)
	assert.Equal(t, 50.0, summary.TotalDeposits)
	assert.Equal(t, 30.0, summary.TotalWithdrawals)
	assert.True(t, summary.LastTransactionDate.Equal(t3))
}

func TestGetAccountSummary_NotFound(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := NewTransactionService(repo)

	_, err := svc.GetAccountSummary(context.Background(), "ACC-00000")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
