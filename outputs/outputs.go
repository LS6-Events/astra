package outputs

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/outputs/azureFunctions"
	"github.com/ls6-events/astra/outputs/json"
	"github.com/ls6-events/astra/outputs/openapi"
)

const (
	OutputModeAzureFunctions astra.OutputMode = "azureFunctions" // Azure Functions HTTP Trigger Bindings
	OutputModeJSON           astra.OutputMode = "json"           // JSON file - primarily used for debugging
	OutputModeOpenAPI        astra.OutputMode = "openapi"        // OpenAPI 3.0 file
)

func addOutput(mode astra.OutputMode, generate astra.ServiceFunction, configuration astra.IOConfiguration) astra.FunctionalOption {
	return func(s *astra.Service) {
		s.Outputs = append(s.Outputs, astra.Output{
			Mode:          mode,
			Generate:      generate,
			Configuration: configuration,
		})
	}
}

type AzureFunctionsOutputOption struct{}

// WithAzureFunctionsOutput adds Azure Functions HTTP Trigger Bindings as an output to the service
func (o AzureFunctionsOutputOption) With(directoryPath string) astra.FunctionalOption {
	return addOutput(
		OutputModeAzureFunctions,
		azureFunctions.Generate(directoryPath),
		astra.IOConfiguration{
			astra.IOConfigurationKeyDirectoryPath: directoryPath,
		},
	)
}

func (o AzureFunctionsOutputOption) LoadFromPlugin(s *astra.Service, p *astra.ConfigurationPlugin) error {
	directoryPathSymbol, found := p.Lookup("OutputAzureFunctionsDirectoryPath")
	if found {
		if directoryPath, ok := directoryPathSymbol.(string); ok {
			o.With(directoryPath)(s)
		}
	}

	return nil
}

type OpenAPIOutputOption struct{}

// WithOpenAPIOutput adds an OpenAPI specification as an output to the service
// It will generate a JSON/YAML file (based on file path [default JSON]) with the routes and components
// It should also contain the configuration for the file path to store in the cache for CLI usage
func (o OpenAPIOutputOption) With(filePath string) astra.FunctionalOption {
	return addOutput(
		OutputModeOpenAPI,
		openapi.Generate(filePath),
		astra.IOConfiguration{
			astra.IOConfigurationKeyFilePath: filePath,
		},
	)
}

func (o OpenAPIOutputOption) LoadFromPlugin(s *astra.Service, p *astra.ConfigurationPlugin) error {
	filePathSymbol, found := p.Lookup("OutputOpenAPIFilePath")
	if found {
		if filePath, ok := filePathSymbol.(string); ok {
			o.With(filePath)(s)
		}
	}

	return nil
}

type JSONOutputOption struct{}

// WithJSONOutput adds JSON as an output to the service
// It will generate a JSON file with the routes and components
// It should also contain the configuration for the file path to store in the cache for CLI usage
func (o JSONOutputOption) With(filePath string) astra.FunctionalOption {
	return addOutput(
		OutputModeJSON,
		json.Generate(filePath),
		astra.IOConfiguration{
			astra.IOConfigurationKeyFilePath: filePath,
		},
	)
}

func (o JSONOutputOption) LoadFromPlugin(s *astra.Service, p *astra.ConfigurationPlugin) error {
	filePathSymbol, found := p.Lookup("OutputJSONFilePath")
	if found {
		if filePath, ok := filePathSymbol.(string); ok {
			o.With(filePath)(s)
		}
	}

	return nil
}

func init() {
	astra.RegisterOption(AzureFunctionsOutputOption{})
	astra.RegisterOption(OpenAPIOutputOption{})
	astra.RegisterOption(JSONOutputOption{})
}
