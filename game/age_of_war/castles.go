package age_of_war

import "github.com/Miniand/brdg.me/render"

const (
	ClanOda = iota
	ClanTokugawa
	ClanUesugi
	ClanMori
	ClanChosokabe
	ClanShimazu
)

var ClanSetPoints = map[int]int{
	ClanOda: 10,
	ClanTokugawa: 8,
	ClanUesugi: 8,
	ClanMori: 5,
	ClanChosokabe: 4,
	ClanShimazu: 3,
}

var ClanNames = map[int]string{
	ClanOda: "Oda",
	ClanTokugawa: "Tokugawa",
	ClanUesugi: "Uesugi",
	ClanMori: "Mori",
	ClanChosokabe: "Chosokabe",
	ClanShimazu: "Shimazu",
}

var ClanColours = map[int]string{
	ClanOda: render.Yellow,
	ClanTokugawa: render.Gray,
	ClanUesugi: render.Magenta,
	ClanMori: render.Red,
	ClanChosokabe: render.Black,
	ClanShimazu: render.Green,
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
	// Clan Oda
	{
		Clan:   ClanOda,
		Name:   "Azuchi",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceArchery}},
			{Symbols: []int{DiceCavalry, DiceCavalry}},
			{Infantry: 5},
		},
	},
	{
		Clan:   ClanOda,
		Name:   "Matsumoto",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceArchery}},
			{Symbols: []int{DiceArchery}},
			{Infantry: 7},
		},
	},
	{
		Clan:   ClanOda,
		Name:   "Odani",
		Points: 1,
		Lines: []Line{
			{Infantry: 10},
		},
	},
	{
		Clan:   ClanOda,
		Name:   "Gifu",
		Points: 1,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery}},
			{Symbols: []int{DiceCavalry}},
		},
	},

	// Clan Tokugawa
	{
		Clan:   ClanTokugawa,
		Name:   "Edo",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Infantry: 3},
		},
	},
	{
		Clan:   ClanTokugawa,
		Name:   "Kiyosu",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery}},
			{Symbols: []int{DiceCavalry}},			
			{Infantry: 3},
		},
	},
	{
		Clan:   ClanTokugawa,
		Name:   "Inuyama",
		Points: 1,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery, DiceArchery}},
		},
	},


	// Clan Uesugi
	{
		Clan:   ClanUesugi,
		Name:   "Kasugayama",
		Points: 4,
		Lines: []Line{
			{Symbols: []int{DiceArchery, DiceArchery}},
			{Symbols: []int{DiceCavalry, DiceCavalry}},
		},
	},
		{
		Clan:   ClanUesugi,
		Name:   "Kitanosho",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Symbols: []int{DiceArchery, DiceCavalry}},
			{Infantry: 6},
		},
	},
	// Clan Mori
	{
		Clan:   ClanMori,
		Name:   "Gassantoda",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Infantry: 8},
		},
	},
		{
		Clan:   ClanMori,
		Name:   "Takahashi",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceCavalry, DiceCavalry}},
			{Infantry: 5},
			{Infantry: 2},			
			
		},
	},
	// Clan Chosokabe
	{
		Clan:   ClanChosokabe,
		Name:   "Matsuyama",
		Points: 2,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo}},
			{Infantry: 4},
			{Infantry: 4},
		},
	},
		{
		Clan:   ClanChosokabe,
		Name:   "Marugame",
		Points: 1,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo, DiceDaimyo}},
			{Symbols: []int{DiceCavalry}},
		},
	},
	// Clan Shimazu
	{
		Clan:   ClanShimazu,
		Name:   "Kumamoto",
		Points: 3,
		Lines: []Line{
			{Symbols: []int{DiceDaimyo, DiceDaimyo}},
			{Symbols: []int{DiceCavalry}},
			{Symbols: []int{DiceArchery}},
			{Infantry: 4},
		},
	},
}
