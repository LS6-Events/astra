package cmd

import (
	"fmt"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/cli"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var (
	cacheFile string = ".gengo/cache.json" // Location of the cache.json file
	cwd       string = "."                 // Current working directory (where main.go is located)
)

// generateCmd represents the generate command
// It is used to generate the service from a cache file
// It requires the cache file to be passed in, and the working directory of the main.go file
// By default the cache file is .gengo/cache.json and the working directory is the current directory
// Example: gengo generate -c .gengo/cache.json -d .
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

		s := gengo.New(cli.WithCLIBuilder(), gengo.WithCustomWorkDir(cwd))

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
