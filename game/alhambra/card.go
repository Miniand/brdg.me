package alhambra

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/render"
)

const (
	CurrencyBlue = iota
	CurrencyGreen
	CurrencyRed
	CurrencyYellow
)

var Currencies = []int{
	CurrencyBlue,
	CurrencyGreen,
	CurrencyRed,
	CurrencyYellow,
}

var CurrencyNames = map[int]string{
	CurrencyBlue:   "blue",
	CurrencyGreen:  "green",
	CurrencyRed:    "red",
	CurrencyYellow: "yellow",
}

var CurrencyAbbr = map[int]string{
	CurrencyBlue:   "B",
	CurrencyGreen:  "G",
	CurrencyRed:    "R",
	CurrencyYellow: "Y",
}

var CurrencyColours = map[int]string{
	CurrencyBlue:   render.Blue,
	CurrencyGreen:  render.Green,
	CurrencyRed:    render.Red,
	CurrencyYellow: render.Yellow,
}

type Card struct {
	Currency, Value int
}

func (c Card) Compare(other card.Comparer) (int, bool) {
	oc := other.(Card)
	curDiff := c.Currency - oc.Currency
	if curDiff != 0 {
		return curDiff, true
	}
	return c.Value - oc.Value, true
}

type ScoringCard struct{}

func (c ScoringCard) Compare(other card.Comparer) (int, bool) {
	return 0, false
}

func (c Card) String() string {
	return render.Markup(fmt.Sprintf(
		"%s%d",
		CurrencyAbbr[c.Currency],
		c.Value,
	), CurrencyColours[c.Currency], true)
}

func Deck(players int) card.Deck {
	deck := card.Deck{}
	n := 3
	if players == 2 {
		n = 2
	}
	for _, c := range Currencies {
		for v := 1; v <= 9; v++ {
			for i := 0; i < n; i++ {
				deck = deck.Push(Card{c, v})
			}
		}
	}
	return deck
}
