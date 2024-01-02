package astra

import "errors"

// SetupParse sets up the parse by calling the Setup function and then the CreateRoutes function.
// It will create the routes from the inputs but not parse them.
// SetupParse should be called before CompleteParse.
// SetupParse should be used if you are using the CLI or if you want to parse the routes yourself.
func (s *Service) SetupParse() error {
	s.Log.Info().Msg("Setting up parse")

	if len(s.Inputs) == 0 {
		err := errors.New("input not set")
		s.Log.Error().Err(err).Msg("Error setting up parse")
		return err
	}

	if len(s.Outputs) == 0 {
		err := errors.New("output not set")
		s.Log.Error().Err(err).Msg("Error setting up parse")
		return err
	}

	err := s.Setup()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error parsing")
		return err
	}

	err = s.CreateRoutes()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error creating routes from inputs")
		return err
	}

	s.Log.Info().Msg("Setting up parse complete")

	return nil
}

// CompleteParse completes the parse by calling the ParseRoutes function, Process function, Clean function and Generate function.
// It will parse the routes from the files identified, process the found definitions, clean up the structs and generate the outputs.
// CompleteParse should be called after SetupParse.
func (s *Service) CompleteParse() error {
	s.Log.Info().Msg("Completing parse")

	if len(s.Inputs) == 0 {
		err := errors.New("input not set")
		s.Log.Error().Err(err).Msg("Error completing parse")
		return err
	}

	if len(s.Outputs) == 0 {
		err := errors.New("output not set")
		s.Log.Error().Err(err).Msg("Error completing parse")
		return err
	}

	err := s.ParseRoutes()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error parsing routes from inputs")
		return err
	}

	err = s.Clean()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error cleaning up structs")
		return err
	}

	err = s.Generate()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error generating outputs")
		return err
	}

	err = s.Teardown()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error tearing down")
		return err
	}

	s.Log.Info().Msg("Completing parse complete")

	return nil
}

// Parse parses the inputs and generates the outputs.
// It will call the SetupParse function and then the CompleteParse function.
// It should be used if you don't intend to use the CLI.
func (s *Service) Parse() error {
	s.Log.Info().Msg("Begin parsing")

	err := s.SetupParse()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error setting up parse")
		return err
	}

	err = s.CompleteParse()
	if err != nil {
		s.Log.Error().Err(err).Msg("Error completing parse")
		return err
	}

	s.Log.Info().Msg("Parsing complete")

	return nil
}
