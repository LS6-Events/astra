package gengo

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"reflect"
	"strings"
)

// Process is the function that processes the type definitions of the specified return types
// It does this by loading the packages that are needed for the generator
// It then finds the types in the packages that have returned
// It will then process them using the types package
func (s *Service) Process() error {
	s.Log.Info().Msg("Begin processing found definitions")

	newPass := true
	for newPass {
		newPass = false

		s.Log.Debug().Int("len", len(s.ToBeProcessed)).Msg("Processing flagged definitions")
		pkgs, err := s.loadPackages()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error loading packages")
			return err
		}
		s.Log.Debug().Int("len", len(s.ToBeProcessed)).Msg("Loaded packages")

		s.ToBeProcessed = make([]Processable, 0)

		for _, pkg := range pkgs {
			s.Log.Debug().Str("pkg", pkg.PkgPath).Msg("Processing package")
			for _, name := range s.typesByName[pkg.PkgPath] {
				s.Log.Debug().Str("name", name).Msg("Processing type")
				targetType := pkg.Types.Scope().Lookup(name)
				if targetType == nil {
					err := fmt.Errorf("type %s not found in package %s", name, pkg.PkgPath)
					s.Log.Error().Err(err).Msg("Error processing type")
					return err
				}

				if !targetType.Exported() {
					err := fmt.Errorf("type %s is not exported", name)
					s.Log.Error().Err(err).Msg("Error processing type")
					return err
				}

				var newPassElem bool
				var resultField Field

				underlyingType := targetType.Type().Underlying()

				pointerPass := true
				for pointerPass {
					pointerPass = false
					switch t := underlyingType.(type) {
					case *types.Pointer:
						underlyingType = t.Elem()
						pointerPass = true
					case *types.Struct:
						newPassElem, resultField = s.processStruct(t, pkg)
					case *types.Map:
						newPassElem, resultField = s.processMap(t, pkg)
					case *types.Slice:
						newPassElem, resultField = s.processSlice(t, pkg)
					default:
						newPassForType, relativePkg, relativeName := s.processType(t.String(), pkg.PkgPath)
						if newPassForType {
							newPassElem = true
						}
						resultField = Field{
							Package: relativePkg,
							Type:    relativeName,
						}
					}
				}

				resultField.Name = name

				s.Components = append(s.Components, resultField)
				s.Log.Debug().Str("name", name).Msg("Processed type")

				if newPassElem {
					newPass = true
				}
			}
		}
	}

	s.Log.Info().Msg("Processing found definitions complete")

	if s.CacheEnabled {
		err := s.Cache()
		if err != nil {
			s.Log.Error().Err(err).Msg("Error caching")
			return err
		}
	}

	return nil
}

// processStruct processes a struct type
// It will return a bool indicating if a new pass is needed and a Field that represents the struct
// It can look at all the fields of the struct and process them individually with all their names and tags
// It will also process embedded structs
// TODO: Add support for validation tags
func (s *Service) processStruct(t *types.Struct, pkg *packages.Package) (bool, Field) {
	newPass := false
	fields := make(map[string]Field)

	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)

		// Get if "binding:required" tag is present and "json" as well
		isRequired := false
		isEmbedded := f.Embedded()
		isShown := true
		name := f.Id()

		tag := t.Tag(i)
		if tag != "" {
			binding := reflect.StructTag(tag).Get("binding")
			if binding != "" {
				isRequired = strings.Contains(binding, "required")
			}

			yaml := reflect.StructTag(tag).Get("yaml")
			if yaml != "" && yaml != "-" {
				isShown = true
				name = strings.Split(yaml, ",")[0]
			} else if yaml == "-" && isShown {
				isShown = false
			}

			xml := reflect.StructTag(tag).Get("xml")
			if xml != "" && xml != "-" {
				isShown = true
				name = strings.Split(xml, ",")[0]
			} else if xml == "-" && isShown {
				isShown = false
			}

			form := reflect.StructTag(tag).Get("form")
			if form != "" && form != "-" {
				isShown = true
				name = strings.Split(form, ",")[0]
			} else if form == "-" && isShown {
				isShown = false
			}

			json := reflect.StructTag(tag).Get("json")
			if json != "" && json != "-" {
				isShown = true
				name = strings.Split(json, ",")[0]
			} else if json == "-" && isShown {
				isShown = false
			}
		}

		if !isShown {
			continue
		}

		switch f.Type().(type) {
		case *types.Basic:
			// If the field is a basic type, we don't need to add a package
			fields[name] = Field{
				Name:       f.Id(),
				Type:       f.Type().String(),
				IsRequired: isRequired,
				IsEmbedded: isEmbedded,
			}
		case *types.Named:
			switch f.Type().Underlying().(type) {
			case *types.Map:
				newPassForMap, mapField := s.processMap(f.Type().Underlying().(*types.Map), pkg)
				if newPassForMap {
					newPass = true
				}
				mapField.IsRequired = isRequired
				mapField.IsEmbedded = isEmbedded
				fields[name] = mapField
			case *types.Slice:
				newPassForSlice, sliceField := s.processSlice(f.Type().Underlying().(*types.Slice), pkg)
				if newPassForSlice {
					newPass = true
				}
				sliceField.IsRequired = isRequired
				sliceField.IsEmbedded = isEmbedded
				fields[name] = sliceField
			default:
				newPassForType, relativePkg, relativeName := s.processType(f.Type().String(), pkg.PkgPath)
				if newPassForType {
					newPass = true
				}
				fields[name] = Field{
					Package:    relativePkg,
					Type:       relativeName,
					IsRequired: isRequired,
					IsEmbedded: isEmbedded,
					Name:       f.Id(),
				}
			}
		}
	}

	return newPass, Field{
		Type:         "struct",
		StructFields: fields,
		Package:      pkg.PkgPath,
	}
}

// processMap processes a map type
// It will return a bool indicating if a new pass is needed and a Field that represents the map
// It will also process the key and value types of the map
func (s *Service) processMap(t *types.Map, pkg *packages.Package) (bool, Field) {
	newPassKey, relativeKeyPkg, relativeKeyName := s.processType(t.Key().String(), pkg.PkgPath)
	newPassValue, relativeValuePkg, relativeValueName := s.processType(t.Elem().String(), pkg.PkgPath)

	resultField := Field{
		Package:   relativeValuePkg,
		Type:      "map",
		MapKeyPkg: relativeKeyPkg,
		MapKey:    relativeKeyName,
		MapValue:  relativeValueName,
	}

	return newPassKey || newPassValue, resultField
}

// processSlice processes a slice type
// It will return a bool indicating if a new pass is needed and a Field that represents the slice
// It will also process the element type of the slice
func (s *Service) processSlice(t *types.Slice, pkg *packages.Package) (bool, Field) {
	newPass, relativePkg, relativeName := s.processType(t.Elem().String(), pkg.PkgPath)
	resultField := Field{
		Package:   relativePkg,
		Type:      "slice",
		SliceType: relativeName,
	}

	return newPass, resultField
}

// processType processes a type
// It will return a bool indicating if a new pass is needed and a Field that represents the type
// It will also process the package and name of the type
// If the type has a dot, it will split it and use the first part as the package and the second part as the name
// Therefore it will return the package and name of the type
// Otherwise it will use the default package and the type as the name
func (s *Service) processType(t string, defaultPkg string) (newPass bool, relativePkg string, relativeName string) {
	split := strings.Split(t, ".")
	if len(split) > 1 {
		relativePkg = strings.Join(split[:len(split)-1], ".")
		relativeName = split[len(split)-1]
	} else {
		relativePkg = defaultPkg
		relativeName = split[0]
	}

	if !IsAcceptedType(relativeName) {
		// We don't mind the duplicate records here, as we check the existing components before adding to queue
		newPass = s.HasAddedToBeProcessed(relativePkg, relativeName)
	} else {
		newPass = false
	}

	return
}
