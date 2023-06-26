package cmd

import (
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/inputs/gin"
	"github.com/ls6-events/gengo/outputs/json"
	"github.com/ls6-events/gengo/outputs/openapi"
)

func rebindOptions(s *gengo.Service) error {
	err := rebindInputs(s)
	if err != nil {
		return err
	}

	err = rebindOutputs(s)
	if err != nil {
		return err
	}

	return nil
}

func rebindInputs(s *gengo.Service) error {
	inputsCopy := s.Inputs
	s.Inputs = nil
	for _, input := range inputsCopy {
		switch input.Mode {
		case gengo.InputModeGin:
			gin.WithGinInput(nil)(s)
		default:
			return gengo.ErrInputModeNotFound
		}
	}

	return nil
}

func rebindOutputs(s *gengo.Service) error {
	outputsCopy := s.Outputs
	s.Outputs = nil
	for _, output := range outputsCopy {
		switch output.Mode {
		case gengo.OutputModeJSON:
			filePath, ok := output.Configuration[gengo.IOConfigurationKeyFilePath].(string)
			if !ok || filePath == "" {
				return gengo.ErrOutputFilePathRequired
			}
			json.WithJSONOutput(filePath)(s)
		case gengo.OutputModeOpenAPI:
			filePath, ok := output.Configuration[gengo.IOConfigurationKeyFilePath].(string)
			if !ok || filePath == "" {
				return gengo.ErrOutputFilePathRequired
			}
			openapi.WithOpenAPIOutput(filePath)(s)
		default:
			return gengo.ErrOutputModeNotFound
		}
	}

	return nil
}
