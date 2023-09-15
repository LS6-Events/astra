package utils

import (
	"github.com/ls6-events/astra"
)

// AddReturnType adds a return type to a slice of return types if it doesn't already exist
// It uses the field type, package and status code to determine if the return type already exists
func AddReturnType(prev []astra.ReturnType, n ...astra.ReturnType) []astra.ReturnType {
	for _, newReturn := range n {
		if len(prev) == 0 {
			prev = append(prev, newReturn)
		} else {
			for _, existingReturn := range prev {
				if newReturn.Field.Type != existingReturn.Field.Type || newReturn.Field.Package != existingReturn.Field.Package || newReturn.StatusCode != existingReturn.StatusCode {
					prev = append(prev, newReturn)
					break
				}
			}
		}
	}

	return prev
}
