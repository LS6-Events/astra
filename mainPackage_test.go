package astra

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestService_setupTempMainPackage(t *testing.T) {
	t.Run("creates temp main package", func(t *testing.T) {
		service := &Service{
			WorkDir: "./",
		}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err = service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		err = service.setupTempMainPackage()
		require.NoError(t, err)
		require.NotEmpty(t, service.tempMainPackageName)
	})

	t.Run("copies all .go files in working directory to temp directory", func(t *testing.T) {
		service := &Service{
			WorkDir: "./",
		}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err = service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		err = service.setupTempMainPackage()
		require.NoError(t, err)

		wdFiles, err := os.ReadDir("./")
		require.NoError(t, err)

		files, err := os.ReadDir("./.astra/astramain")
		require.NoError(t, err)

		filesByName := make([]string, len(files))
		for i, file := range files {
			filesByName[i] = file.Name()
		}

		require.NotEmpty(t, files)

		for _, file := range wdFiles {
			if path.Ext(file.Name()) == ".go" {
				require.Contains(t, filesByName, file.Name())
			}
		}
	})
}

func TestService_cleanupTempMainPackage(t *testing.T) {
	service := &Service{
		WorkDir: "./",
	}

	err := service.setupAstraDir()
	require.NoError(t, err)
	defer func() {
		err = service.cleanupAstraDir()
		require.NoError(t, err)
	}()

	err = service.setupTempMainPackage()
	require.NoError(t, err)

	err = service.cleanupTempMainPackage()
	require.NoError(t, err)

	_, err = os.Stat(path.Join(service.getAstraDirPath(), mainPackageReplacement))
	require.Error(t, err)
}

func TestService_GetMainPackageName(t *testing.T) {
	t.Run("sets up the temp main package if it doesn't exist", func(t *testing.T) {
		service := &Service{
			WorkDir: "./",
		}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err = service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		require.Empty(t, service.tempMainPackageName)

		mainPkgName, err := service.GetMainPackageName()
		require.NoError(t, err)

		require.NotEmpty(t, mainPkgName)
		require.NotEmpty(t, service.tempMainPackageName)

		require.Equal(t, service.tempMainPackageName, mainPkgName)
	})

	t.Run("returns the name of the temp main package", func(t *testing.T) {
		service := &Service{
			WorkDir:             "./",
			tempMainPackageName: "test",
		}

		mainPkgName, err := service.GetMainPackageName()
		require.NoError(t, err)

		require.Equal(t, "test", mainPkgName)
	})
}
