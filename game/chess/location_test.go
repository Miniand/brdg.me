package chess

func inLocations(needle Location, haystack []Location) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}
	return false
}
