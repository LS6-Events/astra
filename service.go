package astra

import (
	"github.com/rs/zerolog"
)

// CLIMode is the mode the CLI is running in
type CLIMode string

const (
	CLIModeNone    CLIMode = ""        // Not running in CLI mode
	CLIModeSetup   CLIMode = "setup"   // Running in setup mode - used in the project code to setup the routes
	CLIModeBuilder CLIMode = "builder" // Running in builder mode - used in the project code to build the routes and genearate the types (this is not needed to be used by other developers, the CLI will use this mode)
)

// Service is the main struct for the generator
type Service struct {
	Inputs  []Input  `json:"inputs" yaml:"inputs"`
	Outputs []Output `json:"outputs" yaml:"outputs"`

	Log zerolog.Logger `json:"-"`

	Config *Config `json:"config" yaml:"config"`

	Routes []Route `json:"routes" yaml:"routes"`

	Components []Field `json:"components" yaml:"components"`

	tempMainPackageName string
	WorkDir             string `json:"-" yaml:"-"`

	CacheEnabled bool    `json:"-"`
	CachePath    string  `json:"-"`
	CLIMode      CLIMode `json:"-"`

	PathDenyList []func(string) bool `json:"-" yaml:"-"`

	CustomFuncs []CustomFunc `json:"-" yaml:"-"`

	ConfigurationPlugin *ConfigurationPlugin `json:"-" yaml:"-"`
}
