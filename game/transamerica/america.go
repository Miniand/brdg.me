package transamerica

import "github.com/Miniand/brdg.me/render"

var America = map[Loc]*Node{}

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
			America[l] = &Node{}
		}
	}

	America[MustParseLoc("AB")].City = render.Green
	America[MustParseLoc("BA")].City = render.Green
	America[MustParseLoc("BA")].City = render.Green
	America[MustParseLoc("DA")].City = render.Green
	America[MustParseLoc("FA")].City = render.Green
	America[MustParseLoc("GA")].City = render.Green
	America[MustParseLoc("JB")].City = render.Green
	America[MustParseLoc("KC")].City = render.Green

	America[MustParseLoc("BD")].City = render.Blue
	America[MustParseLoc("BH")].City = render.Blue
	America[MustParseLoc("BK")].City = render.Blue
	America[MustParseLoc("CK")].City = render.Blue
	America[MustParseLoc("DM")].City = render.Blue
	America[MustParseLoc("CP")].City = render.Blue
	America[MustParseLoc("FN")].City = render.Blue

	America[MustParseLoc("ED")].City = render.Yellow
	America[MustParseLoc("FF")].City = render.Yellow
	America[MustParseLoc("EI")].City = render.Yellow
	America[MustParseLoc("GJ")].City = render.Yellow
	America[MustParseLoc("GL")].City = render.Yellow
	America[MustParseLoc("II")].City = render.Yellow
	America[MustParseLoc("IF")].City = render.Yellow

	America[MustParseLoc("JD")].City = render.Red
	America[MustParseLoc("LF")].City = render.Red
	America[MustParseLoc("KJ")].City = render.Red
	America[MustParseLoc("MJ")].City = render.Red
	America[MustParseLoc("JL")].City = render.Red
	America[MustParseLoc("ML")].City = render.Red
	America[MustParseLoc("KN")].City = render.Red

	America[MustParseLoc("CR")].City = render.Magenta
	America[MustParseLoc("EQ")].City = render.Magenta
	America[MustParseLoc("FP")].City = render.Magenta
	America[MustParseLoc("HP")].City = render.Magenta
	America[MustParseLoc("IO")].City = render.Magenta
	America[MustParseLoc("KP")].City = render.Magenta
	America[MustParseLoc("MO")].City = render.Magenta
}
