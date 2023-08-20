package gengo

import (
	"errors"
	"io/fs"
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
	if err := os.MkdirAll(s.getGenGoDirPath(), 0755); err != nil {
		return err
	}

	if err := s.setupGitIgnore(); err != nil {
		return err
	}

	return nil
}

func (s *Service) SetupTempOutputDir(dirPath string) (string, error) {
	tempDir := path.Join(s.getGenGoDirPath(), dirPath)

	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func (s *Service) MoveTempOutputDir(dirPath, outputDir string) error {
	sourcePath := path.Join(s.getGenGoDirPath(), dirPath)
	destinationPath := path.Join(s.WorkDir, outputDir)

	_, err := os.Stat(destinationPath)
	if !errors.Is(err, fs.ErrNotExist) {
		err := os.RemoveAll(destinationPath)
		if err != nil {
			return err
		}
	}

	err = os.Rename(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	return nil
}

// setupGitIgnore creates the .gitignore file, by default it will ignore all files in the directory
func (s *Service) setupGitIgnore() error {
	gitIgnorePath := path.Join(s.getGenGoDirPath(), ".gitignore")
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
	if err := os.RemoveAll(s.getGenGoDirPath()); err != nil {
		return err
	}

	return nil
}
