package gengo

func New(cfgs ...Option) *Service {
	s := &Service{}

	for _, cfg := range cfgs {
		cfg(s)
	}

	s.Routes = make([]Route, 0)
	s.ToBeProcessed = make([]Processable, 0)
	s.ReturnTypes = make([]Field, 0)

	return s
}
