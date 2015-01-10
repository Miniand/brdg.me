package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardGasRefinery struct{}

func (c CardGasRefinery) Name() string {
	return "Gas Refinery"
}

func (c CardGasRefinery) Description() string {
	return fmt.Sprintf(
		"{{b}}%s and deal 3 damage to all other monsters.{{_b}}",
		RenderVPChange(2),
	)
}

func (c CardGasRefinery) Cost() int {
	return 6
}

func (c CardGasRefinery) Kind() int {
	return CardKindDiscard
}

func (c CardGasRefinery) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	game.Boards[player].ModifyVP(2)
	for p, _ := range game.Players {
		if p != player {
			game.DealDamage(player, p, 3, DefenderActionNone)
		}
	}
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s and dealt 3 damage to other monsters ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(1),
		c.Name(),
	)))
}
