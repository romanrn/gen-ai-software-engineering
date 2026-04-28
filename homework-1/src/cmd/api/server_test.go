package main

import (
	"net/http/httptest"
	"testing"

	handler "github.com/reznikrn/banking-api/internal/adapters/in/http"
	"github.com/reznikrn/banking-api/internal/adapters/out/memory"
	"github.com/reznikrn/banking-api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestNewServer_HealthRoute(t *testing.T) {
	repo := memory.NewTransactionRepository()
	svc := service.NewTransactionService(repo)
	txHandler := handler.NewTransactionHandler(svc)
	accHandler := handler.NewAccountHandler(svc)

	app := newServer(txHandler, accHandler)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}
