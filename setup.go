package astra

import (
	"os"
)

// Setup sets up the service by creating the temp dir and caching
// Setup should be called before anything else in the service
func (s *Service) Setup() error {
	s.Log.Info().Msg("Setting up")

	s.Log.Info().Msg("Creating temp dir")
	err := s.setupAstraDir()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error creating temp dir")
		return err
	}
	s.Log.Info().Msg("Creating temp dir complete")

	if s.WorkDir == "" {
		s.Log.Info().Msg("Noting current dir")
		cwd, err := os.Getwd()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error noting current dir")
			return err
		}
		s.WorkDir = cwd

		s.Log.Info().Msg("Noting current dir complete")
	}

	if s.CacheEnabled {
		err := s.Cache()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error caching")
			return err
		}
	}

	s.Log.Info().Msg("Setting up complete")
	return nil
}
