package util

func SliceContains(elements []interface{}, value interface{}) bool {
	for _, a := range elements {
		if a == value {
			return true
		}
	}
	return false
}

func SliceStringContains(elements []string, value string) bool {
	for _, a := range elements {
		if a == value {
			return true
		}
	}
	return false
}
