package astra

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWithCustomLogger(t *testing.T) {
	service := &Service{}

	customLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	require.NotEqual(t, customLogger, service.Log)

	WithCustomLogger(customLogger)(service)

	require.Equal(t, customLogger, service.Log)
}

func TestWithCustomLogLevel(t *testing.T) {
	service := &Service{}

	t.Run("default", func(t *testing.T) {
		require.NotEqual(t, zerolog.InfoLevel, service.Log.GetLevel())

		WithCustomLogLevel("")(service)

		require.Equal(t, zerolog.InfoLevel, service.Log.GetLevel())
	})

	t.Run("trace", func(t *testing.T) {
		require.NotEqual(t, zerolog.TraceLevel, service.Log.GetLevel())

		WithCustomLogLevel("trace")(service)

		require.Equal(t, zerolog.TraceLevel, service.Log.GetLevel())
	})

	t.Run("debug", func(t *testing.T) {
		require.NotEqual(t, zerolog.DebugLevel, service.Log.GetLevel())

		WithCustomLogLevel("debug")(service)

		require.Equal(t, zerolog.DebugLevel, service.Log.GetLevel())
	})

	t.Run("info", func(t *testing.T) {
		require.NotEqual(t, zerolog.InfoLevel, service.Log.GetLevel())

		WithCustomLogLevel("info")(service)

		require.Equal(t, zerolog.InfoLevel, service.Log.GetLevel())
	})

	t.Run("warn", func(t *testing.T) {
		require.NotEqual(t, zerolog.WarnLevel, service.Log.GetLevel())

		WithCustomLogLevel("warn")(service)

		require.Equal(t, zerolog.WarnLevel, service.Log.GetLevel())
	})

	t.Run("error", func(t *testing.T) {
		require.NotEqual(t, zerolog.ErrorLevel, service.Log.GetLevel())

		WithCustomLogLevel("error")(service)

		require.Equal(t, zerolog.ErrorLevel, service.Log.GetLevel())
	})

	t.Run("fatal", func(t *testing.T) {
		require.NotEqual(t, zerolog.FatalLevel, service.Log.GetLevel())

		WithCustomLogLevel("fatal")(service)

		require.Equal(t, zerolog.FatalLevel, service.Log.GetLevel())
	})

	t.Run("panic", func(t *testing.T) {
		require.NotEqual(t, zerolog.PanicLevel, service.Log.GetLevel())

		WithCustomLogLevel("panic")(service)

		require.Equal(t, zerolog.PanicLevel, service.Log.GetLevel())
	})
}
