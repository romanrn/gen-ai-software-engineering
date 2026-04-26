package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrValidation_Error(t *testing.T) {
	err := &ErrValidation{}
	assert.Equal(t, "validation failed", err.Error())
}
