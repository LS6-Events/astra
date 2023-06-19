package gengo

import (
	"os"
	"path"
)

const gengoDir = ".gengo"

func getGenGoDirPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path.Join(cwd, gengoDir)
}

func setupGenGoDir() error {
	tempDirPath := getGenGoDirPath()
	if err := os.MkdirAll(tempDirPath, 0755); err != nil {
		return err
	}

	if err := setupGitIgnore(); err != nil {
		return err
	}

	return nil
}

func setupGitIgnore() error {
	tempDirPath := getGenGoDirPath()
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

func cleanupGenGoDir() error {
	tempDirPath := getGenGoDirPath()
	if err := os.RemoveAll(tempDirPath); err != nil {
		return err
	}

	return nil
}
