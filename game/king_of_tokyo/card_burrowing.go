package king_of_tokyo

type CardBurrowing struct{}

func (c CardBurrowing) Name() string {
	return "Burrowing"
}

func (c CardBurrowing) Description() string {
	return "{{b}}Deal 1 extra damage when in Tokyo. Deal 1 damage when yielding Tokyo{{_b}} to the monster taking it."
}

func (c CardBurrowing) Cost() int {
	return 5
}

func (c CardBurrowing) Kind() int {
	return CardKindKeep
}

func (c CardBurrowing) ModifyAttack(game *Game, attack int) int {
	if game.PlayerLocation(game.CurrentPlayer) != LocationOutside {
		attack += 1
	}
	return attack
}
