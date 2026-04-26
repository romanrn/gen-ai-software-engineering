package domain

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal error")
)

type ValidationDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrValidation struct {
	Details []ValidationDetail
}

func (e *ErrValidation) Error() string {
	return "validation failed"
}
