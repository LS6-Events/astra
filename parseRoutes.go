package gengo

// ParseRoutes iterates over the inputs and parses the routes from them
// CreateRoutes should be called before ParseRoutes
func (s *Service) ParseRoutes() error {
	s.Log.Info().Msg("Parsing routes from inputs")
	for _, input := range s.Inputs {
		s.Log.Info().Str("mode", string(input.Mode)).Msg("Parsing routes from input")
		err := input.ParseRoutes(s)
		if err != nil {
			s.Log.Error().Err(err).Str("mode", string(input.Mode)).Msg("Error creating routes from input")
			return err
		}
		s.Log.Info().Str("mode", string(input.Mode)).Msg("Parsing routes from input complete")
	}
	s.Log.Info().Msg("Parsing routes from inputs complete")

	if s.cacheEnabled {
		err := s.Cache()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error caching")
			return err
		}
	}

	return nil
}
