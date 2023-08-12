package json

import (
	"encoding/json"
	"github.com/ls6-events/gengo"
	"os"
	"path"
	"strings"
)

// JSONOutput is the output of the JSON output
// It will in essence be a copy of the service's output (routes and components)
type JSONOutput struct {
	Routes     []gengo.Route `json:"routes"`
	Components []gengo.Field `json:"components"`
}

// generate the JSON output
// It will marshal the JSONOutput struct and write it to a file
func Generate(filePath string) gengo.ServiceFunction {
	return func(s *gengo.Service) error {
		s.Log.Info().Msg("Generating JSON output")
		output := JSONOutput{
			Routes:     s.Routes,
			Components: s.Components,
		}

		s.Log.Debug().Msg("Generated JSON output")
		file, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to marshal JSON output")
			return err
		}

		if !strings.HasSuffix(filePath, ".json") {
			s.Log.Debug().Str("filePath", filePath).Msg("Adding .json suffix to file path")
			filePath += ".json"
		}

		s.Log.Debug().Str("filePath", filePath).Msg("Writing JSON output to file")
		filePath = path.Join(s.WorkDir, filePath)
		err = os.WriteFile(filePath, file, 0644)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to write JSON output to file")
			return err
		}

		s.Log.Info().Msg("Generated JSON output")
		return nil
	}
}
