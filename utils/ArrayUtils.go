package utils

func containsString(list []string, expected string) bool {
	for _, current := range list {
		if current == expected {
			return true
		}
	}
	return false
}
