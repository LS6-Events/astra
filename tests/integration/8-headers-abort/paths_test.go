package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeadersAndAbort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	require.NotNil(t, testAstra)

	paths := testAstra.Path("paths")

	t.Run("Headers", func(t *testing.T) {
		// Get
		require.True(t, paths.Exists("/headers", "get", "parameters", "0"))
		require.Equal(t, "X-Test-Header", paths.Path("/headers.get.parameters.0.name").Data().(string))
		require.Equal(t, "header", paths.Path("/headers.get.parameters.0.in").Data().(string))
		require.Equal(t, "string", paths.Path("/headers.get.parameters.0.schema.type").Data().(string))

		// Set
		require.True(t, paths.Exists("/headers", "post", "responses", "200", "headers", "X-Test-Header"))
		require.Equal(t, "string", paths.Path("/headers.post.responses.200.headers.X-Test-Header.schema.type").Data().(string))
	})

	t.Run("Abort", func(t *testing.T) {
		// AbortWithError
		require.True(t, paths.Exists("/abort-with-error", "get", "responses", "400"))

		// AbortWithStatus
		require.True(t, paths.Exists("/abort-with-status", "get", "responses", "401"))

		// AbortWithStatusJSON
		require.True(t, paths.Exists("/abort-with-status-json", "get", "responses", "402"))
		require.Equal(t, "#/components/schemas/gin.H", paths.Search("/abort-with-status-json", "get", "responses", "402", "content", "application/json", "schema", "$ref").Data().(string))
	})
}
