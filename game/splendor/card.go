package splendor

import (
	"math/rand"
	"time"
)

const (
	Diamond = iota
	Sapphire
	Emerald
	Ruby
	Onyx
	Gold
	Prestige
)

type Amount map[int]int

type Card struct {
	Resource int
	Prestige int
	Cost     Amount
}

var Resources = []int{
	Diamond,
	Sapphire,
	Emerald,
	Ruby,
	Onyx,
	Gold,
	Prestige,
}

var Gems = []int{
	Diamond,
	Sapphire,
	Emerald,
	Ruby,
	Onyx,
}

func ShuffleCards(cards []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(cards)
	shuffled := make([]Card, l)
	for i, c := range r.Perm(l) {
		shuffled[i] = cards[c]
	}
	return shuffled
}

func Level1Cards() []Card {
	return []Card{
		{
			Diamond,
			0,
			Amount{
				Sapphire: 1,
				Emerald:  1,
				Ruby:     1,
				Onyx:     1,
			},
		},
		{
			Diamond,
			0,
			Amount{
				Sapphire: 1,
				Emerald:  2,
				Ruby:     1,
				Onyx:     1,
			},
		},
		{
			Diamond,
			0,
			Amount{
				Diamond:  3,
				Sapphire: 1,
				Onyx:     1,
			},
		},
		{
			Diamond,
			0,
			Amount{
				Sapphire: 2,
				Emerald:  2,
				Onyx:     1,
			},
		},
		{
			Diamond,
			0,
			Amount{
				Sapphire: 2,
				Onyx:     2,
			},
		},
		{
			Diamond,
			0,
			Amount{
				Ruby: 2,
				Onyx: 1,
			},
		},
		{
			Diamond,
			1,
			Amount{
				Emerald: 4,
			},
		},
		{
			Diamond,
			0,
			Amount{
				Sapphire: 3,
			},
		},
		{
			Sapphire,
			0,
			Amount{
				Diamond: 1,
				Emerald: 1,
				Ruby:    1,
				Onyx:    1,
			},
		},
		{
			Sapphire,
			0,
			Amount{
				Diamond: 1,
				Emerald: 1,
				Ruby:    2,
				Onyx:    1,
			},
		},
		{
			Sapphire,
			0,
			Amount{
				Sapphire: 1,
				Emerald:  3,
				Ruby:     1,
			},
		},
		{
			Sapphire,
			0,
			Amount{
				Diamond: 1,
				Emerald: 2,
				Ruby:    2,
			},
		},
		{
			Sapphire,
			0,
			Amount{
				Emerald: 2,
				Onyx:    2,
			},
		},
		{
			Sapphire,
			0,
			Amount{
				Diamond: 1,
				Onyx:    2,
			},
		},
		{
			Sapphire,
			1,
			Amount{
				Ruby: 4,
			},
		},
		{
			Sapphire,
			0,
			Amount{
				Onyx: 3,
			},
		},
		{
			Onyx,
			0,
			Amount{
				Diamond:  1,
				Sapphire: 2,
				Emerald:  1,
				Ruby:     1,
			},
		},
		{
			Onyx,
			0,
			Amount{
				Diamond:  1,
				Sapphire: 1,
				Emerald:  1,
				Ruby:     1,
			},
		},
		{
			Onyx,
			0,
			Amount{
				Emerald: 1,
				Ruby:    3,
				Onyx:    1,
			},
		},
		{
			Onyx,
			0,
			Amount{
				Diamond:  2,
				Sapphire: 2,
				Ruby:     1,
			},
		},
		{
			Onyx,
			0,
			Amount{
				Diamond: 2,
				Emerald: 2,
			},
		},
		{
			Onyx,
			0,
			Amount{
				Emerald: 2,
				Ruby:    1,
			},
		},
		{
			Onyx,
			0,
			Amount{
				Emerald: 3,
			},
		},
		{
			Onyx,
			1,
			Amount{
				Sapphire: 4,
			},
		},
		{
			Ruby,
			0,
			Amount{
				Diamond:  2,
				Sapphire: 1,
				Emerald:  1,
				Onyx:     1,
			},
		},
		{
			Ruby,
			0,
			Amount{
				Diamond:  1,
				Sapphire: 1,
				Emerald:  1,
				Onyx:     1,
			},
		},
		{
			Ruby,
			0,
			Amount{
				Diamond: 1,
				Ruby:    1,
				Onyx:    3,
			},
		},
		{
			Ruby,
			0,
			Amount{
				Diamond: 2,
				Ruby:    2,
			},
		},
		{
			Ruby,
			0,
			Amount{
				Sapphire: 2,
				Emerald:  1,
			},
		},
		{
			Ruby,
			1,
			Amount{
				Diamond: 4,
			},
		},
		{
			Ruby,
			0,
			Amount{
				Diamond: 3,
			},
		},
		{
			Ruby,
			0,
			Amount{
				Diamond: 2,
				Emerald: 1,
				Onyx:    2,
			},
		},
		{
			Emerald,
			0,
			Amount{
				Diamond:  1,
				Sapphire: 1,
				Ruby:     1,
				Onyx:     1,
			},
		},
		{
			Emerald,
			0,
			Amount{
				Diamond:  1,
				Sapphire: 1,
				Ruby:     1,
				Onyx:     2,
			},
		},
		{
			Emerald,
			0,
			Amount{
				Diamond:  1,
				Sapphire: 3,
				Emerald:  1,
			},
		},
		{
			Emerald,
			0,
			Amount{
				Sapphire: 1,
				Ruby:     2,
				Onyx:     2,
			},
		},
		{
			Emerald,
			0,
			Amount{
				Sapphire: 2,
				Ruby:     2,
			},
		},
		{
			Emerald,
			0,
			Amount{
				Diamond:  2,
				Sapphire: 1,
			},
		},
		{
			Emerald,
			0,
			Amount{
				Ruby: 3,
			},
		},
		{
			Emerald,
			1,
			Amount{
				Onyx: 4,
			},
		},
	}
}

func Level2Cards() []Card {
	return []Card{
		{
			Diamond,
			1,
			Amount{
				Emerald: 3,
				Ruby:    2,
				Onyx:    2,
			},
		},
		{
			Diamond,
			1,
			Amount{
				Diamond:  2,
				Sapphire: 3,
				Ruby:     3,
			},
		},
		{
			Diamond,
			2,
			Amount{
				Emerald: 1,
				Ruby:    4,
				Onyx:    2,
			},
		},
		{
			Diamond,
			2,
			Amount{
				Ruby: 5,
				Onyx: 3,
			},
		},
		{
			Diamond,
			2,
			Amount{
				Ruby: 5,
			},
		},
		{
			Diamond,
			3,
			Amount{
				Diamond: 6,
			},
		},
		{
			Sapphire,
			1,
			Amount{
				Sapphire: 2,
				Emerald:  2,
				Ruby:     3,
			},
		},
		{
			Sapphire,
			1,
			Amount{
				Sapphire: 2,
				Emerald:  3,
				Onyx:     3,
			},
		},
		{
			Sapphire,
			2,
			Amount{
				Diamond:  5,
				Sapphire: 3,
			},
		},
		{
			Sapphire,
			2,
			Amount{
				Diamond: 2,
				Ruby:    1,
				Onyx:    4,
			},
		},
		{
			Sapphire,
			3,
			Amount{
				Sapphire: 6,
			},
		},
		{
			Sapphire,
			2,
			Amount{
				Sapphire: 5,
			},
		},
		{
			Onyx,
			1,
			Amount{
				Diamond: 3,
				Emerald: 3,
				Onyx:    2,
			},
		},
		{
			Onyx,
			2,
			Amount{
				Sapphire: 1,
				Emerald:  4,
				Ruby:     2,
			},
		},
		{
			Onyx,
			1,
			Amount{
				Diamond:  3,
				Sapphire: 2,
				Emerald:  2,
			},
		},
		{
			Onyx,
			2,
			Amount{
				Emerald: 5,
				Ruby:    3,
			},
		},
		{
			Onyx,
			2,
			Amount{
				Diamond: 5,
			},
		},
		{
			Onyx,
			3,
			Amount{
				Onyx: 6,
			},
		},
		{
			Ruby,
			1,
			Amount{
				Diamond: 2,
				Ruby:    2,
				Onyx:    3,
			},
		},
		{
			Ruby,
			1,
			Amount{
				Sapphire: 3,
				Ruby:     2,
				Onyx:     3,
			},
		},
		{
			Ruby,
			2,
			Amount{
				Diamond:  1,
				Sapphire: 4,
				Emerald:  2,
			},
		},
		{
			Ruby,
			2,
			Amount{
				Diamond: 3,
				Onyx:    5,
			},
		},
		{
			Ruby,
			2,
			Amount{
				Onyx: 5,
			},
		},
		{
			Ruby,
			3,
			Amount{
				Ruby: 6,
			},
		},
		{
			Emerald,
			2,
			Amount{
				Emerald: 5,
			},
		},
		{
			Emerald,
			2,
			Amount{
				Sapphire: 5,
				Emerald:  3,
			},
		},
		{
			Emerald,
			3,
			Amount{
				Emerald: 6,
			},
		},
		{
			Emerald,
			1,
			Amount{
				Diamond:  2,
				Sapphire: 3,
				Onyx:     2,
			},
		},
		{
			Emerald,
			1,
			Amount{
				Diamond: 3,
				Emerald: 2,
				Ruby:    3,
			},
		},
		{
			Emerald,
			2,
			Amount{
				Diamond:  4,
				Sapphire: 2,
				Onyx:     1,
			},
		},
	}
}

func Level3Cards() []Card {
	return []Card{
		{
			Diamond,
			4,
			Amount{
				Diamond: 3,
				Ruby:    3,
				Onyx:    6,
			},
		},
		{
			Diamond,
			4,
			Amount{
				Onyx: 7,
			},
		},
		{
			Diamond,
			5,
			Amount{
				Diamond: 3,
				Onyx:    7,
			},
		},
		{
			Diamond,
			3,
			Amount{
				Sapphire: 3,
				Emerald:  3,
				Ruby:     5,
				Onyx:     3,
			},
		},
		{
			Sapphire,
			5,
			Amount{
				Diamond:  7,
				Sapphire: 3,
			},
		},
		{
			Sapphire,
			4,
			Amount{
				Diamond:  6,
				Sapphire: 3,
				Onyx:     3,
			},
		},
		{
			Sapphire,
			3,
			Amount{
				Diamond: 3,
				Emerald: 3,
				Ruby:    3,
				Onyx:    5,
			},
		},
		{
			Sapphire,
			4,
			Amount{
				Diamond: 7,
			},
		},
		{
			Onyx,
			4,
			Amount{
				Ruby: 7,
			},
		},
		{
			Onyx,
			4,
			Amount{
				Emerald: 3,
				Ruby:    6,
				Onyx:    3,
			},
		},
		{
			Onyx,
			5,
			Amount{
				Ruby: 7,
				Onyx: 3,
			},
		},
		{
			Onyx,
			3,
			Amount{
				Diamond:  3,
				Sapphire: 3,
				Emerald:  5,
				Ruby:     3,
			},
		},
		{
			Ruby,
			4,
			Amount{
				Emerald: 7,
			},
		},
		{
			Ruby,
			3,
			Amount{
				Diamond:  3,
				Sapphire: 5,
				Emerald:  3,
				Onyx:     3,
			},
		},
		{
			Ruby,
			5,
			Amount{
				Emerald: 7,
				Ruby:    3,
			},
		},
		{
			Ruby,
			4,
			Amount{
				Sapphire: 3,
				Emerald:  6,
				Ruby:     3,
			},
		},
		{
			Emerald,
			4,
			Amount{
				Diamond:  3,
				Sapphire: 6,
				Emerald:  3,
			},
		},
		{
			Emerald,
			4,
			Amount{
				Sapphire: 7,
			},
		},
		{
			Emerald,
			5,
			Amount{
				Sapphire: 7,
				Emerald:  3,
			},
		},
		{
			Emerald,
			3,
			Amount{
				Diamond:  5,
				Sapphire: 3,
				Ruby:     3,
				Onyx:     3,
			},
		},
	}
}
