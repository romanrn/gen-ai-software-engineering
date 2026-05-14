package main

import (
	"net/http/httptest"
	handler "support-tickets/internal/adapters/in/http"
	"support-tickets/internal/adapters/out/memory"
	"support-tickets/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupServer_RegistersRoutes(t *testing.T) {
	repo := memory.New()
	svc := service.NewTicketService(repo)
	hdlr := handler.New(svc)

	app := setupServer(hdlr)
	require.NotNil(t, app)

	// Verify routes are registered by checking the route count
	routes := app.GetRoutes()
	assert.True(t, len(routes) > 0)

	// Confirm expected paths are registered
	paths := make(map[string]bool)
	for _, r := range routes {
		paths[r.Path] = true
	}
	assert.True(t, paths["/health"])
	assert.True(t, paths["/tickets"])
	assert.True(t, paths["/tickets/import"])
}

func TestSetupServer_HealthEndpoint(t *testing.T) {
	repo := memory.New()
	svc := service.NewTicketService(repo)
	hdlr := handler.New(svc)

	app := setupServer(hdlr)

	// Hit health endpoint directly via app.Test
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
