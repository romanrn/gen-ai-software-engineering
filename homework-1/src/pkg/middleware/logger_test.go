package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLogger_PassesThroughResponse(t *testing.T) {
	app := fiber.New()
	app.Use(Logger())
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusCreated).SendString("ok")
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
