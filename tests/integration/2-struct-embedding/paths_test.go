package integration

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	paths := testAstra.Path("paths")

	// GET /cats/{id}
	require.True(t, paths.Exists("/cats/{id}", "get"))
	require.Equal(t, "string", paths.Path("/cats/{id}.get.parameters.0.schema.type").Data().(string))
	require.Equal(t, "#/components/schemas/2-struct-embedding.Cat", paths.Path("/cats/{id}.get.responses.200.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/cats/{id}.get.responses.400.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/cats/{id}.get.responses.404.content.application/json.schema.$ref").Data().(string))

	// GET /dogs/{id}
	require.True(t, paths.Exists("/dogs/{id}", "get"))
	require.Equal(t, "string", paths.Path("/dogs/{id}.get.parameters.0.schema.type").Data().(string))
	require.Equal(t, "#/components/schemas/2-struct-embedding.Dog", paths.Path("/dogs/{id}.get.responses.200.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/dogs/{id}.get.responses.400.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/dogs/{id}.get.responses.404.content.application/json.schema.$ref").Data().(string))

}
