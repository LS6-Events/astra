package gengo

import (
	"os"
	"path"
)

// This creates a new directory in the current working directory called .gengo
// The idea is that this directory will be used to store the cache file and any other files that are needed for the generator
// It will create a .gitignore file in the directory so that it is not committed to git

const gengoDir = ".gengo"

// getGenGoDirPath returns the path to the .gengo directory
func (s *Service) getGenGoDirPath() string {
	return path.Join(s.WorkDir, gengoDir)
}

// setupGenGoDir creates the .gengo directory and the .gitignore file
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

// setupGitIgnore creates the .gitignore file, by default it will ignore all files in the directory
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

// cleanupGenGoDir removes the .gengo directory
func (s *Service) cleanupGenGoDir() error {
	tempDirPath := s.getGenGoDirPath()
	if err := os.RemoveAll(tempDirPath); err != nil {
		return err
	}

	return nil
}
