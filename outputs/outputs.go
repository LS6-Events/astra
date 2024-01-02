package outputs

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/outputs/azureFunctions"
	"github.com/ls6-events/astra/outputs/json"
	"github.com/ls6-events/astra/outputs/openapi"
)

const (
	OutputModeAzureFunctions astra.OutputMode = "azureFunctions" // Azure Functions HTTP Trigger Bindings.
	OutputModeJSON           astra.OutputMode = "json"           // JSON file - primarily used for debugging.
	OutputModeOpenAPI        astra.OutputMode = "openapi"        // OpenAPI 3.0 file.
)

func addOutput(mode astra.OutputMode, generate astra.ServiceFunction, configuration astra.IOConfiguration) astra.Option {
	return func(s *astra.Service) {
		s.Outputs = append(s.Outputs, astra.Output{
			Mode:          mode,
			Generate:      generate,
			Configuration: configuration,
		})
	}
}

// WithAzureFunctionsOutput adds Azure Functions HTTP Trigger Bindings as an output to the service.
func WithAzureFunctionsOutput(directoryPath string) astra.Option {
	return addOutput(
		OutputModeAzureFunctions,
		azureFunctions.Generate(directoryPath),
		astra.IOConfiguration{
			astra.IOConfigurationKeyDirectoryPath: directoryPath,
		},
	)
}

// WithOpenAPIOutput adds an OpenAPI specification as an output to the service.
// It will generate a JSON/YAML file (based on file path [default JSON]) with the routes and components.
// It should also contain the configuration for the file path to store in the cache for CLI usage.
func WithOpenAPIOutput(filePath string) astra.Option {
	return addOutput(
		OutputModeOpenAPI,
		openapi.Generate(filePath),
		astra.IOConfiguration{
			astra.IOConfigurationKeyFilePath: filePath,
		},
	)
}

// WithJSONOutput adds JSON as an output to the service.
// It will generate a JSON file with the routes and components.
// It should also contain the configuration for the file path to store in the cache for CLI usage.
func WithJSONOutput(filePath string) astra.Option {
	return addOutput(
		OutputModeJSON,
		json.Generate(filePath),
		astra.IOConfiguration{
			astra.IOConfigurationKeyFilePath: filePath,
		},
	)
}
