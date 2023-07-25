package gengo

// Clean cleans up the structs
// At the moment it only changes the package name of the main package to "main"
// It also handles the "special" types
// It also caches the service after cleaning
func (s *Service) Clean() error {
	s.Log.Info().Msg("Cleaning up structs")

	mainPkg, err := s.GetMainPackageName()
	if err != nil {
		return err
	}

	for i := 0; i < len(s.Components); i++ {
		f := s.Components[i]

		s.handleSpecialType(&f)

		if f.Package == mainPkg {
			f.Package = "main"
		}

		for k, v := range f.StructFields {
			err := cleanField(k, v, mainPkg, f.StructFields)
			if err != nil {
				return err
			}
		}

		s.Components[i] = f
	}

	s.Log.Info().Msg("Cleaning up structs complete")

	if s.CacheEnabled {
		err := s.Cache()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error caching")
			return err
		}
	}

	return nil
}

func cleanField(k string, f Field, mainPkg string, returnTypes map[string]Field) error {
	if f.Package == mainPkg {
		f.Package = "main"
	}

	for k, v := range f.StructFields {
		err := cleanField(k, v, mainPkg, f.StructFields)
		if err != nil {
			return err
		}
	}

	returnTypes[k] = f

	return nil
}
