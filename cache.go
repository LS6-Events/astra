package gengo

import (
	"encoding/json"
	"os"
	"path"
)

const cacheFileName = "cache.json"

func (s *Service) Cache() error {
	s.Log.Debug().Msg("Caching service")

	cachePath := path.Join(getGenGoDirPath(), cacheFileName)
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

	cachePath := path.Join(getGenGoDirPath(), cacheFileName)
	if _, err := os.Stat(cachePath); err != nil {
		return nil
	}

	f, err := os.Open(cachePath)
	if err != nil {
		return err
	}

	var service Service
	err = json.NewDecoder(f).Decode(&service)
	if err != nil {
		return err
	}

	s.Config = service.Config
	s.Routes = service.Routes
	s.Components = service.Components
	s.ToBeProcessed = service.ToBeProcessed
	s.Log.Debug().Msg("Loaded cached service")
	return nil
}

func (s *Service) ClearCache() error {
	s.Log.Debug().Msg("Clearing cached service")

	cachePath := path.Join(getGenGoDirPath(), cacheFileName)
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
