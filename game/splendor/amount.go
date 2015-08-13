package splendor

import "github.com/Miniand/brdg.me/game/cost"

func CanAfford(a, c cost.Cost) bool {
	short := 0
	for g, n := range c {
		if a[g] < n {
			short += n - a[g]
		}
	}
	return a[Gold]-c[Gold] >= short
}
