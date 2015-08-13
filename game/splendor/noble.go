package splendor

import (
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/game/cost"
)

type Noble struct {
	Prestige int
	Cost     cost.Cost
}

func ShuffleNobles(nobles []Noble) []Noble {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(nobles)
	shuffled := make([]Noble, l)
	for i, n := range r.Perm(l) {
		shuffled[i] = nobles[n]
	}
	return shuffled
}

func NobleCards() []Noble {
	return []Noble{
		{
			3,
			cost.Cost{
				Emerald:  3,
				Sapphire: 3,
				Diamond:  3,
			},
		},
		{
			3,
			cost.Cost{
				Emerald:  3,
				Sapphire: 3,
				Ruby:     3,
			},
		},
		{
			3,
			cost.Cost{
				Onyx:    3,
				Ruby:    3,
				Diamond: 3,
			},
		},
		{
			3,
			cost.Cost{
				Onyx:     3,
				Sapphire: 3,
				Diamond:  3,
			},
		},
		{
			3,
			cost.Cost{
				Onyx:    3,
				Ruby:    3,
				Emerald: 3,
			},
		},
		{
			3,
			cost.Cost{
				Onyx: 4,
				Ruby: 4,
			},
		},
		{
			3,
			cost.Cost{
				Onyx:    4,
				Diamond: 4,
			},
		},
		{
			3,
			cost.Cost{
				Sapphire: 4,
				Diamond:  4,
			},
		},
		{
			3,
			cost.Cost{
				Sapphire: 4,
				Emerald:  4,
			},
		},
		{
			3,
			cost.Cost{
				Ruby:    4,
				Emerald: 4,
			},
		},
	}
}
