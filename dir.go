package gengo

import (
	"os"
	"path"
)

const gengoDir = ".gengo"

func (s *Service) getGenGoDirPath() string {
	return path.Join(s.WorkDir, gengoDir)
}

func (s *Service) setupGenGoDir() error {
	tempDirPath := s.getGenGoDirPath()
	if err := os.MkdirAll(tempDirPath, 0755); err != nil {
		return err
	}

	if err := s.setupGitIgnore(); err != nil {
		return err
	}

	return nil
}

func (s *Service) setupGitIgnore() error {
	tempDirPath := s.getGenGoDirPath()
	gitIgnorePath := path.Join(tempDirPath, ".gitignore")
	if _, err := os.Stat(gitIgnorePath); err == nil {
		return nil
	}

	f, err := os.Create(gitIgnorePath)
	if err != nil {
		return err
	}

	_, err = f.WriteString("*\n")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) cleanupGenGoDir() error {
	tempDirPath := s.getGenGoDirPath()
	if err := os.RemoveAll(tempDirPath); err != nil {
		return err
	}

	return nil
}
