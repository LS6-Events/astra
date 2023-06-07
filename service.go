package gengo

import "github.com/rs/zerolog"

type Option func(*Service)

type Service struct {
	Inputs  []Input
	Outputs []Output

	Log zerolog.Logger

	Config *Config

	Routes []Route

	ToBeProcessed []Processable
	typesByName   map[string][]string

	tempMainPackageName string

	ReturnTypes []Field
}
