package cmd

import (
	"github.com/spf13/cobra"
	"withcobra/routes"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		router := routes.GetRouter()

		err := router.Run(":8000")
		if err != nil {
			panic(err)
		}
	},
}
