package azureFunctions

import "github.com/ls6-events/gengo"

func WithAzureFunctionsOutput(directoryPath string) gengo.Option {
	return func(s *gengo.Service) {
		s.Outputs = append(s.Outputs, gengo.Output{
			Mode:     gengo.OutputModeAzureFunctions,
			Generate: generate(directoryPath),
		})
	}
}
