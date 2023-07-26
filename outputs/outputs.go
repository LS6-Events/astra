package outputs

import (
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/outputs/json"
	"github.com/ls6-events/gengo/outputs/openapi"
)

const (
	OutputModeJSON    gengo.OutputMode = "json"    // JSON file - primarily used for debugging
	OutputModeOpenAPI gengo.OutputMode = "openapi" // OpenAPI 3.0 file
)

func addOutput(mode gengo.OutputMode, generate gengo.ServiceFunction, configuration gengo.IOConfiguration) gengo.Option {
	return func(s *gengo.Service) {
		s.Outputs = append(s.Outputs, gengo.Output{
			Mode:          mode,
			Generate:      generate,
			Configuration: configuration,
		})
	}
}

// WithOpenAPIOutput adds an OpenAPI specification as an output to the service
// It will generate a JSON/YAML file (based on file path [default JSON]) with the routes and components
// It should also contain the configuration for the file path to store in the cache for CLI usage
func WithOpenAPIOutput(filePath string) gengo.Option {
	return addOutput(
		OutputModeOpenAPI,
		openapi.Generate(filePath),
		gengo.IOConfiguration{
			gengo.IOConfigurationKeyFilePath: filePath,
		},
	)
}

// WithJSONOutput adds JSON as an output to the service
// It will generate a JSON file with the routes and components
// It should also contain the configuration for the file path to store in the cache for CLI usage
func WithJSONOutput(filePath string) gengo.Option {
	return addOutput(
		OutputModeJSON,
		json.Generate(filePath),
		gengo.IOConfiguration{
			gengo.IOConfigurationKeyFilePath: filePath,
		},
	)

}
