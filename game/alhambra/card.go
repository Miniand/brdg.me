package alhambra

import (
	"math/rand"
	"time"
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

type Card struct {
	Currency, Value int
}

func Deck() []Card {
	deck := []Card{}
	for _, c := range Currencies {
		for v := 1; v <= 9; v++ {
			for i := 0; i < 3; i++ {
				deck = append(deck, Card{c, v})
			}
		}
	}
	return deck
}

func ShuffleCards(cards []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(cards)
	shuffled := make([]Card, l)
	for i, k := range r.Perm(l) {
		shuffled[i] = cards[k]
	}
	return shuffled
}
