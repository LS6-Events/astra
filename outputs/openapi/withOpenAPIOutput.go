package openapi

import (
	"github.com/ls6-events/gengo"
)

func WithOpenAPIOutput(filePath string) gengo.Option {
	return func(s *gengo.Service) {
		s.Outputs = append(s.Outputs, gengo.Output{
			Mode:     gengo.OutputModeOpenAPI,
			Generate: generate(filePath),
		})
	}
}
