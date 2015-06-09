package age_of_war

const (
	ClanFarty = iota
	ClanEgg
	ClanBacon
)

var ClanSetPoints = map[int]int{
	ClanFarty: 7,
	ClanEgg:   8,
	ClanBacon: 9,
}

var ClanNames = map[int]string{
	ClanFarty: "Farty",
	ClanEgg:   "Egg",
	ClanBacon: "Bacon",
}

type Castle struct {
	Clan   int
	Name   string
	Points int
	// Lines are from top to bottom on the card, not including the special Daimyo for stealing.
	Lines []Line
}

type Line struct {
	Infantry int
	Symbols  []int
}

// Definitions of all the castles.
var Castles = []Castle{
	{
		Clan:   ClanFarty,
		Name:   "Fart Castle",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Infantry: 6},
		},
	},
	{
		Clan:   ClanEgg,
		Name:   "Egg Castle",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo, DiceCavalry}},
			{Symbols: []int{DiceArchery, DiceArchery}},
			{Infantry: 3},
		},
	},
}
