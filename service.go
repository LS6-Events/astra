package gengo

import (
	"github.com/rs/zerolog"
)

// Option is a function that can be used to configure the generator
type Option func(*Service)

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

	ToBeProcessed []Processable `json:"-"`
	typesByName   map[string][]string

	Components []Field `json:"components" yaml:"components"`

	tempMainPackageName string
	WorkDir             string `json:"-"`

	CacheEnabled bool    `json:"-"`
	CachePath    string  `json:"-"`
	CLIMode      CLIMode `json:"-"`

	PathBlacklist []func(string) bool `json:"-"`
}
