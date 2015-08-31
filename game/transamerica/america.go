package transamerica

import "github.com/Miniand/brdg.me/render"

var America = &Board{
	Nodes:   map[Loc]*Node{},
	Terrain: map[Edge]int{},
}

func NodeRow(y, fromX, untilX byte) []Loc {
	locs := []Loc{}
	for x := fromX; x <= untilX; x++ {
		locs = append(locs, Loc{x, y})
	}
	return locs
}

func init() {
	// America
	for _, row := range [][]Loc{
		NodeRow('A', 'B', 'K'),
		NodeRow('B', 'A', 'L'),
		NodeRow('B', 'P', 'Q'),
		NodeRow('C', 'A', 'M'),
		NodeRow('C', 'P', 'R'),
		NodeRow('D', 'A', 'Q'),
		NodeRow('E', 'A', 'Q'),
		NodeRow('F', 'A', 'P'),
		NodeRow('G', 'A', 'P'),
		NodeRow('H', 'A', 'P'),
		NodeRow('I', 'B', 'Q'),
		NodeRow('J', 'B', 'P'),
		NodeRow('K', 'C', 'P'),
		NodeRow('L', 'D', 'O'),
		NodeRow('M', 'G', 'O'),
	} {
		for _, l := range row {
			America.Nodes[l] = &Node{}
		}
	}

	for loc, city := range map[string]string{
		"AB": render.Green,
		"BA": render.Green,
		"DA": render.Green,
		"FA": render.Green,
		"GA": render.Green,
		"JB": render.Green,
		"KC": render.Green,

		"BD": render.Blue,
		"BH": render.Blue,
		"BK": render.Blue,
		"CK": render.Blue,
		"DM": render.Blue,
		"CP": render.Blue,
		"FN": render.Blue,

		"ED": render.Yellow,
		"FF": render.Yellow,
		"EI": render.Yellow,
		"GJ": render.Yellow,
		"GL": render.Yellow,
		"II": render.Yellow,
		"IF": render.Yellow,

		"JD": render.Red,
		"LF": render.Red,
		"KJ": render.Red,
		"MJ": render.Red,
		"JL": render.Red,
		"ML": render.Red,
		"KN": render.Red,

		"CR": render.Magenta,
		"EQ": render.Magenta,
		"FP": render.Magenta,
		"HP": render.Magenta,
		"IO": render.Magenta,
		"KP": render.Magenta,
		"MO": render.Magenta,
	} {
		America.Nodes[MustParseLoc(loc)].City = city
	}

	for e, t := range map[string]int{
		"AB AC": TerrainMountain,
	} {
		America.Terrain[MustParseEdge(e)] = t
	}
}
