package splendor

type Noble struct {
	Prestige int
	Cost     Cost
}

func NobleCards() []Noble {
	return []Noble{
		{
			3,
			Cost{
				Emerald:  3,
				Sapphire: 3,
				Diamond:  3,
			},
		},
		{
			3,
			Cost{
				Emerald:  3,
				Sapphire: 3,
				Ruby:     3,
			},
		},
		{
			3,
			Cost{
				Onyx:    3,
				Ruby:    3,
				Diamond: 3,
			},
		},
		{
			3,
			Cost{
				Onyx:     3,
				Sapphire: 3,
				Diamond:  3,
			},
		},
		{
			3,
			Cost{
				Onyx:    3,
				Ruby:    3,
				Emerald: 3,
			},
		},
		{
			3,
			Cost{
				Onyx: 4,
				Ruby: 4,
			},
		},
		{
			3,
			Cost{
				Onyx:    4,
				Diamond: 4,
			},
		},
		{
			3,
			Cost{
				Sapphire: 4,
				Diamond:  4,
			},
		},
		{
			3,
			Cost{
				Sapphire: 4,
				Emerald:  4,
			},
		},
		{
			3,
			Cost{
				Ruby:    4,
				Emerald: 4,
			},
		},
	}
}
