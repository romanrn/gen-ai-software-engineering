package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	handler "github.com/reznikrn/banking-api/internal/adapters/in/http"
	"github.com/reznikrn/banking-api/pkg/errorhandler"
	"github.com/reznikrn/banking-api/pkg/middleware"
)

// healthHandler godoc
// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} handler.HealthResponse
// @Router /health [get]
func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func newServer(txHandler *handler.TransactionHandler, accHandler *handler.AccountHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.Handler,
	})

	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())

	app.Get("/health", healthHandler)

	app.Post("/transactions", txHandler.Create)
	app.Get("/transactions", txHandler.List)
	app.Get("/transactions/:id", txHandler.GetByID)

	app.Get("/accounts/:accountId/balance", accHandler.GetBalance)
	app.Get("/accounts/:accountId/summary", accHandler.GetSummary)

	return app
}
