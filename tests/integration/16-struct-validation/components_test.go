package integration

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

	// Pet
	require.True(t, schemas.Exists("petstore.Pet"))
	require.ElementsMatch(t, []interface{}{"id", "name"}, schemas.Search("petstore.Pet", "required").Data().([]interface{}))
	require.Equal(t, map[string]interface{}{"id": map[string]interface{}{"allOf": []interface{}{map[string]interface{}{"format": "uuid"}}, "type": "string"}, "name": map[string]interface{}{"type": "string"}, "photoUrls": map[string]interface{}{"allOf": []interface{}{map[string]interface{}{}}, "items": map[string]interface{}{"type": "string", "allOf": []interface{}{map[string]interface{}{"format": "uri"}}}, "type": "array"}, "status": map[string]interface{}{"type": "string"}, "tags": map[string]interface{}{"items": map[string]interface{}{"$ref": "#/components/schemas/petstore.Tag"}, "type": "array"}}, schemas.Search("petstore.Pet", "properties").Data().(map[string]interface{}))

	// PetDTO
	require.True(t, schemas.Exists("petstore.PetDTO"))
	require.ElementsMatch(t, []interface{}{"name"}, schemas.Search("petstore.PetDTO", "required").Data().([]interface{}))
	require.Equal(t, map[string]interface{}{"name": map[string]interface{}{"type": "string"}, "photoUrls": map[string]interface{}{"allOf": []interface{}{map[string]interface{}{}}, "items": map[string]interface{}{"type": "string", "allOf": []interface{}{map[string]interface{}{"format": "uri"}}}, "type": "array"}, "status": map[string]interface{}{"type": "string"}, "tags": map[string]interface{}{"items": map[string]interface{}{"$ref": "#/components/schemas/petstore.Tag"}, "type": "array"}}, schemas.Search("petstore.PetDTO", "properties").Data().(map[string]interface{}))

	// Tag
	require.True(t, schemas.Exists("petstore.Tag"))
	require.ElementsMatch(t, []interface{}{"id", "name"}, schemas.Search("petstore.Tag", "required").Data().([]interface{}))
	require.Equal(t, map[string]interface{}{"id": map[string]interface{}{"allOf": []interface{}{map[string]interface{}{"format": "uuid"}}, "type": "string"}, "name": map[string]interface{}{"type": "string"}}, schemas.Search("petstore.Tag", "properties").Data().(map[string]interface{}))
}
