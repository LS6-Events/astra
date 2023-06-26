package gengo

func (s *Service) Generate() error {
	s.Log.Info().Msg("Generating outputs")
	for _, output := range s.Outputs {
		s.Log.Info().Str("mode", string(output.Mode)).Msg("Generating output")
		err := output.Generate(s)
		if err != nil {
			s.Log.Error().Err(err).Str("mode", string(output.Mode)).Msg("Error generating output")
			return err
		}
		s.Log.Info().Str("mode", string(output.Mode)).Msg("Generating output complete")
	}

	s.Log.Info().Msg("Generating outputs complete")

	if s.cacheEnabled {
		err := s.Cache()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error caching")
			return err
		}
	}

	return nil
}
