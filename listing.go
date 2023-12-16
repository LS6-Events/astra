package astra

import (
	"regexp"
)

type PathDenyListOption struct{}

func (o PathDenyListOption) With(denyList string) FunctionalOption {
	return func(s *Service) {
		s.PathDenyList = append(s.PathDenyList, func(path string) bool {
			return path == denyList
		})
	}
}

func (o PathDenyListOption) WithRegex(regex *regexp.Regexp) FunctionalOption {
	return func(s *Service) {
		s.PathDenyList = append(s.PathDenyList, func(path string) bool {
			return regex.MatchString(path)
		})
	}
}

func (o PathDenyListOption) WithFunc(denyList func(string) bool) FunctionalOption {
	return func(s *Service) {
		s.PathDenyList = append(s.PathDenyList, denyList)
	}
}

func (o PathDenyListOption) LoadFromPlugin(s *Service, p *ConfigurationPlugin) error {
	denyListSymbol, found := p.Lookup("PathDenyList")
	if found {
		if denyList, ok := denyListSymbol.(string); ok {
			o.With(denyList)(s)
		}
	}

	denyListSymbolRegex, found := p.Lookup("PathDenyListRegex")
	if found {
		if denyList, ok := denyListSymbolRegex.(*regexp.Regexp); ok {
			o.WithRegex(denyList)(s)
		}
	}

	denyListSymbolFunc, found := p.Lookup("PathDenyListFunc")
	if found {
		if denyList, ok := denyListSymbolFunc.(func(string) bool); ok {
			o.WithFunc(denyList)(s)
		}
	}

	return nil
}

func init() {
	RegisterOption(PathDenyListOption{})
}
