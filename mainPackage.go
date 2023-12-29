package astra

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// This file contains functionality for copying the main package to a temporary directory and replacing the package name with a different name.
// This is required so that the generator can parse the types in the main package should they be required, and follow any functions that are required.
// The package is cleaned up after the generator has finished running.

const mainPackageReplacement = "astramain"
const mainPackageReplacementPath = astraDir + "/" + mainPackageReplacement

// setupTempMainPackage copies the main package to a temporary directory and replaces the package name with a different name.
func (s *Service) setupTempMainPackage() error {
	var pkgName string

	newMainPkgPath := path.Join(s.getAstraDirPath(), mainPackageReplacement)
	if _, err := os.Stat(newMainPkgPath); err == nil {
		err := os.RemoveAll(newMainPkgPath)
		if err != nil {
			return err
		}
	}
	err := os.Mkdir(newMainPkgPath, 0755)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(s.WorkDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") {
			fileName := path.Join(s.WorkDir, file.Name())
			fileData, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}

			fileData = []byte(strings.ReplaceAll(string(fileData), "package main", fmt.Sprintf("package %s", mainPackageReplacement)))

			newFilePath := path.Join(newMainPkgPath, file.Name())
			err = os.WriteFile(newFilePath, fileData, 0644)
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(file.Name(), ".mod") {
			fileName := path.Join(s.WorkDir, file.Name())
			fileData, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}

			pkgName = strings.Replace(strings.Split(string(fileData), "\n")[0], "module ", "", 1)
		}
	}

	newMainPkg := path.Join(pkgName, mainPackageReplacementPath)

	s.tempMainPackageName = newMainPkg

	return nil
}

// cleanupTempMainPackage removes the temporary main package.
func (s *Service) cleanupTempMainPackage() error {
	newMainPkgPath := path.Join(s.getAstraDirPath(), mainPackageReplacement)
	if _, err := os.Stat(newMainPkgPath); err == nil {
		err := os.RemoveAll(newMainPkgPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMainPackageName returns the name of the temporary main package.
func (s *Service) GetMainPackageName() (string, error) {
	if s.tempMainPackageName == "" {
		err := s.setupTempMainPackage()
		if err != nil {
			return "", errors.Join(errors.New("failed to setup temporary main package"), err)
		}
	}
	return s.tempMainPackageName, nil
}
