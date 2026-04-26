package memory

import (
	"context"
	"testing"

	"github.com/reznikrn/banking-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepository_SaveFindByIDFindAll(t *testing.T) {
	repo := NewTransactionRepository()
	ctx := context.Background()

	tx := &domain.Transaction{ID: "tx-1", FromAccount: "ACC-12345", ToAccount: "ACC-67890"}

	require.NoError(t, repo.Save(ctx, tx))

	found, err := repo.FindByID(ctx, "tx-1")
	require.NoError(t, err)
	assert.Equal(t, tx, found)

	all, err := repo.FindAll(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, "tx-1", all[0].ID)
}

func TestTransactionRepository_FindByIDNotFound(t *testing.T) {
	repo := NewTransactionRepository()

	_, err := repo.FindByID(context.Background(), "missing")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
