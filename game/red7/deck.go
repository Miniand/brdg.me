package red7

import "github.com/Miniand/brdg.me/render"

const (
	SuitViolet = 1 << iota
	SuitIndigo
	SuitBlue
	SuitGreen
	SuitYellow
	SuitOrange
	SuitRed

	Rank1
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
)

const (
	SuitMask = SuitViolet | SuitIndigo | SuitBlue | SuitGreen | SuitYellow | SuitOrange | SuitRed
	RankMask = Rank1 | Rank2 | Rank3 | Rank4 | Rank5 | Rank6 | Rank7
)

var Suits = []int{
	SuitViolet,
	SuitIndigo,
	SuitBlue,
	SuitGreen,
	SuitYellow,
	SuitOrange,
	SuitRed,
}

var Ranks = []int{
	Rank1,
	Rank2,
	Rank3,
	Rank4,
	Rank5,
	Rank6,
	Rank7,
}

var RankVal = map[int]int{
	Rank1: 1,
	Rank2: 2,
	Rank3: 3,
	Rank4: 4,
	Rank5: 5,
	Rank6: 6,
	Rank7: 7,
}

var SuitStr = map[int]string{
	SuitViolet: "Violet",
	SuitIndigo: "Indigo",
	SuitBlue:   "Blue",
	SuitGreen:  "Green",
	SuitYellow: "Yellow",
	SuitOrange: "Orange",
	SuitRed:    "Red",
}

var SuitAbbr = map[int]string{
	SuitViolet: "V",
	SuitIndigo: "I",
	SuitBlue:   "B",
	SuitGreen:  "G",
	SuitYellow: "Y",
	SuitOrange: "O",
	SuitRed:    "R",
}

var SuitCol = map[int]string{
	SuitViolet: render.Magenta,
	SuitIndigo: render.Blue,
	SuitBlue:   render.Cyan,
	SuitGreen:  render.Green,
	SuitYellow: render.Yellow,
	SuitOrange: render.Gray,
	SuitRed:    render.Red,
}

func CardValues(card int) (suit, rank int) {
	return card & SuitMask, card & RankMask
}

var Deck []int

func init() {
	Deck = make([]int, 0, len(Suits)*len(Ranks))
	for _, s := range Suits {
		for _, r := range Ranks {
			Deck = append(Deck, s|r)
		}
	}
}

func ParseCard(input string) (int, bool) {
	return 0, false
}
