package cache

import "github.com/ls6-events/gengo"

// WithCache Option to enable the cache
func WithCache() gengo.Option {
	return func(s *gengo.Service) {
		s.CacheEnabled = true
	}
}

// WithCustomCachePath Option to enable the cache with a custom path
func WithCustomCachePath(cachePath string) gengo.Option {
	return func(s *gengo.Service) {
		s.CacheEnabled = true
		s.CachePath = cachePath
	}
}
