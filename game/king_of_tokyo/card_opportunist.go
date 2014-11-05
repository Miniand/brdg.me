package king_of_tokyo

type CardOpportunist struct{}

func (c CardOpportunist) Name() string {
	return "Opportunist"
}

func (c CardOpportunist) Description() string {
	return "{{b}}Whenever a new card is revealed you have the option of purchasing it{{_b}} as soon as it is revealed."
}

func (c CardOpportunist) Cost() int {
	return 3
}

func (c CardOpportunist) Kind() int {
	return CardKindKeep
}
