package gengo

func (s *Service) AddToBeProcessed(pkg string, name string) {
	_ = s.HasAddedToBeProcessed(pkg, name)
}

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
