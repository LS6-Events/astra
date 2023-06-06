package gengo

import (
	"golang.org/x/tools/go/packages"
	"os"
)

func (s *Service) loadPackages() ([]*packages.Package, error) {
	s.typesByName = make(map[string][]string, 0)
	patterns := make([]string, 0)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for _, structToProcess := range s.ToBeProcessed {
		if structToProcess.Pkg == "main" {
			structToProcess.Pkg, err = s.GetMainPackageName()
			if err != nil {
				return nil, err
			}
		}

		patterns = append(patterns, structToProcess.Pkg)
		if _, ok := s.typesByName[structToProcess.Pkg]; !ok {
			s.typesByName[structToProcess.Pkg] = make([]string, 0)
		}
		s.typesByName[structToProcess.Pkg] = append(s.typesByName[structToProcess.Pkg], structToProcess.Name)
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps | packages.NeedName,
		Dir:  cwd,
	}, patterns...)

	if err != nil {
		return nil, err
	}

	return pkgs, nil
}
