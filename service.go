package gengo

import (
	"github.com/rs/zerolog"
)

type Option func(*Service)

type CLIMode string

const (
	CLIModeNone    CLIMode = ""
	CLIModeSetup   CLIMode = "setup"
	CLIModeBuilder CLIMode = "builder"
)

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
	WorkDir             string

	cacheEnabled bool
	cachePath    string
	CLIMode      CLIMode `json:"-"`
}
