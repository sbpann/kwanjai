package libraries

// Find function to check if element is in string, returns index (int) and  status (bool).
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
