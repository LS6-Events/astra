package azureFunctions

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/ls6-events/astra"
)

const tempOutputDir = "azure"
const outputFile = "function.json"

type AzureFunctionsBinding struct {
	Type      string   `json:"type"`
	Direction string   `json:"direction"`
	Name      string   `json:"name"`
	Methods   []string `json:"methods,omitempty"`
	Route     string   `json:"route,omitempty"`
}

type AzureFunctionsOutput struct {
	Bindings []AzureFunctionsBinding `json:"bindings"`
}

func Generate(directoryPath string) astra.ServiceFunction {
	return func(s *astra.Service) error {
		s.Log.Debug().Msg("Generating Azure Functions output")

		tempOutputDirectoryPath, err := s.SetupTempOutputDir(tempOutputDir)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to create directory")
			return err
		}

		for _, route := range s.Routes {
			splitHandler := strings.Split(route.Handlers[len(route.Handlers)-1].Name, ".")
			funcName := splitHandler[len(splitHandler)-1]

			functionDirectoryPath := path.Join(tempOutputDirectoryPath, funcName)
			filePath := path.Join(functionDirectoryPath, outputFile)

			err := os.Mkdir(functionDirectoryPath, 0755)
			if err != nil {
				s.Log.Error().Err(err).Msg("Failed to create directory")
				return err
			}

			output := AzureFunctionsOutput{
				Bindings: []AzureFunctionsBinding{
					{
						Type:      "httpTrigger",
						Direction: "in",
						Name:      "req",
						Methods:   associatedMethods(route.Method),
						Route:     convertRoute(route),
					},
					{
						Type:      "http",
						Direction: "out",
						Name:      "res",
					},
				},
			}

			s.Log.Debug().Msg("Generated JSON output")
			file, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				s.Log.Error().Err(err).Msg("Failed to marshal JSON output")
				return err
			}

			s.Log.Debug().Str("filePath", filePath).Msg("Writing JSON output to file")
			err = os.WriteFile(filePath, file, 0644)
			if err != nil {
				s.Log.Error().Err(err).Msg("Failed to write JSON output to file")
				return err
			}
		}

		err = s.MoveTempOutputDir(tempOutputDir, directoryPath)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to move temporary output directory to final location")
			return err
		}

		s.Log.Info().Msg("Generated JSON output")
		return nil
	}
}
