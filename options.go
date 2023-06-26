package gengo

func WithCustomWorkDir(wd string) Option {
	return func(s *Service) {
		s.WorkDir = wd
	}
}
