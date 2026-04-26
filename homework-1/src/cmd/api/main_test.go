package main

import (
	"errors"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_DefaultPort(t *testing.T) {
	require.NoError(t, os.Unsetenv("PORT"))

	called := false
	err := run(func(app *fiber.App, addr string) error {
		called = true
		assert.Equal(t, ":8080", addr)
		return nil
	})

	require.NoError(t, err)
	assert.True(t, called)
}

func TestRun_UsesEnvPortAndPropagatesListenError(t *testing.T) {
	require.NoError(t, os.Setenv("PORT", "9090"))
	defer os.Unsetenv("PORT")

	expected := errors.New("listen failed")
	err := run(func(app *fiber.App, addr string) error {
		assert.Equal(t, ":9090", addr)
		return expected
	})

	assert.ErrorIs(t, err, expected)
}
