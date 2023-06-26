package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gengo",
	Short: "Generate specifications for your Go services",
	Long:  `Generate specifications for your Go services, importing from your Go web server code and exporting to OpenAPI and JSON`,
}

func Execute() error {
	return rootCmd.Execute()
}
