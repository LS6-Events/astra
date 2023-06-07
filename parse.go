package gengo

import (
	"errors"
)

func (s *Service) Parse() error {
	s.Log.Info().Msg("Begin parsing")

	if len(s.Inputs) == 0 {
		err := errors.New("input not set")
		s.Log.Error().Err(err).Msg("Error parsing")
		return err
	}

	if len(s.Outputs) == 0 {
		err := errors.New("output not set")
		s.Log.Error().Err(err).Msg("Error parsing")
		return err
	}

	s.Log.Info().Msg("Creating temp dir")
	err := setupTempDir()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error creating temp dir")
		return err
	}
	defer func() {
		s.Log.Info().Msg("Cleaning up temp dir")
		err := cleanupTempDir()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error cleaning up temp dir")
		} else {
			s.Log.Info().Msg("Cleaning up temp dir complete")
		}
	}()
	s.Log.Info().Msg("Creating temp dir complete")

	s.Log.Info().Msg("Populating inputs")
	for _, input := range s.Inputs {
		s.Log.Info().Str("mode", string(input.Mode)).Msg("Populating input")
		err = input.Populate(s)
		if err != nil {
			s.Log.Error().Err(err).Str("mode", string(input.Mode)).Msg("Error populating input")
			return err
		}
		s.Log.Info().Str("mode", string(input.Mode)).Msg("Populating input complete")
	}
	s.Log.Info().Msg("Populating inputs complete")

	s.Log.Info().Msg("Processing found types")
	err = s.process()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error processing found types")
		return err
	}
	s.Log.Info().Msg("Processing found types complete")

	s.Log.Info().Msg("Cleaning up structs")
	err = s.clean()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error cleaning up structs")
		return err
	}
	s.Log.Info().Msg("Cleaning up structs complete")

	s.Log.Info().Msg("Generating outputs")
	for _, output := range s.Outputs {
		s.Log.Info().Str("mode", string(output.Mode)).Msg("Generating output")
		err = output.Generate(s)
		if err != nil {
			s.Log.Error().Err(err).Str("mode", string(output.Mode)).Msg("Error generating output")
			return err
		}
		s.Log.Info().Str("mode", string(output.Mode)).Msg("Generating output complete")
	}
	s.Log.Info().Msg("Generating outputs complete")

	s.Log.Info().Msg("Parsing complete")

	return nil
}
