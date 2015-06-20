package king_of_tokyo

type CardNovaBreath struct{}

func (c CardNovaBreath) Name() string {
	return "Nova Breath"
}

func (c CardNovaBreath) Description() string {
	return "{{b}}Your attacks damage all other monsters.{{_b}}"
}

func (c CardNovaBreath) Cost() int {
	return 7
}

func (c CardNovaBreath) Kind() int {
	return CardKindKeep
}

func (c CardNovaBreath) ModifyAttackTargets(game *Game, player int, targets []int) []int {
	targets = []int{}
	for p := range game.Players {
		if p != player {
			targets = append(targets, p)
		}
	}
	return targets
}
