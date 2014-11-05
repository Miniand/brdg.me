package king_of_tokyo

type CardSpikedTail struct{}

func (c CardSpikedTail) Name() string {
	return "Spiked Tail"
}

func (c CardSpikedTail) Description() string {
	return "When you attack {{b}}deal 1 extra damage{{_b}}."
}

func (c CardSpikedTail) Cost() int {
	return 5
}

func (c CardSpikedTail) Kind() int {
	return CardKindKeep
}
