package astra

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

// This creates a new directory in the current working directory called .astra.
// The idea is that this directory will be used to store the cache file and any other files that are needed for the generator.
// It will create a .gitignore file in the directory so that it is not committed to git.

const astraDir = ".astra"

// getAstraDirPath returns the path to the .astra directory.
func (s *Service) getAstraDirPath() string {
	return path.Join(s.WorkDir, astraDir)
}

// setupAstraDir creates the .astra directory and the .gitignore file.
func (s *Service) setupAstraDir() error {
	if err := os.MkdirAll(s.getAstraDirPath(), 0755); err != nil {
		return err
	}

	if err := s.setupGitIgnore(); err != nil {
		return err
	}

	return nil
}

func (s *Service) SetupTempOutputDir(dirPath string) (string, error) {
	tempDir := path.Join(s.getAstraDirPath(), dirPath)

	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func (s *Service) MoveTempOutputDir(dirPath, outputDir string) error {
	sourcePath := path.Join(s.getAstraDirPath(), dirPath)
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

// setupGitIgnore creates the .gitignore file, by default it will ignore all files in the directory.
func (s *Service) setupGitIgnore() error {
	gitIgnorePath := path.Join(s.getAstraDirPath(), ".gitignore")
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

// cleanupAstraDir removes the .astra directory.
func (s *Service) cleanupAstraDir() error {
	if err := os.RemoveAll(s.getAstraDirPath()); err != nil {
		return err
	}

	return nil
}
