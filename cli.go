package gengo

func WithCLI() Option {
	return func(s *Service) {
		s.cacheEnabled = true
		s.CLIMode = CLIModeSetup
	}
}

func WithCLIBuilder() Option {
	return func(s *Service) {
		s.CLIMode = CLIModeBuilder
	}
}
