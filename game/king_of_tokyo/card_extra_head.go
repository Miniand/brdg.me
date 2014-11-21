package king_of_tokyo

type CardExtraHead struct{}

func (c CardExtraHead) Name() string {
	return "Extra Head"
}

func (c CardExtraHead) Description() string {
	return "{{b}}You get 1 extra die.{{_b}}"
}

func (c CardExtraHead) Cost() int {
	return 7
}

func (c CardExtraHead) Kind() int {
	return CardKindKeep
}

func (c CardExtraHead) ModifyDiceCount(game *Game, player, diceCount int) int {
	return diceCount + 1
}
