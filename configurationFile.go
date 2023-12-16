package astra

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"plugin"
	"runtime"
	"strings"
)

const CONFIGURATION_FILE_NAME = "astra.config.go"
const CONFIGURATION_FILE_PATH_COMPILED = "config/astra.config.so"

var (
	ErrConfigurationFileDoesNotContainAnyValidDeclarations = errors.New("configuration file does not contain any valid declarations")
	ErrGoOSNotSupported                                    = errors.New("GOOS is not supported")
)

func (s *Service) doesConfigurationFileExist() (bool, error) {
	// Scan the main package for the configuration file
	// If it exists, return true
	// If it doesn't exist, return false
	_, err := os.Stat(path.Join(s.WorkDir, CONFIGURATION_FILE_NAME))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

// Get every declaration in the configuration file and see if any of them have the ASTRA_ prefix
func (s *Service) hasConfigurationFileGotAnyValidDeclarationPrefixes() (bool, error) {
	fset := token.FileSet{}

	file, err := parser.ParseFile(&fset, path.Join(s.WorkDir, CONFIGURATION_FILE_NAME), nil, parser.AllErrors)
	if err != nil {
		return false, err
	}

	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			funcDecl := decl.(*ast.FuncDecl)
			if strings.HasPrefix(funcDecl.Name.Name, CONFIGURATION_FUNCTION_PREFIX) {
				return true, nil
			}
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.ValueSpec:
					valueSpec := spec.(*ast.ValueSpec)
					for _, ident := range valueSpec.Names {
						if strings.HasPrefix(ident.Name, CONFIGURATION_FUNCTION_PREFIX) {
							return true, nil
						}
					}
				}
			}
		}
	}

	return false, nil
}

func (s *Service) compileConfigurationFile() error {
	// If the OS is not Linux, FreeBSD or macOS then plugins are not supported
	// https://golang.org/pkg/plugin/
	// So we will return an error here
	if osType := runtime.GOOS; osType != "linux" && osType != "freebsd" && osType != "darwin" {
		s.Log.Error().Str("GOOS", osType).Msg("GOOS is not supported")
		return ErrGoOSNotSupported
	}

	// If the configuration plugin file already exists, remove it
	if _, err := os.Stat(path.Join(s.getAstraDirPath(), CONFIGURATION_FILE_PATH_COMPILED)); err == nil {
		err := os.Remove(path.Join(s.getAstraDirPath(), CONFIGURATION_FILE_PATH_COMPILED))
		if err != nil {
			return err
		}
	}

	dirPath := s.getAstraDirPath()
	configurationOutputFilePath := path.Join(dirPath, CONFIGURATION_FILE_PATH_COMPILED)
	// We compile the entire directory to allow for imports from other files
	err := exec.Command("go", "build", "-buildmode=plugin", "-o", configurationOutputFilePath, ".").Run()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) loadPluginFile() error {
	// Load the plugin file
	// If it fails, return an error
	// If it succeeds, return nil
	dirPath := s.getAstraDirPath()
	configurationOutputFilePath := path.Join(dirPath, CONFIGURATION_FILE_PATH_COMPILED)
	configPlugin, err := plugin.Open(configurationOutputFilePath)
	if err != nil {
		return err
	}

	s.ConfigurationPlugin = NewConfigurationPlugin(configPlugin)

	return nil
}

func (s *Service) loadPluginOptions() error {
	spew.Dump(pluginOptionsList)
	for _, pluginOptions := range pluginOptionsList {
		err := pluginOptions.LoadFromPlugin(s, s.ConfigurationPlugin)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) TestConfigurationFile() error {
	doesExist, err := s.doesConfigurationFileExist()
	if err != nil {
		return err
	}

	if doesExist {
		containsConfig, err := s.hasConfigurationFileGotAnyValidDeclarationPrefixes()
		if err != nil {
			return err
		}

		s.Log.Info().Msg("Configuration file exists")

		if !containsConfig {
			return ErrConfigurationFileDoesNotContainAnyValidDeclarations
		}

		s.Log.Info().Msg("Configuration file contains config")

		err = s.compileConfigurationFile()
		if err != nil {
			return err
		}
		s.Log.Info().Msg("Configuration file compiled")

		err = s.loadPluginFile()
		if err != nil {
			return err
		}
		s.Log.Info().Msg("Configuration file loaded")

		err = s.loadPluginOptions()
		if err != nil {
			return err
		}

		s.Log.Info().Msg("Configuration file options loaded")
	} else {
		s.Log.Info().Msg("Configuration file does not exist")
	}

	return nil
}
