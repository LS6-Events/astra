package gengo

// CreateRoutes creates routes from the inputs
func (s *Service) CreateRoutes() error {
	s.Log.Info().Msg("Creating routes from inputs")
	for _, input := range s.Inputs {
		s.Log.Info().Str("mode", string(input.Mode)).Msg("Creating routes from input")
		err := input.CreateRoutes(s)
		if err != nil {
			s.Log.Error().Err(err).Str("mode", string(input.Mode)).Msg("Error creating routes from input")
			return err
		}
		s.Log.Info().Str("mode", string(input.Mode)).Msg("Creating routes from input complete")
	}
	s.Log.Info().Msg("Creating routes from inputs complete")

	if s.CacheEnabled {
		err := s.Cache()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error caching")
			return err
		}
	}

	return nil
}
