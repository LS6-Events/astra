package json

import (
	"github.com/ls6-events/gengo"
)

// WithJSONOutput adds JSON as an output to the service
// It will generate a JSON file with the routes and components
// It should also contain the configuration for the file path to store in the cache for CLI usage
func WithJSONOutput(filePath string) gengo.Option {
	return func(s *gengo.Service) {
		s.Outputs = append(s.Outputs, gengo.Output{
			Mode:     gengo.OutputModeJSON,
			Generate: generate(filePath),
			Configuration: gengo.IOConfiguration{
				gengo.IOConfigurationKeyFilePath: filePath,
			},
		})
	}
}
