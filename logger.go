package astra

import (
	"github.com/rs/zerolog"
)

type CustomLoggerOption struct{}

// WithCustomLogger enables the custom logger for the generator
// It needs to pass in a zerolog.Logger
func (o CustomLoggerOption) With(logger zerolog.Logger) FunctionalOption {
	return func(s *Service) {
		s.Log = logger
	}
}

// WithCustomLogLevel enables the custom log level for the generator
// This allows for easier debugging and controlling of log levels
func (o CustomLoggerOption) WithLevel(level string) FunctionalOption {
	return func(s *Service) {
		var logLevel zerolog.Level
		switch level {
		case zerolog.LevelTraceValue:
			logLevel = zerolog.TraceLevel
		case zerolog.LevelDebugValue:
			logLevel = zerolog.DebugLevel
		case zerolog.LevelInfoValue:
			logLevel = zerolog.InfoLevel
		case zerolog.LevelWarnValue:
			logLevel = zerolog.WarnLevel
		case zerolog.LevelErrorValue:
			logLevel = zerolog.ErrorLevel
		case zerolog.LevelFatalValue:
			logLevel = zerolog.FatalLevel
		case zerolog.LevelPanicValue:
			logLevel = zerolog.PanicLevel
		default:
			logLevel = zerolog.InfoLevel
		}
		s.Log = s.Log.Level(logLevel)
	}
}

func (o CustomLoggerOption) LoadFromPlugin(s *Service, p *ConfigurationPlugin) error {
	loggerSymbol, found := p.Lookup("Logger")
	if found {
		if logger, ok := loggerSymbol.(zerolog.Logger); ok {
			o.With(logger)(s)
		}
	}

	levelSymbol, found := p.Lookup("LogLevel")
	if found {
		if level, ok := levelSymbol.(string); ok {
			o.WithLevel(level)(s)
		}
	}

	return nil
}

func init() {
	RegisterOption(CustomLoggerOption{})
}
