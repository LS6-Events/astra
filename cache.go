package gengo

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

const cacheFileName = "cache.json"

func (s *Service) Cache() error {
	s.Log.Debug().Msg("Caching service")

	var cachePath string
	if s.cachePath != "" {
		cachePath = s.cachePath
	} else {
		cachePath = path.Join(s.getGenGoDirPath(), cacheFileName)
	}

	if _, err := os.Stat(cachePath); err == nil {
		err := s.ClearCache()
		if err != nil {
			return err
		}
	}

	f, err := os.Create(cachePath)
	if err != nil {
		return err
	}

	cacheStr, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = f.Write(cacheStr)
	if err != nil {
		return err
	}

	s.Log.Debug().Msg("Cached service")
	return nil
}

func (s *Service) LoadCache() error {
	s.Log.Debug().Msg("Loading cached service")

	var cachePath string
	if s.cachePath != "" {
		cachePath = s.cachePath
	} else {
		cachePath = path.Join(s.getGenGoDirPath(), cacheFileName)
	}

	if _, err := os.Stat(cachePath); err != nil {
		return nil
	}

	if err := s.LoadCacheFromCustomPath(cachePath); err != nil {
		return err
	}

	s.Log.Debug().Msg("Loaded cached service")
	return nil
}

func (s *Service) LoadCacheFromCustomPath(cachePath string) error {
	f, err := os.Open(cachePath)
	if err != nil {
		return err
	}

	var service Service
	if strings.HasSuffix(cachePath, ".json") {
		err = json.NewDecoder(f).Decode(&service)
	} else if strings.HasSuffix(cachePath, ".yaml") || strings.HasSuffix(cachePath, ".yml") {
		err = yaml.NewDecoder(f).Decode(&service)
	} else {
		err = errors.New("unsupported file format")
	}

	if err != nil {
		return err
	}

	s.Inputs = service.Inputs
	s.Outputs = service.Outputs
	s.Config = service.Config
	s.Routes = service.Routes
	s.Components = service.Components
	s.ToBeProcessed = service.ToBeProcessed
	return nil
}

func (s *Service) ClearCache() error {
	s.Log.Debug().Msg("Clearing cached service")

	var cachePath string
	if s.cachePath != "" {
		cachePath = s.cachePath
	} else {
		cachePath = path.Join(s.getGenGoDirPath(), cacheFileName)
	}

	if _, err := os.Stat(cachePath); err != nil {
		return nil
	}

	if err := os.RemoveAll(cachePath); err != nil {
		return err
	}

	s.Log.Debug().Msg("Cleared cached service")
	return nil
}

func (s *Service) IsCacheEnabled() bool {
	return s.cacheEnabled
}

func WithCache() Option {
	return func(s *Service) {
		s.cacheEnabled = true
	}
}

func WithCustomCachePath(cachePath string) Option {
	return func(s *Service) {
		s.cacheEnabled = true
		s.cachePath = cachePath
	}
}
