package domain

import "fmt"

type ValidationDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrValidation struct {
	Details []ValidationDetail
}

func (e *ErrValidation) Error() string {
	if len(e.Details) == 0 {
		return "validation error"
	}
	return fmt.Sprintf("validation error: %d field(s) invalid", len(e.Details))
}

type ErrNotFound struct {
	ID string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("ticket not found: %s", e.ID)
}

type ErrInvalidInput struct {
	Message string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Message)
}
