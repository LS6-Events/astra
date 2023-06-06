package gengo

func (s *Service) AddToBeProcessed(pkg string, name string) {
	for _, p := range s.ToBeProcessed {
		if p.Name == name && p.Pkg == pkg {
			return
		}
	}
	for _, v := range s.ReturnTypes {
		if v.Name == name && v.Package == pkg {
			return
		}
	}
	s.ToBeProcessed = append(s.ToBeProcessed, Processable{
		Name: name,
		Pkg:  pkg,
	})
}
