package cli

import "github.com/ls6-events/gengo"

// WithCLI enables the CLI mode for the generator
// This will run the generator to only crawl for the file names, function names and line numbers of the functions that need to be analysed.
func WithCLI() gengo.Option {
	return func(s *gengo.Service) {
		s.CacheEnabled = true
		s.CLIMode = gengo.CLIModeSetup
	}
}

// WithCLIBuilder enables the CLI mode for the generator
// This will run the generator utilising existing cache for the file names, function names and line numbers of the functions that need to be analysed to generate the code.
func WithCLIBuilder() gengo.Option {
	return func(s *gengo.Service) {
		s.CLIMode = gengo.CLIModeBuilder
	}
}
