package gengo

// WithCLI enables the CLI mode for the generator
// This will run the generator to only crawl for the file names, function names and line numbers of the functions that need to be analysed.
func WithCLI() Option {
	return func(s *Service) {
		s.cacheEnabled = true
		s.CLIMode = CLIModeSetup
	}
}

// WithCLIBuilder enables the CLI mode for the generator
// This will run the generator utilising existing cache for the file names, function names and line numbers of the functions that need to be analysed to generate the code.
func WithCLIBuilder() Option {
	return func(s *Service) {
		s.CLIMode = CLIModeBuilder
	}
}
