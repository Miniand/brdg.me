package king_of_tokyo

import "fmt"

type CardNuclearPowerPlant struct{}

func (c CardNuclearPowerPlant) Name() string {
	return "Nuclear Power Plant"
}

func (c CardNuclearPowerPlant) Description() string {
	return fmt.Sprintf(
		"{{b}}%s and heal 3 damage.{{_b}}",
		RenderVPChange(2),
	)
}

func (c CardNuclearPowerPlant) Cost() int {
	return 6
}

func (c CardNuclearPowerPlant) Kind() int {
	return CardKindDiscard
}

func (c CardNuclearPowerPlant) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	game.Boards[player].ModifyVP(2)
	game.Boards[player].ModifyHealth(3)
}
