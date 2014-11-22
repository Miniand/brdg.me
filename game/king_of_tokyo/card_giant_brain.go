package king_of_tokyo

type CardGiantBrain struct{}

func (c CardGiantBrain) Name() string {
	return "Giant Brain"
}

func (c CardGiantBrain) Description() string {
	return "{{b}}You have 1 extra reroll{{_b}} each turn"
}

func (c CardGiantBrain) Cost() int {
	return 5
}

func (c CardGiantBrain) Kind() int {
	return CardKindKeep
}

func (c CardGiantBrain) ModifyRollCount(game *Game, player, rollCount int) int {
	return rollCount + 1
}
