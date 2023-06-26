package gengo

func (s *Service) Teardown() error {
	s.Log.Info().Msg("Tearing down")

	if !s.cacheEnabled && s.CLIMode == CLIModeNone {
		s.Log.Info().Msg("Cleaning up temp dir")
		err := s.cleanupGenGoDir()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error cleaning up temp dir")
		} else {
			s.Log.Info().Msg("Cleaning up temp dir complete")
		}
	} else {
		s.Log.Info().Msg("Leaving temp dir but cleaning up main package")
		err := s.cleanupTempMainPackage()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error cleaning up main package")
		} else {
			s.Log.Info().Msg("Cleaning up main package complete")
		}
	}

	s.Log.Info().Msg("Tearing down complete")
	return nil
}
