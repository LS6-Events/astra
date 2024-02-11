package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInlineFunctionsAreMapped(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	paths := testAstra.Path("paths")

	// GET /inline
	require.True(t, paths.Exists("/inline", "get"))
	require.Equal(t, "#/components/schemas/petstore.Pet", paths.Path("/inline.get.responses.200.content.application/json.schema.$ref").Data().(string))

	// POST /inline
	require.True(t, paths.Exists("/inline", "post"))
	require.Equal(t, "#/components/schemas/petstore.PetDTO", paths.Path("/inline.post.requestBody.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/petstore.Pet", paths.Path("/inline.post.responses.200.content.application/json.schema.$ref").Data().(string))

	// GET /inline/{param}
	require.True(t, paths.Exists("/inline/{param}", "get"))
	// Path parameter
	require.Equal(t, "param", paths.Path("/inline/{param}.get.parameters.1.name").Data().(string))
	require.Equal(t, "path", paths.Path("/inline/{param}.get.parameters.1.in").Data().(string))
	require.True(t, paths.Path("/inline/{param}.get.parameters.1.required").Data().(bool))
	require.Equal(t, "string", paths.Path("/inline/{param}.get.parameters.1.schema.type").Data().(string))
	// Query Parameter
	require.Equal(t, "query", paths.Path("/inline/{param}.get.parameters.0.in").Data().(string))
	require.Equal(t, "name", paths.Path("/inline/{param}.get.parameters.0.name").Data().(string))
	require.Equal(t, "string", paths.Path("/inline/{param}.get.parameters.0.schema.type").Data().(string))
	// Response
	require.Equal(t, "#/components/schemas/petstore.Pet", paths.Path("/inline/{param}.get.responses.200.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/inline/{param}.get.responses.400.content.application/json.schema.$ref").Data().(string))
}
