package king_of_tokyo

type CardHighAltitudeBombing struct{}

func (c CardHighAltitudeBombing) Name() string {
	return "High Altitude Bombing"
}

func (c CardHighAltitudeBombing) Description() string {
	return "{{b}}All monsters{{_b}} (including you) {{b}}take 3 damage.{{_b}}"
}

func (c CardHighAltitudeBombing) Cost() int {
	return 4
}

func (c CardHighAltitudeBombing) Kind() int {
	return CardKindDiscard
}

func (c CardHighAltitudeBombing) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	for p, _ := range game.Players {
		if game.Boards[p].Health > 0 {
			game.TakeDamage(p, 3)
		}
	}
}
