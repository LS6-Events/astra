package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComponentsDocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	components := testAstra.Path("components")

	schemas := components.Path("schemas")

	// gin.H
	require.True(t, schemas.Exists("gin.H", "description"))
	require.Equal(t, "H is a shortcut for map[string]any", schemas.Search("gin.H", "description").Data().(string))

	// petstore.Pet
	require.True(t, schemas.Exists("petstore.Pet", "description"))
	require.Equal(t, "Pet the pet model.", schemas.Search("petstore.Pet", "description").Data().(string))

	// petstore.PetDTO
	require.True(t, schemas.Exists("petstore.PetDTO", "description"))
	require.Equal(t, "PetDTO the pet dto.", schemas.Search("petstore.PetDTO", "description").Data().(string))

	// petstore.Tag
	require.True(t, schemas.Exists("petstore.Tag", "description"))
	require.Equal(t, "Tag the tag model.", schemas.Search("petstore.Tag", "description").Data().(string))

	// 4-doc-comments.NoDescription
	require.True(t, schemas.Exists("4-doc-comments.NoDescription"))
	require.False(t, schemas.Exists("4-doc-comments.NoDescription", "description"))
}
