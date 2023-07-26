package cmd

import (
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/inputs"
	"github.com/ls6-events/gengo/outputs"
)

// These functions are used to rebind the inputs and outputs to the service, as the JSON unmarshalling does not call the functions to bind the inputs and outputs, and loses all their referenced functions

// rebindOptions is used to rebind the inputs and outputs to the service
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

// rebindInputs is used to rebind the inputs to the service
// It will have to be updated if more input modes are added
func rebindInputs(s *gengo.Service) error {
	inputsCopy := s.Inputs
	s.Inputs = nil
	for _, input := range inputsCopy {
		switch input.Mode {
		case inputs.InputModeGin:
			inputs.WithGinInput(nil)(s)
		default:
			return gengo.ErrInputModeNotFound
		}
	}

	return nil
}

// rebindOutputs is used to rebind the outputs to the service
// It will have to be updated if more output modes are added
// It utilises the configuration keys to get the file path for the output
func rebindOutputs(s *gengo.Service) error {
	outputsCopy := s.Outputs
	s.Outputs = nil
	for _, output := range outputsCopy {
		switch output.Mode {
		case outputs.OutputModeAzureFunctions:
			directoryPath, ok := output.Configuration[gengo.IOConfigurationKeyDirectoryPath].(string)
			if ok || directoryPath == "" {
				return gengo.ErrOutputDirectoryPathRequired
			}
			outputs.WithAzureFunctionsOutput(directoryPath)(s)
		case outputs.OutputModeJSON:
			filePath, ok := output.Configuration[gengo.IOConfigurationKeyFilePath].(string)
			if !ok || filePath == "" {
				return gengo.ErrOutputFilePathRequired
			}
			outputs.WithJSONOutput(filePath)(s)
		case outputs.OutputModeOpenAPI:
			filePath, ok := output.Configuration[gengo.IOConfigurationKeyFilePath].(string)
			if !ok || filePath == "" {
				return gengo.ErrOutputFilePathRequired
			}
			outputs.WithOpenAPIOutput(filePath)(s)
		default:
			return gengo.ErrOutputModeNotFound
		}
	}

	return nil
}
