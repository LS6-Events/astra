package astra

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSerivce_getAstraDirPath(t *testing.T) {
	t.Run("returns the path to the .astra directory relatively", func(t *testing.T) {
		service := &Service{}

		require.Equal(t, ".astra", service.getAstraDirPath())
	})

	t.Run("returns the path to the .astra directory relative to the working directory", func(t *testing.T) {
		service := &Service{
			WorkDir: "/test/",
		}

		require.Equal(t, "/test/.astra", service.getAstraDirPath())
	})
}

func TestService_setupAstraDir(t *testing.T) {
	t.Run("creates the .astra directory", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		_, err = os.Stat(".astra")
		require.NoError(t, err)
	})

	t.Run("creates the .gitignore file", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		_, err = os.Stat(".astra/.gitignore")
		require.NoError(t, err)
	})
}

func TestService_SetupTempOutputDir(t *testing.T) {
	t.Run("creates the directory", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		tempDir, err := service.SetupTempOutputDir("test")
		require.NoError(t, err)

		_, err = os.Stat(tempDir)
		require.NoError(t, err)
	})

	t.Run("returns the path to the directory", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		tempDir, err := service.SetupTempOutputDir("test")
		require.NoError(t, err)

		require.Equal(t, ".astra/test", tempDir)
	})
}

func TestService_MoveTempOutputDir(t *testing.T) {
	t.Run("moves the directory", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		_, err = service.SetupTempOutputDir("test")
		require.NoError(t, err)

		err = service.MoveTempOutputDir("test", "output")
		require.NoError(t, err)

		_, err = os.Stat("output")
		require.NoError(t, err)
	})

	t.Run("removes the temp directory", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		tempDir, err := service.SetupTempOutputDir("test")
		require.NoError(t, err)

		err = service.MoveTempOutputDir("test", "output")
		require.NoError(t, err)
		defer func() {
			err := os.Remove("output")
			require.NoError(t, err)
		}()

		_, err = os.Stat(tempDir)
		require.Error(t, err)
	})
}

func TestService_setupGitIgnore(t *testing.T) {
	t.Run("creates the .gitignore file", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		err = service.setupGitIgnore()
		require.NoError(t, err)

		_, err = os.Stat(".astra/.gitignore")
		require.NoError(t, err)
	})

	t.Run("does not create the .gitignore file if it already exists", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		err = service.setupGitIgnore()
		require.NoError(t, err)

		f, err := os.Create(".astra/.gitignore")
		require.NoError(t, err)
		defer func() {
			err := f.Close()
			require.NoError(t, err)
		}()

		err = service.setupGitIgnore()
		require.NoError(t, err)
	})

	t.Run("has * as content", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)
		defer func() {
			err := service.cleanupAstraDir()
			require.NoError(t, err)
		}()

		err = service.setupGitIgnore()
		require.NoError(t, err)

		f, err := os.Open(".astra/.gitignore")
		require.NoError(t, err)
		defer func() {
			err := f.Close()
			require.NoError(t, err)
		}()

		content := make([]byte, 1)
		_, err = f.Read(content)
		require.NoError(t, err)

		require.Equal(t, "*", string(content))
	})
}

func TestService_cleanupAstraDir(t *testing.T) {
	t.Run("removes the .astra directory", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)

		err = service.cleanupAstraDir()
		require.NoError(t, err)

		_, err = os.Stat(".astra")
		require.Error(t, err)
	})

	t.Run("removes the .gitignore file", func(t *testing.T) {
		service := &Service{}

		err := service.setupAstraDir()
		require.NoError(t, err)

		err = service.cleanupAstraDir()
		require.NoError(t, err)

		_, err = os.Stat(".astra/.gitignore")
		require.Error(t, err)
	})
}
