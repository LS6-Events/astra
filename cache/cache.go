package cache

import "github.com/ls6-events/astra"

// WithCache Option to enable the cache
func WithCache() astra.Option {
	return func(s *astra.Service) {
		s.CacheEnabled = true
	}
}

// WithCustomCachePath Option to enable the cache with a custom path
func WithCustomCachePath(cachePath string) astra.Option {
	return func(s *astra.Service) {
		s.CacheEnabled = true
		s.CachePath = cachePath
	}
}
