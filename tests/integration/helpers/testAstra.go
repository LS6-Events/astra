package helpers

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/inputs"
	"github.com/ls6-events/astra/outputs"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func SetupTestAstraWithDefaultConfig(t *testing.T, r *gin.Engine) (*gabs.Container, error) {
	t.Helper()

	config := &astra.Config{
		Host: "localhost",
		Port: 8000,
	}

	return SetupTestAstra(t, r, config)
}

func SetupTestAstra(t *testing.T, r *gin.Engine, config *astra.Config) (*gabs.Container, error) {
	t.Helper()

	gen := astra.New(inputs.WithGinInput(r), outputs.WithOpenAPIOutput("./output.json"))

	gen.SetConfig(config)

	err := gen.Parse()
	require.NoError(t, err)

	fileContents, err := os.ReadFile("./output.json")
	require.NoError(t, err)

	return gabs.ParseJSON(fileContents)
}
