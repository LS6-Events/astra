package astra

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// New creates a new generator service.
// It takes in a list of options that can be used to configure the generator.
// It will also setup the logger for the generator and setup the slices that are used to store the routes, inputs, outputs and components.
func New(opts ...Option) *Service {
	s := &Service{}

	//nolint:reassign // This is needed to set the stack trace marshaling function.
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	s.Log = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.InfoLevel)

	s.Routes = make([]Route, 0)
	s.Components = make([]Field, 0)
	s.CustomTypeMapping = make(map[string]TypeFormat)

	for _, opt := range opts {
		opt(s)
	}

	s.Log.Debug().Msg("Service created")

	return s
}
