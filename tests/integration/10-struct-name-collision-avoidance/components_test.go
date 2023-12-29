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

	schemas := testAstra.Path("components.schemas")

	require.NotNil(t, schemas)

	require.True(t, schemas.Exists("10-struct-name-collision-avoidance.types.TestType"))
	require.Equal(t, "object", schemas.Search("10-struct-name-collision-avoidance.types.TestType", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("10-struct-name-collision-avoidance.types.TestType", "properties", "topLevelType", "type").Data().(string))

	require.True(t, schemas.Exists("nested.types.TestType"))
	require.Equal(t, "object", schemas.Search("nested.types.TestType", "type").Data().(string))
	require.Equal(t, "string", schemas.Search("nested.types.TestType", "properties", "nestedField", "type").Data().(string))
}
