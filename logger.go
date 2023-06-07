package gengo

import "github.com/rs/zerolog"

func WithCustomLogger(logger zerolog.Logger) Option {
	return func(s *Service) {
		s.Log = logger
	}
}

func WithCustomLogLevel(level string) Option {
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
