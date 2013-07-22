package card

type Comparer interface {
	// < 0 means the context card is smaller, the argument card is larger
	// 0 means the card is equal
	// > 1 means the context card is larger, the argument card is smaller
	// The second value denotes whether the cards are comparable
	Compare(Comparer) (int, bool)
}

type Card interface {
	Comparer
}
