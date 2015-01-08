package king_of_tokyo

type CardFrenzy struct{}

func (c CardFrenzy) Name() string {
	return "Frenzy"
}

func (c CardFrenzy) Description() string {
	return "When you purchase this card {{b}}immediately take another turn{{_b}}."
}

func (c CardFrenzy) Cost() int {
	return 7
}

func (c CardFrenzy) Kind() int {
	return CardKindDiscard
}

func (c CardFrenzy) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	// As per FAQ, restart turn immediately without the end of turn phase.
	game.RollPhaceNDice(len(game.CurrentRoll))
}
