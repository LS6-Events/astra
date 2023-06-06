package gengo

import (
	"errors"
	"fmt"
)

func (s *Service) Parse() error {
	if len(s.Inputs) == 0 {
		return errors.New("input not set")
	}

	if len(s.Outputs) == 0 {
		return errors.New("output not set")
	}

	err := setupTempDir()
	if err != nil {
		return err
	}
	defer func() {
		err := cleanupTempDir()
		if err != nil {
			fmt.Printf("error cleaning up temp dir: %s\n", err)
		}
	}()

	for _, input := range s.Inputs {
		err = input.Populate(s)
		if err != nil {
			return err
		}
	}

	err = s.process()
	if err != nil {
		return err
	}

	err = s.clean()
	if err != nil {
		return err
	}

	for _, output := range s.Outputs {
		err = output.Generate(s)
		if err != nil {
			return err
		}
	}

	return nil
}
