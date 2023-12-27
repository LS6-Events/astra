package snapshot

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/inputs"
	"github.com/ls6-events/astra/outputs"
	"github.com/ls6-events/astra/tests/snapshot/comparison"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/pets", getAllPets)
	r.GET("/pets/:id", getPetByID)
	r.POST("/pets", createPet)
	r.DELETE("/pets/:id", deletePet)

	return r
}

func TestSnapshot(t *testing.T) {
	r := setupRouter()

	config := &astra.Config{
		Host: "localhost",
		Port: 8000,
	}

	gen := astra.New(inputs.WithGinInput(r), outputs.WithOpenAPIOutput("./output.yaml"))

	gen.SetConfig(config)

	err := gen.Parse()
	require.NoError(t, err)

	if os.Getenv("GENERATE_SNAPSHOT") != "true" {
		// Compare the generated snapshot with the existing one
		comparison.CompareYAML(t, "./snapshot.yaml", "./output.yaml")
	} else {
		// Overwrite the existing snapshot with the new one
		err = os.Rename("./output.yaml", "./snapshot.yaml")
		require.NoError(t, err)
	}
}
