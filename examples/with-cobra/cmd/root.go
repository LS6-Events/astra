package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "with-cobra",
	Short: "with-cobra is a sample application that demonstrates how to use gengo with cobra to generate an OpenAPI specification for a Go application.",
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(generateCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
