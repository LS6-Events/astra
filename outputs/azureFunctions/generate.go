package azureFunctions

import (
	"encoding/json"
	"github.com/ls6-events/gengo"
	"os"
	"path"
	"strings"
)

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

func generate(directoryPath string) gengo.ServiceFunction {
	return func(s *gengo.Service) error {
		s.Log.Debug().Msg("Generating Azure Functions output")

		tempOutputDirectoryPath, err := s.SetupTempOutputDir(gengo.OutputModeAzureFunctions)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to create directory")
			return err
		}

		for _, route := range s.Routes {
			splitHandler := strings.Split(route.Handler, ".")
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

		err = s.MoveTempOutputDir(gengo.OutputModeAzureFunctions, directoryPath)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to move temporary output directory to final location")
			return err
		}

		s.Log.Info().Msg("Generated JSON output")
		return nil
	}
}
