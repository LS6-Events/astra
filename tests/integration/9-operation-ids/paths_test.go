package integration

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOperationIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	paths := testAstra.Path("paths")

	// GET /pets
	require.Equal(t, "getAllPets", paths.Path("/pets.get.operationId").Data().(string))

	// POST /pets
	require.Equal(t, "createPet", paths.Path("/pets.post.operationId").Data().(string))

	// GET /pets/{id}
	require.Equal(t, "getPetById", paths.Path("/pets/{id}.get.operationId").Data().(string))

	// DELETE /pets/{id}
	require.Equal(t, "deletePet", paths.Path("/pets/{id}.delete.operationId").Data().(string))
}
