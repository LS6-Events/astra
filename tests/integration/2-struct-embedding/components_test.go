package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmbeddedStructsAllOf(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	components := testAstra.Path("components")

	schemas := components.Path("schemas")

	// Cat
	require.True(t, schemas.Exists("2-struct-embedding.Cat"))
	require.Equal(t, "object", schemas.Search("2-struct-embedding.Cat", "type").Data().(string))
	// Does reference Pet
	require.Equal(t, "#/components/schemas/petstore.Pet", schemas.Search("2-struct-embedding.Cat", "allOf", "0", "$ref").Data().(string))
	// Breed
	require.True(t, schemas.Exists("2-struct-embedding.Cat", "allOf", "1", "properties", "breed"))
	require.Equal(t, "string", schemas.Search("2-struct-embedding.Cat", "allOf", "1", "properties", "breed", "type").Data().(string))
	// IsIndependent
	require.True(t, schemas.Exists("2-struct-embedding.Cat", "allOf", "1", "properties", "isIndependent"))
	require.Equal(t, "boolean", schemas.Search("2-struct-embedding.Cat", "allOf", "1", "properties", "isIndependent", "type").Data().(string))

	// Dog
	require.True(t, schemas.Exists("2-struct-embedding.Dog"))
	require.Equal(t, "object", schemas.Search("2-struct-embedding.Dog", "type").Data().(string))
	// Does reference Pet
	require.Equal(t, "#/components/schemas/petstore.Pet", schemas.Search("2-struct-embedding.Dog", "allOf", "0", "$ref").Data().(string))
	// Breed
	require.True(t, schemas.Exists("2-struct-embedding.Dog", "allOf", "1", "properties", "breed"))
	require.Equal(t, "string", schemas.Search("2-struct-embedding.Dog", "allOf", "1", "properties", "breed", "type").Data().(string))
	// IsTrained
	require.True(t, schemas.Exists("2-struct-embedding.Dog", "allOf", "1", "properties", "isTrained"))
	require.Equal(t, "boolean", schemas.Search("2-struct-embedding.Dog", "allOf", "1", "properties", "isTrained", "type").Data().(string))
}
