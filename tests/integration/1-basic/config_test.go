package integration

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	t.Run("OpenAPI Version", func(t *testing.T) {
		require.Equal(t, "3.0.0", testAstra.Path("openapi").Data().(string))
	})

	t.Run("OpenAPI Info", func(t *testing.T) {
		require.Equal(t, "Generated by astra", testAstra.Path("info.description").Data().(string))
	})

	t.Run("OpenAPI Servers", func(t *testing.T) {
		require.Equal(t, "http://localhost:8000", testAstra.Path("servers.0.url").Data().(string))
	})
}

func TestCustomConfig(t *testing.T) {
	r := setupRouter()

	config := &astra.Config{
		Title:       "Test Title",
		Description: "Test Description",
		Version:     "0.0.1",
		Contact: astra.Contact{
			Name:  "John Doe",
			URL:   "https://www.google.com",
			Email: "john@doe.com",
		},
		License: astra.License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
		Secure:   false,
		Host:     "localhost",
		BasePath: "/base-path",
		Port:     8000,
	}

	testAstra, err := helpers.SetupTestAstra(t, r, config)
	require.NoError(t, err)

	t.Run("OpenAPI Version", func(t *testing.T) {
		require.Equal(t, "3.0.0", testAstra.Path("openapi").Data().(string))
	})

	t.Run("OpenAPI Info", func(t *testing.T) {
		require.Equal(t, "Test Title", testAstra.Path("info.title").Data().(string))
		require.Equal(t, "0.0.1", testAstra.Path("info.version").Data().(string))
		require.Equal(t, "John Doe", testAstra.Path("info.contact.name").Data().(string))
		require.Equal(t, "https://www.google.com", testAstra.Path("info.contact.url").Data().(string))
		require.Equal(t, "Test Description", testAstra.Path("info.description").Data().(string))
	})

	t.Run("OpenAPI Servers", func(t *testing.T) {
		require.Equal(t, "http://localhost:8000/base-path", testAstra.Path("servers.0.url").Data().(string))
	})
}
