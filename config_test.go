package astra

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	t.Run("valid full config", func(t *testing.T) {
		config := &Config{
			Title:       "test",
			Description: "test description",
			Version:     "0.0.1",
			Contact: Contact{
				Name:  "John Doe",
				URL:   "https://www.google.com",
				Email: "john@doe.com",
			},
			License: License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Secure:   false,
			Host:     "localhost",
			BasePath: "/base-path",
			Port:     8000,
		}

		err := config.Validate()
		require.NoError(t, err)
	})

	t.Run("valid config with defaults", func(t *testing.T) {
		config := &Config{
			Port: 8000,
		}

		err := config.Validate()
		require.NoError(t, err)

		require.Equal(t, "localhost", config.Host)
		require.Equal(t, "/", config.BasePath)
	})

	t.Run("invalid config", func(t *testing.T) {
		config := &Config{}

		err := config.Validate()
		require.Error(t, err)
		require.Equal(t, ErrConfigPortRequired, err)
	})
}

func TestService_SetConfig(t *testing.T) {
	service := &Service{}
	config := &Config{
		Port: 8000,
	}

	service.SetConfig(config)

	require.Equal(t, config, service.Config)
	require.Equal(t, config.Port, service.Config.Port)
}

func TestWithConfig(t *testing.T) {
	service := &Service{}
	config := &Config{
		Port: 8000,
	}

	WithConfig(config)(service)

	require.Equal(t, config, service.Config)
	require.Equal(t, config.Port, service.Config.Port)
}
