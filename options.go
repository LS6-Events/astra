package astra

// WithCustomWorkDir is an option to set the working directory of the service to a custom directory.
func WithCustomWorkDir(wd string) Option {
	return func(s *Service) {
		s.WorkDir = wd
	}
}
