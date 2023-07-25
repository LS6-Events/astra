package gengo

// AddToBeProcessed adds a type to the list of types that have to be processed by the generator
// It ensures that the type is not already in the list of types to be processed or already processed
func (s *Service) AddToBeProcessed(pkg string, name string) {
	_ = s.HasAddedToBeProcessed(pkg, name)
}

// HasAddedToBeProcessed checks if a type has already been added to the list of types to be processed
// It also checks if the type has already been processed
// If the type is successfully added to the list of types to be processed, it returns true otherwise false
func (s *Service) HasAddedToBeProcessed(pkg string, name string) bool {
	s.Log.Debug().Str("pkg", pkg).Str("name", name).Msg("Adding to be processed")

	for _, p := range s.ToBeProcessed {
		if p.Name == name && p.Pkg == pkg {
			s.Log.Trace().Str("pkg", pkg).Str("name", name).Msg("Already in to be processed queue")
			return false
		}
	}
	for _, v := range s.Components {
		if v.Name == name && v.Package == pkg {
			s.Log.Trace().Str("pkg", pkg).Str("name", name).Msg("Already processed")
			return false
		}
	}
	s.ToBeProcessed = append(s.ToBeProcessed, Processable{
		Name: name,
		Pkg:  pkg,
	})

	return true
}
