package cmd

import (
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/inputs"
	"github.com/ls6-events/gengo/outputs"
	"github.com/spf13/cobra"
	"withcobra/routes"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the OpenAPI specification using gengo",
	Run: func(cmd *cobra.Command, args []string) {
		router := routes.GetRouter()

		gen := gengo.New(inputs.WithGinInput(router), outputs.WithOpenAPIOutput("openapi.generated.yaml"))

		config := gengo.Config{
			Title:   "Example API with Cache",
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
