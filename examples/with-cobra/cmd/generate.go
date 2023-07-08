package cmd

import (
	"github.com/ls6-events/gengo"
	gengoGin "github.com/ls6-events/gengo/inputs/gin"
	"github.com/ls6-events/gengo/outputs/openapi"
	"github.com/spf13/cobra"
	"withcobra/routes"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the OpenAPI specification using gengo",
	Run: func(cmd *cobra.Command, args []string) {
		router := routes.GetRouter()

		gen := gengo.New(gengoGin.WithGinInput(router), openapi.WithOpenAPIOutput("openapi.generated.yaml"))

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
