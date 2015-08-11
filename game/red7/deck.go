package red7

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

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

var RankValMap map[int]int

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

var SuitAbbrMap map[string]int

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

	SuitAbbrMap = map[string]int{}
	for k, v := range SuitAbbr {
		SuitAbbrMap[v] = k
	}

	RankValMap = map[int]int{}
	for k, v := range RankVal {
		RankValMap[v] = k
	}
}

func ParseCard(input string) (card int, ok bool) {
	if len(input) != 2 {
		return
	}
	suit, ok := SuitAbbrMap[strings.ToUpper(input[0:1])]
	if !ok {
		return
	}
	rank, err := strconv.Atoi(input[1:2])
	if err != nil || rank < 1 || rank > 7 {
		ok = false
		return
	}
	card = suit | RankValMap[rank]
	return
}

func Points(cards []int) int {
	points := 0
	for _, c := range cards {
		points += RankVal[c&RankMask]
	}
	return points
}

type BySuit []int

func (bs BySuit) Len() int      { return len(bs) }
func (bs BySuit) Swap(i, j int) { bs[i], bs[j] = bs[j], bs[i] }
func (bs BySuit) Less(i, j int) bool {
	return bs[i]&SuitMask < bs[j]&SuitMask ||
		bs[i]&SuitMask == bs[j]&SuitMask &&
			bs[i]&RankMask < bs[j]&RankMask
}

func SortBySuit(cards []int) []int {
	sortedCards := append([]int{}, cards...)
	sort.Sort(BySuit(sortedCards))
	return sortedCards
}
