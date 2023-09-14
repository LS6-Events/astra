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
		s.Components[i] = s.cleanField(s.Components[i], mainPkg)
	}

	for i := 0; i < len(s.Routes); i++ {
		for j := 0; j < len(s.Routes[i].ReturnTypes); j++ {
			s.Routes[i].ReturnTypes[j].Field = s.cleanField(s.Routes[i].ReturnTypes[j].Field, mainPkg)
		}
		for j := 0; j < len(s.Routes[i].PathParams); j++ {
			s.Routes[i].PathParams[j].Field = s.cleanField(s.Routes[i].PathParams[j].Field, mainPkg)
		}
		for j := 0; j < len(s.Routes[i].QueryParams); j++ {
			s.Routes[i].QueryParams[j].Field = s.cleanField(s.Routes[i].QueryParams[j].Field, mainPkg)
		}
		for j := 0; j < len(s.Routes[i].Body); j++ {
			s.Routes[i].Body[j].Field = s.cleanField(s.Routes[i].Body[j].Field, mainPkg)
		}
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

func (s *Service) cleanField(f Field, mainPkg string) Field {
	s.HandleSpecialType(&f)

	if f.Package == mainPkg {
		f.Package = "main"
	}
	if f.MapKeyPackage == mainPkg {
		f.MapKeyPackage = "main"
	}

	for k, v := range f.StructFields {
		f.StructFields[k] = s.cleanField(v, mainPkg)
	}

	return f
}
