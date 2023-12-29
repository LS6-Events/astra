package astra

import "github.com/rs/zerolog"

// WithCustomLogger enables the custom logger for the generator.
// It needs to pass in a zerolog.Logger.
func WithCustomLogger(logger zerolog.Logger) Option {
	return func(s *Service) {
		s.Log = logger
	}
}

// WithCustomLogLevel enables the custom log level for the generator.
// This allows for easier debugging and controlling of logs.
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
