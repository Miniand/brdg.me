package red7

import "github.com/Miniand/brdg.me/game/helper"

var SuitRules = map[int]func([]int) []int{
	SuitRed:    HighestCard,
	SuitOrange: CardsOfOneNumber,
	SuitYellow: CardsOfOneColor,
	SuitGreen:  MostEvenCards,
	SuitBlue:   CardsOfDifferentColors,
	SuitIndigo: CardsThatFormARun,
	SuitViolet: MostCardsBelow4,
}

var SuitRulesStrs = map[int]string{
	SuitRed:    "Highest card",
	SuitOrange: "Same number",
	SuitYellow: "Same color",
	SuitGreen:  "Even cards",
	SuitBlue:   "Most colors",
	SuitIndigo: "In a row",
	SuitViolet: "Below 4",
}

func HighestCard(cards []int) []int {
	return []int{helper.IntMax(cards...)}
}

func CardsOfOneNumber(cards []int) []int {
	curRank := 0
	most := []int{}
	cur := []int{}
	for _, c := range helper.IntReverse(helper.IntSort(cards)) {
		rank := c & RankMask
		if rank != curRank {
			if len(cur) > len(most) {
				most = cur
			}
			cur = []int{c}
			curRank = rank
		} else {
			cur = append(cur, c)
		}
	}
	if len(cur) > len(most) {
		most = cur
	}
	return most
}

func CardsOfOneColor(cards []int) []int {
	bySuit := [][]int{}
	suitMap := map[int]int{}
	n := 0
	for _, c := range helper.IntReverse(helper.IntSort(cards)) {
		suit := c & SuitMask
		if _, ok := suitMap[suit]; !ok {
			suitMap[suit] = len(bySuit)
			bySuit = append(bySuit, []int{})
		}
		index := suitMap[suit]
		bySuit[index] = append(bySuit[index], c)
		if l := len(bySuit[index]); l > n {
			n = l
		}
	}
	for _, s := range bySuit {
		if len(s) == n {
			return s
		}
	}
	return []int{}
}

func MostEvenCards(cards []int) []int {
	even := []int{}
	for _, c := range cards {
		if RankVal[c&RankMask]%2 == 0 {
			even = append(even, c)
		}
	}
	return even
}

func CardsOfDifferentColors(cards []int) []int {
	usedSuits := map[int]bool{}
	matching := []int{}
	for _, c := range helper.IntReverse(helper.IntSort(cards)) {
		suit := c & SuitMask
		if !usedSuits[suit] {
			usedSuits[suit] = true
			matching = append(matching, c)
		}
	}
	return matching
}

func CardsThatFormARun(cards []int) []int {
	lastRank := 0
	cur := []int{}
	longest := []int{}
	for _, c := range helper.IntReverse(helper.IntSort(cards)) {
		rank := RankVal[c&RankMask]
		switch rank {
		case lastRank:
			continue
		case lastRank - 1:
			cur = append(cur, c)
		default:
			if len(cur) > len(longest) {
				longest = cur
			}
			cur = []int{c}
		}
		lastRank = rank
	}
	if len(cur) > len(longest) {
		longest = cur
	}
	return longest
}

func MostCardsBelow4(cards []int) []int {
	below := []int{}
	for _, c := range helper.IntReverse(helper.IntSort(cards)) {
		if RankVal[c&RankMask] < 4 {
			below = append(below, c)
		}
	}
	return below
}

func Leader(palettes [][]int) (leader int, leaderPalette []int) {
	l := len(palettes)
	if l == 0 {
		return
	}
	leaderPalette = palettes[0]
	for i := 1; i < l; i++ {
		lLen := len(leaderPalette)
		iLen := len(palettes[i])
		if iLen > lLen || iLen == lLen &&
			helper.IntMax(palettes[i]...) > helper.IntMax(leaderPalette...) {
			leader = i
			leaderPalette = palettes[i]
		}
	}
	return
}
