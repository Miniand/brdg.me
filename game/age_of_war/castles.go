package age_of_war

import "github.com/Miniand/brdg.me/render"

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

var ClanColours = map[int]string{
	ClanFarty: render.Blue,
	ClanEgg:   render.Red,
	ClanBacon: render.Yellow,
}

type Castle struct {
	Clan   int
	Name   string
	Points int
	// Lines are from top to bottom on the card, not including the special Daimyo for stealing.
	Lines []Line
}

// MinDice is the minimum dice required to conquer this castle.
func (c Castle) MinDice() int {
	min := 0
	for _, l := range c.Lines {
		min += len(l.Symbols) + (l.Infantry+2)/3
	}
	return min
}

// CalcLines gets the lines for the castle, including the extra daimyo if
// stealing.
func (c Castle) CalcLines(stealing bool) []Line {
	lines := []Line{}
	if c.Lines != nil {
		lines = append(lines, c.Lines...)
	}
	if stealing {
		lines = append(lines, Line{
			Symbols: []int{DiceDaimyo},
		})
	}
	return lines
}

type Line struct {
	Infantry int
	Symbols  []int
}

// Definitions of all the castles.
var Castles = []Castle{
	// Clan Bacon
	{
		Clan:   ClanBacon,
		Name:   "Hamfort",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Infantry: 6},
		},
	},
	{
		Clan:   ClanBacon,
		Name:   "Pig keep",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Infantry: 6},
		},
	},
	{
		Clan:   ClanBacon,
		Name:   "Meat tower",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo, DiceCavalry}},
			{Symbols: []int{DiceArchery, DiceArchery}},
			{Infantry: 3},
		},
	},

	// Clan Farty
	{
		Clan:   ClanFarty,
		Name:   "Fartfort",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Infantry: 6},
		},
	},
	{
		Clan:   ClanFarty,
		Name:   "Stinky keep",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Infantry: 6},
		},
	},

	// Clan Egg
	{
		Clan:   ClanEgg,
		Name:   "Egg town",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo, DiceCavalry}},
			{Symbols: []int{DiceArchery, DiceArchery}},
			{Infantry: 3},
		},
	},
}
