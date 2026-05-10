package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"support-tickets/internal/adapters/in/http"
	"support-tickets/internal/adapters/out/memory"
	"support-tickets/internal/service"
	"support-tickets/pkg/errorhandler"
	"support-tickets/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func newApp() *fiber.App {
	repo := memory.New()
	svc := service.NewTicketService(repo)
	hdlr := handler.New(svc)

	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.Handler})
	app.Use(middleware.RequestID)

	tickets := app.Group("/tickets")
	tickets.Post("", hdlr.CreateTicket)
	tickets.Post("/import", hdlr.BulkImport)
	tickets.Get("", hdlr.ListTickets)
	tickets.Get("/:id", hdlr.GetTicket)
	tickets.Put("/:id", hdlr.UpdateTicket)
	tickets.Delete("/:id", hdlr.DeleteTicket)
	tickets.Post("/:id/auto-classify", hdlr.AutoClassify)

	return app
}

func jsonRequest(method, url string, body any) *http.Request {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, url, r)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func decodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

// fixturesDir returns the absolute path to the fixtures directory.
func fixturesDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "fixtures")
}
