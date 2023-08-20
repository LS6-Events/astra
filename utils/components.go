package utils

import (
	"github.com/ls6-events/gengo"
)

// AddComponent adds a component to a slice of components if it doesn't already exist
// It uses the field type and package to determine if the component already exists
func AddComponent(prev []gengo.Field, n ...gengo.Field) []gengo.Field {
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
