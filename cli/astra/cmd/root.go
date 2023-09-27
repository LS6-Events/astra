package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// The root command is "astra"
var rootCmd = &cobra.Command{
	Use:   "astra",
	Short: "CLI version of Astra tool",
	Long:  `Generate specifications for your Go services, importing from your Go web server code and exporting to the OpenAPI standard`,
}

func Execute() error {
	return rootCmd.Execute()
}
