package gengo

import (
	"os"
	"path"
)

const tempDir = ".gengo_temp"

func getTempDirPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path.Join(cwd, tempDir)
}

func setupTempDir() error {
	tempDirPath := getTempDirPath()
	if err := os.MkdirAll(tempDirPath, 0755); err != nil {
		return err
	}

	return nil
}

func cleanupTempDir() error {
	tempDirPath := getTempDirPath()
	if err := os.RemoveAll(tempDirPath); err != nil {
		return err
	}

	return nil
}
