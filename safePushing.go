package astra

// AddReturnType adds a return type to a slice of return types if it doesn't already exist.
// It uses the field type, package and status code to determine if the return type already exists.
func AddReturnType(prev []ReturnType, n ...ReturnType) []ReturnType {
	for _, newReturn := range n {
		if len(prev) == 0 {
			prev = append(prev, newReturn)
		} else {
			for _, existingReturn := range prev {
				if newReturn.Field.Type != existingReturn.Field.Type || newReturn.Field.Package != existingReturn.Field.Package || newReturn.StatusCode != existingReturn.StatusCode || newReturn.ContentType != existingReturn.ContentType {
					prev = append(prev, newReturn)
					break
				}
			}
		}
	}

	return prev
}

// AddComponent adds a component to a slice of components if it doesn't already exist.
// It uses the field type and package to determine if the component already exists.
func AddComponent(prev []Field, n ...Field) []Field {
	for _, newComponent := range n {
		if len(prev) == 0 {
			prev = append(prev, newComponent)
		} else {
			var found bool
			for _, existingComponent := range prev {
				if newComponent.Name == existingComponent.Name && newComponent.Package == existingComponent.Package {
					found = true
					break
				}
			}
			if !found {
				prev = append(prev, newComponent)
			}
		}
	}

	return prev
}
