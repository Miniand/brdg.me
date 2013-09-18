package chess

func inLocations(needle Location, haystack []Move) bool {
	for _, h := range haystack {
		if needle == h.To {
			return true
		}
	}
	return false
}
