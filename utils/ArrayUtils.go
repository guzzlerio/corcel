package utils

//ContainsString ...
func ContainsString(list []string, expected string) bool {
	for _, current := range list {
		if current == expected {
			return true
		}
	}
	return false
}

//SliceIndex ...
//http://stackoverflow.com/questions/8307478/go-how-to-find-out-element-position-in-slice
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
