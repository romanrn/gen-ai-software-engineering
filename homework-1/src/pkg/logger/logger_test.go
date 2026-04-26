package logger

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetup_HandlesSupportedLogLevels(t *testing.T) {
	levels := []string{"", "DEBUG", "WARN", "ERROR"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			require.NoError(t, os.Setenv("LOG_LEVEL", level))
			Setup()
			require.NotNil(t, slog.Default())
		})
	}
}
