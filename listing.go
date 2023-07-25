package gengo

import "regexp"

func WithPathDenyList(denyList string) Option {
	return func(s *Service) {
		s.PathDenyList = append(s.PathDenyList, func(path string) bool {
			return path == denyList
		})
	}
}

func WithPathDenyListRegex(regex *regexp.Regexp) Option {
	return func(s *Service) {
		s.PathDenyList = append(s.PathDenyList, func(path string) bool {
			return regex.MatchString(path)
		})
	}
}

func WithPathDenyListFunc(denyList func(string) bool) Option {
	return func(s *Service) {
		s.PathDenyList = append(s.PathDenyList, denyList)
	}
}
