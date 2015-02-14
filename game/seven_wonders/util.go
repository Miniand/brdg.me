package seven_wonders

import "github.com/Miniand/brdg.me/game/cost"

func TrimIntMap(m map[int]int) map[int]int {
	n := map[int]int{}
	for k, v := range m {
		if v != 0 {
			n[k] = v
		}
	}
	return n
}

func InInts(needle int, haystack []int) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}
	return false
}

func SliceToCost(ints []int) []cost.Cost {
	l := len(ints)
	c := make([]cost.Cost, l)
	for i := 0; i < l; i++ {
		c[i] = cost.Cost{
			ints[i]: 1,
		}
	}
	return c
}
