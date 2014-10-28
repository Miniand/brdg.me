package category_5

import (
	"sort"

	"github.com/Miniand/brdg.me/render"
)

type Card int

var CardColours = map[int]string{
	7: "magenta",
	5: "red",
	3: "yellow",
	2: "cyan",
	1: "gray",
}

func (c Card) Heads() int {
	switch {
	case c == 55:
		return 7
	case c%11 == 0:
		return 5
	case c%10 == 0:
		return 3
	case c%5 == 0:
		return 2
	default:
		return 1
	}
}

func (c Card) String() string {
	return render.Markup(int(c), CardColours[c.Heads()], true)
}

func Deck() []Card {
	deck := []Card{}
	for i := 1; i <= 104; i++ {
		deck = append(deck, Card(i))
	}
	return deck
}

func Shuffle(deck []Card) []Card {
	l := len(deck)
	shuffled := make([]Card, l)
	for i, p := range r.Perm(l) {
		shuffled[i] = deck[p]
	}
	return shuffled
}

func TakeCards(deck []Card, n int) (taken, remaining []Card) {
	if len(deck) < n {
		panic("not enough cards in deck")
	}
	return deck[:n], deck[n:]
}

func RemoveCard(deck []Card, card Card) ([]Card, bool) {
	for i, c := range deck {
		if c == card {
			return append(deck[:i], deck[i+1:]...), true
		}
	}
	return deck, false
}

func Ctoi(cards []Card) []int {
	l := len(cards)
	cardInts := make([]int, l)
	for i, c := range cards {
		cardInts[i] = int(c)
	}
	return cardInts
}

func Itoc(ints []int) []Card {
	l := len(ints)
	cards := make([]Card, l)
	for i, in := range ints {
		cards[i] = Card(in)
	}
	return cards
}

func SortCards(cards []Card) []Card {
	ints := Ctoi(cards)
	sort.Ints(ints)
	return Itoc(ints)
}

func CardsHeads(cards []Card) int {
	heads := 0
	for _, c := range cards {
		heads += c.Heads()
	}
	return heads
}
