package utils

import (
	"github.com/ls6-events/gengo"
)

func AddReturnType(prev []gengo.ReturnType, n ...gengo.ReturnType) []gengo.ReturnType {
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
