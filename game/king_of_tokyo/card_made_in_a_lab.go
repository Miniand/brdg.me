package king_of_tokyo

type CardMadeInALab struct{}

func (c CardMadeInALab) Name() string {
	return "Made in a Lab"
}

func (c CardMadeInALab) Description() string {
	return "When purchasing cards you can {{b}}peek at and purchase the top card{{_b}} of the deck."
}

func (c CardMadeInALab) Cost() int {
	return 2
}

func (c CardMadeInALab) Kind() int {
	return CardKindKeep
}
