package cmd

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/inputs"
	"github.com/ls6-events/astra/outputs"
	"github.com/spf13/cobra"
	"withcobra/routes"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the OpenAPI specification using astra",
	Run: func(cmd *cobra.Command, args []string) {
		router := routes.GetRouter()

		gen := astra.New(inputs.WithGinInput(router), outputs.WithOpenAPIOutput("openapi.generated.yaml"))

		config := astra.Config{
			Title:   "Example API with Cobra",
			Version: "1.0.0",
			Host:    "localhost",
			Port:    8000,
		}

		gen.SetConfig(&config)

		err := gen.Parse()
		if err != nil {
			panic(err)
		}
	},
}
