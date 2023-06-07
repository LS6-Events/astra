package gengo

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func New(cfgs ...Option) *Service {
	s := &Service{}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	s.Log = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.InfoLevel)

	for _, cfg := range cfgs {
		cfg(s)
	}

	s.Routes = make([]Route, 0)
	s.ToBeProcessed = make([]Processable, 0)
	s.ReturnTypes = make([]Field, 0)

	s.Log.Debug().Msg("Service created")

	return s
}
