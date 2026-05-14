package errorhandler

import (
	"errors"
	"net/http"
	"support-tickets/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error     string                        `json:"error"`
	RequestID string                        `json:"request_id"`
	Details   []domain.ValidationDetail `json:"details,omitempty"`
}

func Handler(c *fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	var details []domain.ValidationDetail

	requestID := c.Locals("request_id")

	var validationErr *domain.ErrValidation
	var notFoundErr *domain.ErrNotFound
	var invalidInputErr *domain.ErrInvalidInput

	if errors.As(err, &validationErr) {
		code = http.StatusBadRequest
		message = "Validation Error"
		details = validationErr.Details
	} else if errors.As(err, &notFoundErr) {
		code = http.StatusNotFound
		message = notFoundErr.Error()
	} else if errors.As(err, &invalidInputErr) {
		code = http.StatusBadRequest
		message = invalidInputErr.Error()
	}

	return c.Status(code).JSON(ErrorResponse{
		Error:     message,
		RequestID: requestID.(string),
		Details:   details,
	})
}
