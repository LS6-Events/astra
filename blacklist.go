package gengo

import "regexp"

func WithPathBlacklist(blacklist string) Option {
	return func(s *Service) {
		s.PathBlacklist = append(s.PathBlacklist, func(path string) bool {
			return path == blacklist
		})
	}
}

func WithPathBlacklistRegex(regex *regexp.Regexp) Option {
	return func(s *Service) {
		s.PathBlacklist = append(s.PathBlacklist, func(path string) bool {
			return regex.MatchString(path)
		})
	}
}

func WithPathBlacklistFunc(blacklist func(string) bool) Option {
	return func(s *Service) {
		s.PathBlacklist = append(s.PathBlacklist, blacklist)
	}
}
