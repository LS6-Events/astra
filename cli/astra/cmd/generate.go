package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/cli"
	"github.com/spf13/cobra"
)

var (
	cacheFile = ".astra/cache.json" // Location of the cache.json file
	cwd       = "."                 // Current working directory (where main.go is located)
)

// generateCmd represents the generate command
// It is used to generate the service from a cache file
// It requires the cache file to be passed in, and the working directory of the main.go file
// By default the cache file is .astra/cache.json and the working directory is the current directory
// Example: astra generate -c .astra/cache.json -d .
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the service",
	Long:  `Generate the service by parsing the inputs and outputs and calculating the routes using the options setup where the service was defined`,
	Run: func(cmd *cobra.Command, args []string) {

		if !path.IsAbs(cwd) {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Failed to get current working directory: %s\n", err.Error())
				os.Exit(1)
			}

			cwd = path.Join(wd, cwd)
		}

		s := astra.New(cli.WithCLIBuilder(), astra.WithCustomWorkDir(cwd))

		err := s.LoadCacheFromCustomPath(cacheFile)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to load cache")
			os.Exit(1)
		}

		err = rebindOptions(s)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to rebind options")
			os.Exit(1)
		}

		err = s.CompleteParse()
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to generate service")
			os.Exit(1)
		}

		s.Log.Info().Msg("Service built")
	},
}

func init() {
	generateCmd.Flags().StringVarP(&cacheFile, "cache", "c", cacheFile, "Location of the cache.json file")
	generateCmd.Flags().StringVarP(&cwd, "dir", "d", cwd, "Current working directory (where main.go is located)")
	rootCmd.AddCommand(generateCmd)
}
