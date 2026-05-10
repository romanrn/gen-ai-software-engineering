package main

import (
	"log"
	"os"
	handler "support-tickets/internal/adapters/in/http"
	"support-tickets/internal/adapters/out/memory"
	"support-tickets/internal/service"

	_ "support-tickets/docs"
)

// @title Intelligent Customer Support System
// @version 1.0
// @description A customer support ticket management system with auto-classification
// @basePath /
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize repository
	repo := memory.New()

	// Initialize service
	svc := service.NewTicketService(repo)

	// Initialize handlers
	hdlr := handler.New(svc)

	// Setup and run server
	app := setupServer(hdlr)

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
