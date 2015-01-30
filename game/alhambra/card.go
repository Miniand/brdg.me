package alhambra

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/render"
)

const (
	CurrencyBlue = iota
	CurrencyGreen
	CurrencyOrange
	CurrencyYellow
)

var Currencies = []int{
	CurrencyBlue,
	CurrencyGreen,
	CurrencyOrange,
	CurrencyYellow,
}

var CurrencyNames = map[int]string{
	CurrencyBlue:   "blue",
	CurrencyGreen:  "green",
	CurrencyOrange: "orange",
	CurrencyYellow: "yellow",
}

var CurrencyAbbr = map[int]string{
	CurrencyBlue:   "B",
	CurrencyGreen:  "G",
	CurrencyOrange: "O",
	CurrencyYellow: "Y",
}

var CurrencyColours = map[int]string{
	CurrencyBlue:   render.Blue,
	CurrencyGreen:  render.Green,
	CurrencyOrange: render.Red,
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

func (c Card) String() string {
	return render.Markup(fmt.Sprintf(
		"%s%d",
		CurrencyAbbr[c.Currency],
		c.Value,
	), CurrencyColours[c.Currency], true)
}

func Deck() card.Deck {
	deck := card.Deck{}
	for _, c := range Currencies {
		for v := 1; v <= 9; v++ {
			for i := 0; i < 3; i++ {
				deck = deck.Push(Card{c, v})
			}
		}
	}
	return deck
}
