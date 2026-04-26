package errorhandler

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/reznikrn/banking-api/internal/domain"
	"github.com/reznikrn/banking-api/pkg/middleware"
	"github.com/stretchr/testify/require"
)

func TestHandler_MapsErrorsToStatusCodes(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: Handler})
	app.Use(middleware.RequestID())

	app.Get("/validation", func(c *fiber.Ctx) error {
		return &domain.ErrValidation{Details: []domain.ValidationDetail{{Field: "amount", Message: "bad"}}}
	})
	app.Get("/notfound", func(c *fiber.Ctx) error { return domain.ErrNotFound })
	app.Get("/invalid", func(c *fiber.Ctx) error { return domain.ErrInvalidInput })
	app.Get("/fiber", func(c *fiber.Ctx) error { return fiber.NewError(fiber.StatusTeapot, "teapot") })
	app.Get("/internal", func(c *fiber.Ctx) error { return errors.New("boom") })

	cases := []struct {
		path string
		code int
	}{
		{path: "/validation", code: fiber.StatusUnprocessableEntity},
		{path: "/notfound", code: fiber.StatusNotFound},
		{path: "/invalid", code: fiber.StatusBadRequest},
		{path: "/fiber", code: fiber.StatusTeapot},
		{path: "/internal", code: fiber.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tc.code, resp.StatusCode)
			require.NotEmpty(t, resp.Header.Get("X-Request-ID"))
		})
	}
}
