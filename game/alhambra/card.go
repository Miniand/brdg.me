package alhambra

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
