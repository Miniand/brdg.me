package sushi_go

import (
	"sort"

	"github.com/Miniand/brdg.me/render"
)

const (
	CardPlayed = iota
	CardTempura
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

var CardStrings = map[int]string{
	CardPlayed:       "played",
	CardTempura:      "tempura",
	CardSashimi:      "sashimi",
	CardDumpling:     "dumpling",
	CardMakiRoll3:    "maki x3",
	CardMakiRoll2:    "maki x2",
	CardMakiRoll1:    "maki x1",
	CardSalmonNigiri: "salmon nigiri",
	CardSquidNigiri:  "squid nigiri",
	CardEggNigiri:    "egg nigiri",
	CardPudding:      "pudding",
	CardWasabi:       "wasabi",
	CardChopsticks:   "chopsticks",
}

var CardColours = map[int]string{
	CardPlayed:       render.Gray,
	CardTempura:      render.Magenta,
	CardSashimi:      render.Magenta,
	CardDumpling:     render.Yellow,
	CardMakiRoll3:    render.Red,
	CardMakiRoll2:    render.Red,
	CardMakiRoll1:    render.Red,
	CardSalmonNigiri: render.Cyan,
	CardSquidNigiri:  render.Cyan,
	CardEggNigiri:    render.Cyan,
	CardPudding:      render.Blue,
	CardWasabi:       render.Green,
	CardChopsticks:   render.Black,
}

var makiExplanation = "most: 6/3"
var CardExplanations = map[int]string{
	CardTempura:      "x2 = 5",
	CardSashimi:      "x3 = 10",
	CardDumpling:     "1 3 6 10 15",
	CardMakiRoll3:    makiExplanation,
	CardMakiRoll2:    makiExplanation,
	CardMakiRoll1:    makiExplanation,
	CardSalmonNigiri: "2",
	CardSquidNigiri:  "3",
	CardEggNigiri:    "1",
	CardPudding:      "end: most 6, least -6",
	CardWasabi:       "next nigiri x3",
	CardChopsticks:   "swap for 2",
}

var CardBaseScores = map[int]int{
	CardSalmonNigiri: 2,
	CardSquidNigiri:  3,
	CardEggNigiri:    1,
}

func CardExplanation(card int, players int) string {
	if card == CardPudding && players > 2 {
		// Special case, no negative points
		return CardExplanations[card] + ", least -6"
	}
	return CardExplanations[card]
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

func TrimPlayed(cards []int) []int {
	newCards := []int{}
	for _, c := range cards {
		if c != CardPlayed {
			newCards = append(newCards, c)
		}
	}
	return newCards
}
