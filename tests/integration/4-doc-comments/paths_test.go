package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPathsDocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	paths := testAstra.Path("paths")

	// GET /no-docs
	require.True(t, paths.Exists("/no-docs", "get"))
	require.False(t, paths.Path("/no-docs.get").Exists("description"))

	// GET /pets
	require.Equal(t, "getAllPets returns all pets.", paths.Path("/pets.get.description").Data().(string))

	// POST /pets
	require.Equal(t, "createPet creates a pet.\nIt takes in a Pet without an ID in the request body.", paths.Path("/pets.post.description").Data().(string))

	// GET /pets/{id}
	require.Equal(t, "getPetByID returns a pet by its ID.\nIt takes in the ID as a path parameter.", paths.Path("/pets/{id}.get.description").Data().(string))

	// DELETE /pets/{id}
	require.Equal(t, "deletePet deletes a pet by its ID.\nIt takes in the ID as a path parameter.", paths.Path("/pets/{id}.delete.description").Data().(string))
}
