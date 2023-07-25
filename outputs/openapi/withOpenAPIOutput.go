package openapi

import (
	"github.com/ls6-events/gengo"
)

// WithOpenAPIOutput adds an OpenAPI specification as an output to the service
// It will generate a JSON/YAML file (based on file path [default JSON]) with the routes and components
// It should also contain the configuration for the file path to store in the cache for CLI usage
func WithOpenAPIOutput(filePath string) gengo.Option {
	return func(s *gengo.Service) {
		s.Outputs = append(s.Outputs, gengo.Output{
			Mode:     gengo.OutputModeOpenAPI,
			Generate: generate(filePath),
			Configuration: gengo.IOConfiguration{
				gengo.IOConfigurationKeyFilePath: filePath,
			},
		})
	}
}
