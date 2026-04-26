package errorhandler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/reznikrn/banking-api/internal/domain"
	"github.com/reznikrn/banking-api/pkg/middleware"
)

func Handler(c *fiber.Ctx, err error) error {
	requestID, _ := c.Locals(middleware.RequestIDKey).(string)

	var valErr *domain.ErrValidation
	if errors.As(err, &valErr) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":      "Validation failed",
			"request_id": requestID,
			"details":    valErr.Details,
		})
	}

	if errors.Is(err, domain.ErrNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":      "Not found",
			"request_id": requestID,
			"details":    []any{},
		})
	}

	if errors.Is(err, domain.ErrInvalidInput) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": requestID,
			"details":    []any{},
		})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"error":      fiberErr.Message,
			"request_id": requestID,
			"details":    []any{},
		})
	}

	slog.Error("unhandled error", "request_id", requestID, "error", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":      "Internal server error",
		"request_id": requestID,
		"details":    []any{},
	})
}
