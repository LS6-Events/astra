package json

import (
	"github.com/ls6-events/gengo"
)

func WithJSONOutput(filePath string) gengo.Option {
	return func(s *gengo.Service) {
		s.Outputs = append(s.Outputs, gengo.Output{
			Mode:     gengo.OutputModeJSON,
			Generate: generate(filePath),
		})
	}
}
