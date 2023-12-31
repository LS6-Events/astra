package astra

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

// The caching mechanism is used to cache the service in a file so that it can be loaded later on
// At the moment it is only used the CLI to load the service and the files that are needed to be crawled by the AST parser with their respective inputs
// Plans are for the future to use it as a change only mechanism to only generate the files that have changed since the last build

const cacheFileName = "cache.json"

// Cache the service in a file
func (s *Service) Cache() error {
	s.Log.Debug().Msg("Caching service")

	var cachePath string
	if s.CachePath != "" {
		cachePath = s.CachePath
	} else {
		cachePath = path.Join(s.getAstraDirPath(), cacheFileName)
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

	var cacheStr []byte
	if strings.HasSuffix(cachePath, ".yaml") || strings.HasSuffix(cachePath, ".yml") {
		cacheStr, err = yaml.Marshal(s)
	} else {
		cacheStr, err = json.Marshal(s)
	}
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

// LoadCache Load the service from a file cache
// If the file does not exist, it will not return an error
func (s *Service) LoadCache() error {
	s.Log.Debug().Msg("Loading cached service")

	var cachePath string
	if s.CachePath != "" {
		cachePath = s.CachePath
	} else {
		cachePath = path.Join(s.getAstraDirPath(), cacheFileName)
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

// LoadCacheFromCustomPath Load the service from a file cache
// If the file does not exist, it will return an error
// Requires the path to the cache file
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
	return nil
}

// ClearCache Clear the cache file
func (s *Service) ClearCache() error {
	s.Log.Debug().Msg("Clearing cached service")

	var cachePath string
	if s.CachePath != "" {
		cachePath = s.CachePath
	} else {
		cachePath = path.Join(s.getAstraDirPath(), cacheFileName)
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
