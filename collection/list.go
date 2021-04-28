package collection

// ListContains returns true if the given list contains the given element
func ListContains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}

	return false
}
