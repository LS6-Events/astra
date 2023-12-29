package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchemas(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	components := testAstra.Path("components")

	schemas := components.Path("schemas")

	// gin.H
	require.True(t, schemas.Exists("gin.H"))
	require.Equal(t, "object", schemas.Search("gin.H", "type").Data().(string))

	// petstore.Pet
	require.True(t, schemas.Exists("petstore.Pet"))
	require.Equal(t, "object", schemas.Search("petstore.Pet", "type").Data().(string))
	require.Equal(t, "integer", schemas.Search("petstore.Pet", "properties", "id", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("petstore.Pet", "properties", "name", "type").Data().(string))
	require.Equal(t, "array", schemas.Search("petstore.Pet", "properties", "photoUrls", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("petstore.Pet", "properties", "photoUrls", "items", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("petstore.Pet", "properties", "status", "type").Data().(string))
	require.Equal(t, "array", schemas.Search("petstore.Pet", "properties", "tags", "type").Data().(string))
	require.Equal(t, "#/components/schemas/petstore.Tag", schemas.Search("petstore.Pet", "properties", "tags", "items", "$ref").Data().(string))

	// petstore.PetDTO
	require.True(t, schemas.Exists("petstore.PetDTO"))
	require.Equal(t, "object", schemas.Search("petstore.PetDTO", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("petstore.PetDTO", "properties", "name", "type").Data().(string))
	require.Equal(t, "array", schemas.Search("petstore.PetDTO", "properties", "photoUrls", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("petstore.PetDTO", "properties", "photoUrls", "items", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("petstore.PetDTO", "properties", "status", "type").Data().(string))
	require.Equal(t, "array", schemas.Search("petstore.PetDTO", "properties", "tags", "type").Data().(string))
	require.Equal(t, "#/components/schemas/petstore.Tag", schemas.Search("petstore.PetDTO", "properties", "tags", "items", "$ref").Data().(string))

	// petstore.Tag
	require.True(t, schemas.Exists("petstore.Tag"))
	require.Equal(t, "object", schemas.Search("petstore.Tag", "type").Data().(string))
	require.Equal(t, "integer", schemas.Search("petstore.Tag", "properties", "id", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("petstore.Tag", "properties", "name", "type").Data().(string))
}
