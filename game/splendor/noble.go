package splendor

import (
	"math/rand"
	"time"
)

type Noble struct {
	Prestige int
	Cost     Amount
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
			Amount{
				Emerald:  3,
				Sapphire: 3,
				Diamond:  3,
			},
		},
		{
			3,
			Amount{
				Emerald:  3,
				Sapphire: 3,
				Ruby:     3,
			},
		},
		{
			3,
			Amount{
				Onyx:    3,
				Ruby:    3,
				Diamond: 3,
			},
		},
		{
			3,
			Amount{
				Onyx:     3,
				Sapphire: 3,
				Diamond:  3,
			},
		},
		{
			3,
			Amount{
				Onyx:    3,
				Ruby:    3,
				Emerald: 3,
			},
		},
		{
			3,
			Amount{
				Onyx: 4,
				Ruby: 4,
			},
		},
		{
			3,
			Amount{
				Onyx:    4,
				Diamond: 4,
			},
		},
		{
			3,
			Amount{
				Sapphire: 4,
				Diamond:  4,
			},
		},
		{
			3,
			Amount{
				Sapphire: 4,
				Emerald:  4,
			},
		},
		{
			3,
			Amount{
				Ruby:    4,
				Emerald: 4,
			},
		},
	}
}
