package service

import (
	"errors"
	"testing"

	"github.com/reznikrn/banking-api/internal/domain"
	portin "github.com/reznikrn/banking-api/internal/ports/in"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateInput_ValidTransfer(t *testing.T) {
	input := portin.CreateTransactionInput{
		FromAccount: "ACC-12345",
		ToAccount:   "ACC-67890",
		Amount:      100.50,
		Currency:    "USD",
		Type:        "transfer",
	}

	err := validateCreateInput(input)
	require.NoError(t, err)
}

func TestValidateCreateInput_NegativeAmount(t *testing.T) {
	input := portin.CreateTransactionInput{
		FromAccount: "ACC-12345",
		ToAccount:   "ACC-67890",
		Amount:      -1,
		Currency:    "USD",
		Type:        "transfer",
	}

	err := validateCreateInput(input)
	require.Error(t, err)

	var valErr *domain.ErrValidation
	require.True(t, errors.As(err, &valErr))
	require.NotEmpty(t, valErr.Details)
	assert.Equal(t, "amount", valErr.Details[0].Field)
}

func TestValidateCreateInput_MaxTwoDecimals(t *testing.T) {
	input := portin.CreateTransactionInput{
		FromAccount: "ACC-12345",
		ToAccount:   "ACC-67890",
		Amount:      1.999,
		Currency:    "USD",
		Type:        "transfer",
	}

	err := validateCreateInput(input)
	require.Error(t, err)
}

func TestValidateCreateInput_TransferRequiresValidToAccount(t *testing.T) {
	input := portin.CreateTransactionInput{
		FromAccount: "ACC-12345",
		ToAccount:   "INVALID",
		Amount:      10,
		Currency:    "USD",
		Type:        "transfer",
	}

	err := validateCreateInput(input)
	require.Error(t, err)
}

func TestValidateCreateInput_InvalidCurrencyAndType(t *testing.T) {
	input := portin.CreateTransactionInput{
		FromAccount: "ACC-12345",
		ToAccount:   "ACC-67890",
		Amount:      10,
		Currency:    "FAKE",
		Type:        "unknown",
	}

	err := validateCreateInput(input)
	require.Error(t, err)

	var valErr *domain.ErrValidation
	require.True(t, errors.As(err, &valErr))
	assert.GreaterOrEqual(t, len(valErr.Details), 2)
}
