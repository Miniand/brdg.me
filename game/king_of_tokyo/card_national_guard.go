package king_of_tokyo

import "fmt"

type CardNationalGuard struct{}

func (c CardNationalGuard) Name() string {
	return "National Guard"
}

func (c CardNationalGuard) Description() string {
	return fmt.Sprintf(
		"{{b}}%s and take 2 damage.{{_b}}",
		RenderVPChange(2),
	)
}

func (c CardNationalGuard) Cost() int {
	return 3
}

func (c CardNationalGuard) Kind() int {
	return CardKindDiscard
}

func (c CardNationalGuard) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	game.Boards[player].ModifyVP(2)
	game.TakeDamage(player, 2)
}
