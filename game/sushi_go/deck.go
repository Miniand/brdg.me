package sushi_go

import "sort"

const (
	CardTempura = iota
	CardSashimi
	CardDumpling
	CardMakiRoll3
	CardMakiRoll2
	CardMakiRoll1
	CardSalmonNigiri
	CardSquidNigiri
	CardEggNigiri
	CardPudding
	CardWasabi
	CardChopsticks
)

var CardTypes = []int{
	CardTempura,
	CardSashimi,
	CardDumpling,
	CardMakiRoll3,
	CardMakiRoll2,
	CardMakiRoll1,
	CardSalmonNigiri,
	CardSquidNigiri,
	CardEggNigiri,
	CardPudding,
	CardWasabi,
	CardChopsticks,
}

var CardCounts = map[int]int{
	CardTempura:      14,
	CardSashimi:      14,
	CardDumpling:     14,
	CardMakiRoll3:    8,
	CardMakiRoll2:    12,
	CardMakiRoll1:    6,
	CardSalmonNigiri: 10,
	CardSquidNigiri:  5,
	CardEggNigiri:    5,
	CardPudding:      10,
	CardWasabi:       6,
	CardChopsticks:   4,
}

var PlayerDrawCounts = map[int]int{
	2: 9, // Usually 10, but we implement the variant.
	3: 9,
	4: 8,
	5: 7,
}

func Deck() []int {
	deck := []int{}
	for _, t := range CardTypes {
		for i := 0; i < CardCounts[t]; i++ {
			deck = append(deck, t)
		}
	}
	return deck
}

func Sort(deck []int) []int {
	sorted := make([]int, len(deck))
	copy(sorted, deck)
	sort.Ints(sorted)
	return sorted
}

func Shuffle(deck []int) []int {
	l := len(deck)
	shuffled := make([]int, l)
	for i, p := range rnd.Perm(l) {
		shuffled[i] = deck[p]
	}
	return shuffled
}

func Contains(needle int, haystack []int) (int, bool) {
	for i, h := range haystack {
		if h == needle {
			return i, true
		}
	}
	return 0, false
}
