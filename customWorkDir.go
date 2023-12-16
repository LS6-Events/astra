package astra

type CustomWorkDirOption struct{}

func (o CustomWorkDirOption) With(wd string) FunctionalOption {
	return func(s *Service) {
		s.WorkDir = wd
	}
}

func (o CustomWorkDirOption) LoadFromPlugin(s *Service, p *ConfigurationPlugin) error {
	wdSymbol, found := p.Lookup("WorkDir")
	if found {
		if wd, ok := wdSymbol.(string); ok {
			o.With(wd)(s)
		}
	}

	return nil
}

func init() {
	RegisterOption(CustomWorkDirOption{})
}
