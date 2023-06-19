package gengo

import "github.com/rs/zerolog"

type Option func(*Service)

type Service struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`

	Log zerolog.Logger `json:"-"`

	Config *Config `json:"config"`

	Routes []Route `json:"routes"`

	ToBeProcessed []Processable `json:"-"`
	typesByName   map[string][]string

	tempMainPackageName string

	cacheEnabled bool

	Components []Field `json:"components"`
}
