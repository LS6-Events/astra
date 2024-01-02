package astra

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("creates a new service", func(t *testing.T) {
		service := New()

		require.NotNil(t, service)
	})

	t.Run("populates service with default values", func(t *testing.T) {
		service := New()

		require.Equal(t, log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.InfoLevel), service.Log)
		require.Equal(t, make([]Route, 0), service.Routes)
		require.Equal(t, make([]Field, 0), service.Components)
		require.Equal(t, make(map[string]TypeFormat), service.CustomTypeMapping)
	})

	t.Run("populates service with option", func(t *testing.T) {
		var optionCalled bool
		testOption := func(service *Service) {
			optionCalled = true
			service.Config = &Config{
				Title: "test ran",
			}
		}

		service := New(testOption)

		require.True(t, optionCalled)
		require.Equal(t, "test ran", service.Config.Title)
	})

	t.Run("populates service with many options", func(t *testing.T) {
		var option1Called, option2Called bool
		testOption1 := func(service *Service) {
			option1Called = true
			service.Config = &Config{
				Title: "test 1 ran",
			}
		}
		testOption2 := func(service *Service) {
			option2Called = true
			service.Config.Description = "test 2 ran"
		}

		service := New(testOption1, testOption2)

		require.True(t, option1Called)
		require.True(t, option2Called)
		require.Equal(t, "test 1 ran", service.Config.Title)
		require.Equal(t, "test 2 ran", service.Config.Description)
	})
}
