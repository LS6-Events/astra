package json

import (
	"encoding/json"
	"github.com/ls6-events/gengo"
	"os"
	"strings"
)

type JSONOutput struct {
	Routes []gengo.Route `json:"routes"`
	Fields []gengo.Field `json:"fields"`
}

func generate(filePath string) gengo.GenerateFunction {
	return func(s *gengo.Service) error {
		output := JSONOutput{
			Routes: s.Routes,
			Fields: s.ReturnTypes,
		}

		file, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}

		if !strings.HasSuffix(filePath, ".json") {
			filePath += ".json"
		}

		err = os.WriteFile(filePath, file, 0644)
		if err != nil {
			return err
		}

		return nil
	}
}
