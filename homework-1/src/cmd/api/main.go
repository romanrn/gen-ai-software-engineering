package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	handler "github.com/reznikrn/banking-api/internal/adapters/in/http"
	"github.com/reznikrn/banking-api/internal/adapters/out/memory"
	"github.com/reznikrn/banking-api/internal/service"
	"github.com/reznikrn/banking-api/pkg/logger"
)

// @title Banking Transactions API
// @version 1.0
// @description REST API for creating and querying banking transactions.
// @BasePath /
// @schemes http

func main() {
	if err := run(func(app *fiber.App, addr string) error {
		return app.Listen(addr)
	}); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func run(listen func(app *fiber.App, addr string) error) error {
	logger.Setup()

	repo := memory.NewTransactionRepository()
	svc := service.NewTransactionService(repo)

	txHandler := handler.NewTransactionHandler(svc)
	accHandler := handler.NewAccountHandler(svc)

	app := newServer(txHandler, accHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "port", port)
	return listen(app, ":"+port)
}
