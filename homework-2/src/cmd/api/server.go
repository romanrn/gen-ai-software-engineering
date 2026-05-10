package main

import (
	handler "support-tickets/internal/adapters/in/http"
	"support-tickets/pkg/errorhandler"
	"support-tickets/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

func setupServer(hdlr *handler.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.Handler,
	})

	// Middleware
	app.Use(middleware.RequestID)
	app.Use(fiberlogger.New())
	app.Use(recover.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes
	api := app.Group("/tickets")
	{
		api.Post("", hdlr.CreateTicket)
		api.Post("/import", hdlr.BulkImport)
		api.Get("", hdlr.ListTickets)
		api.Get("/:id", hdlr.GetTicket)
		api.Put("/:id", hdlr.UpdateTicket)
		api.Delete("/:id", hdlr.DeleteTicket)
		api.Post("/:id/auto-classify", hdlr.AutoClassify)
	}

	return app
}
