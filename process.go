package gengo

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"reflect"
	"strings"
)

func (s *Service) process() error {
	newPass := true
	for newPass {
		newPass = false

		s.Log.Debug().Int("len", len(s.ToBeProcessed)).Msg("Processing flagged types")
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
						newPass, resultField = s.processStruct(t, pkg)
					case *types.Map:
						newPass, resultField = s.processMap(t, pkg)
					case *types.Slice:
						newPass, resultField = s.processSlice(t, pkg)
					default:
						resultField = Field{
							Package: pkg.PkgPath,
							Type:    targetType.Type().String(),
						}
					}
				}

				resultField.Name = name

				s.ReturnTypes = append(s.ReturnTypes, resultField)
				s.Log.Debug().Str("name", name).Msg("Processed type")
			}
		}
	}

	return nil
}

func (s *Service) processStruct(t *types.Struct, pkg *packages.Package) (bool, Field) {
	newPass := false
	fields := make(map[string]Field)
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)

		// Get if "binding:required" tag is present and "json" as well
		isRequired := false
		name := f.Id()

		tag := t.Tag(i)
		if tag != "" {
			binding := reflect.StructTag(tag).Get("binding")
			if binding != "" {
				isRequired = strings.Contains(binding, "required")
			}

			yaml := reflect.StructTag(tag).Get("yaml")
			if yaml != "" {
				name = strings.Split(yaml, ",")[0]
			}

			xml := reflect.StructTag(tag).Get("xml")
			if xml != "" {
				name = strings.Split(xml, ",")[0]
			}

			form := reflect.StructTag(tag).Get("form")
			if form != "" {
				name = strings.Split(form, ",")[0]
			}

			json := reflect.StructTag(tag).Get("json")
			if json != "" {
				name = strings.Split(json, ",")[0]
			}
		}
		switch f.Type().(type) {
		case *types.Struct:
			fields[name] = Field{
				Package:    pkg.PkgPath,
				Name:       f.Id(),
				Type:       "struct",
				IsRequired: isRequired,
			}
			s.AddToBeProcessed(pkg.PkgPath, f.Id())
			newPass = true
		case *types.Map:
			newPassForMap, mapField := s.processMap(f.Type().Underlying().(*types.Map), pkg)
			if newPassForMap {
				newPass = true
			}
			mapField.IsRequired = isRequired
			fields[name] = mapField
		case *types.Slice:
			newPassForSlice, sliceField := s.processSlice(f.Type().Underlying().(*types.Slice), pkg)
			if newPassForSlice {
				newPass = true
			}
			sliceField.IsRequired = isRequired
			fields[name] = sliceField
		default:
			fields[name] = Field{
				Package:    pkg.PkgPath,
				Type:       f.Type().String(),
				IsRequired: isRequired,
			}
		}
	}

	return newPass, Field{
		Type:         "struct",
		StructFields: fields,
		Package:      pkg.PkgPath,
	}
}

func (s *Service) processMap(t *types.Map, pkg *packages.Package) (bool, Field) {
	resultField := Field{
		Package:  pkg.PkgPath,
		Type:     "map",
		MapKey:   t.Key().String(),
		MapValue: t.Elem().String(),
	}

	// If map key is not accepted type, then return error
	if !IsAcceptedType(resultField.MapKey) {
		s.AddToBeProcessed(pkg.PkgPath, resultField.MapKey)
		return true, resultField
	}

	// If map value is custom struct type, then process it
	if !IsAcceptedType(resultField.MapValue) {
		s.AddToBeProcessed(pkg.PkgPath, resultField.MapValue)
		return true, resultField
	}

	return false, resultField
}

func (s *Service) processSlice(t *types.Slice, pkg *packages.Package) (bool, Field) {
	resultField := Field{
		Package:   pkg.PkgPath,
		Type:      "slice",
		SliceType: t.Elem().String(),
	}

	// If slice type is not accepted type, then process it
	if !IsAcceptedType(resultField.SliceType) {
		s.AddToBeProcessed(pkg.PkgPath, resultField.SliceType)
		return true, resultField
	}

	return false, resultField
}
