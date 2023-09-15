package astra

// strSliceContains checks if a string slice contains a value
// Simple utility function
func strSliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}
