package handler

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/reznikrn/banking-api/internal/domain"
	portin "github.com/reznikrn/banking-api/internal/ports/in"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	createFn    func(ctx context.Context, input portin.CreateTransactionInput) (*domain.Transaction, error)
	getFn       func(ctx context.Context, id string) (*domain.Transaction, error)
	listFn      func(ctx context.Context, filter portin.TransactionFilter) ([]domain.Transaction, error)
	balanceFn   func(ctx context.Context, accountID string) (*domain.AccountBalance, error)
	summaryFn   func(ctx context.Context, accountID string) (*domain.AccountSummary, error)
	lastInput   portin.CreateTransactionInput
	lastFilter  portin.TransactionFilter
	lastID      string
	lastAccount string
}

func (m *mockService) CreateTransaction(ctx context.Context, input portin.CreateTransactionInput) (*domain.Transaction, error) {
	m.lastInput = input
	return m.createFn(ctx, input)
}

func (m *mockService) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	m.lastID = id
	return m.getFn(ctx, id)
}

func (m *mockService) ListTransactions(ctx context.Context, filter portin.TransactionFilter) ([]domain.Transaction, error) {
	m.lastFilter = filter
	return m.listFn(ctx, filter)
}

func (m *mockService) GetAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error) {
	m.lastAccount = accountID
	return m.balanceFn(ctx, accountID)
}

func (m *mockService) GetAccountSummary(ctx context.Context, accountID string) (*domain.AccountSummary, error) {
	m.lastAccount = accountID
	return m.summaryFn(ctx, accountID)
}

func newMockService() *mockService {
	return &mockService{
		createFn: func(ctx context.Context, input portin.CreateTransactionInput) (*domain.Transaction, error) {
			return &domain.Transaction{ID: "tx-1"}, nil
		},
		getFn: func(ctx context.Context, id string) (*domain.Transaction, error) {
			return &domain.Transaction{ID: id}, nil
		},
		listFn: func(ctx context.Context, filter portin.TransactionFilter) ([]domain.Transaction, error) {
			return []domain.Transaction{{ID: "tx-1"}}, nil
		},
		balanceFn: func(ctx context.Context, accountID string) (*domain.AccountBalance, error) {
			return &domain.AccountBalance{AccountID: accountID, Balance: 10, Currency: "USD"}, nil
		},
		summaryFn: func(ctx context.Context, accountID string) (*domain.AccountSummary, error) {
			return &domain.AccountSummary{AccountID: accountID, TransactionCount: 1}, nil
		},
	}
}

func TestTransactionHandler_Create(t *testing.T) {
	mock := newMockService()
	h := NewTransactionHandler(mock)
	app := fiber.New()
	app.Post("/transactions", h.Create)

	req := httptest.NewRequest("POST", "/transactions", strings.NewReader(`{"fromAccount":"ACC-12345","toAccount":"ACC-67890","amount":10.5,"currency":"USD","type":"transfer"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	assert.Equal(t, "ACC-12345", mock.lastInput.FromAccount)
}

func TestTransactionHandler_CreateInvalidBody(t *testing.T) {
	h := NewTransactionHandler(newMockService())
	app := fiber.New()
	app.Post("/transactions", h.Create)

	req := httptest.NewRequest("POST", "/transactions", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTransactionHandler_List(t *testing.T) {
	mock := newMockService()
	h := NewTransactionHandler(mock)
	app := fiber.New()
	app.Get("/transactions", h.List)

	req := httptest.NewRequest("GET", "/transactions?accountId=ACC-12345&type=transfer&from=2026-01-01&to=2026-01-31", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "ACC-12345", mock.lastFilter.AccountID)
	assert.Equal(t, "transfer", mock.lastFilter.Type)
	require.NotNil(t, mock.lastFilter.From)
	require.NotNil(t, mock.lastFilter.To)
	assert.Equal(t, 23, mock.lastFilter.To.Hour())
}

func TestTransactionHandler_ListInvalidDates(t *testing.T) {
	h := NewTransactionHandler(newMockService())
	app := fiber.New()
	app.Get("/transactions", h.List)

	reqFrom := httptest.NewRequest("GET", "/transactions?from=2026-99-99", nil)
	respFrom, err := app.Test(reqFrom)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, respFrom.StatusCode)

	reqTo := httptest.NewRequest("GET", "/transactions?to=bad", nil)
	respTo, err := app.Test(reqTo)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, respTo.StatusCode)
}

func TestTransactionHandler_GetByID(t *testing.T) {
	mock := newMockService()
	h := NewTransactionHandler(mock)
	app := fiber.New()
	app.Get("/transactions/:id", h.GetByID)

	req := httptest.NewRequest("GET", "/transactions/tx-123", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "tx-123", mock.lastID)
}

func TestAccountHandler_Endpoints(t *testing.T) {
	mock := newMockService()
	h := NewAccountHandler(mock)
	app := fiber.New()
	app.Get("/accounts/:accountId/balance", h.GetBalance)
	app.Get("/accounts/:accountId/summary", h.GetSummary)

	reqBal := httptest.NewRequest("GET", "/accounts/ACC-12345/balance", nil)
	respBal, err := app.Test(reqBal)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, respBal.StatusCode)

	reqSum := httptest.NewRequest("GET", "/accounts/ACC-12345/summary", nil)
	respSum, err := app.Test(reqSum)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, respSum.StatusCode)
	assert.Equal(t, "ACC-12345", mock.lastAccount)
}

func TestTransactionHandler_ListResponseIsJSON(t *testing.T) {
	mock := newMockService()
	mock.listFn = func(ctx context.Context, filter portin.TransactionFilter) ([]domain.Transaction, error) {
		return []domain.Transaction{{ID: "tx-json", Timestamp: time.Now().UTC()}}, nil
	}
	h := NewTransactionHandler(mock)
	app := fiber.New()
	app.Get("/transactions", h.List)

	req := httptest.NewRequest("GET", "/transactions", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body []map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Len(t, body, 1)
	assert.Equal(t, "tx-json", body[0]["id"])
}
