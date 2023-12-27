package petstore

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

	// GET /pets
	require.True(t, paths.Exists("/pets", "get"))
	require.Equal(t, "array", paths.Path("/pets.get.responses.200.content.application/json.schema.type").Data().(string))
	require.Equal(t, "#/components/schemas/petstore.Pet", paths.Path("/pets.get.responses.200.content.application/json.schema.items.$ref").Data().(string))

	// POST /pets
	require.True(t, paths.Exists("/pets", "post"))
	require.Equal(t, "#/components/schemas/petstore.PetDTO", paths.Path("/pets.post.requestBody.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/petstore.PetDTO", paths.Path("/pets.post.responses.200.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/pets.post.responses.400.content.application/json.schema.$ref").Data().(string))

	// GET /pets/{id}
	require.True(t, paths.Exists("/pets/{id}", "get"))
	require.Equal(t, "string", paths.Path("/pets/{id}.get.parameters.0.schema.type").Data().(string))
	require.Equal(t, "id", paths.Path("/pets/{id}.get.parameters.0.name").Data().(string))
	require.Equal(t, "path", paths.Path("/pets/{id}.get.parameters.in").Data().(string))
	require.True(t, paths.Path("/pets/{id}.get.parameters.0.required").Data().(bool))
	require.Equal(t, "string", paths.Path("/pets/{id}.get.parameters.0.schema.type").Data().(string))
	require.Equal(t, "#/components/schemas/petstore.Pet", paths.Path("/pets/{id}.get.responses.200.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/pets/{id}.get.responses.400.content.application/json.schema.$ref").Data().(string))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/pets/{id}.get.responses.404.content.application/json.schema.$ref").Data().(string))

	// DELETE /pets/{id}
	require.True(t, paths.Exists("/pets/{id}", "delete"))
	require.Equal(t, "id", paths.Path("/pets/{id}.get.parameters.0.name").Data().(string))
	require.Equal(t, "path", paths.Path("/pets/{id}.get.parameters.in").Data().(string))
	require.True(t, paths.Path("/pets/{id}.get.parameters.0.required").Data().(bool))
	require.Equal(t, "string", paths.Path("/pets/{id}.get.parameters.0.schema.type").Data().(string))
	require.True(t, paths.Exists("/pets/{id}", "delete", "responses", "200"))
	require.Equal(t, "#/components/schemas/gin.H", paths.Path("/pets/{id}.delete.responses.400.content.application/json.schema.$ref").Data().(string))
}
