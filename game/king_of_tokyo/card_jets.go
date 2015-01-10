package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardJets struct{}

func (c CardJets) Name() string {
	return "Jets"
}

func (c CardJets) Description() string {
	return "{{b}}You suffer no damage{{_b}} when yielding Tokyo."
}

func (c CardJets) Cost() int {
	return 5
}

func (c CardJets) Kind() int {
	return CardKindKeep
}

func (c CardJets) ModifyDamage(
	game *Game,
	player, attacker, damage, defenderAction int,
) int {
	if defenderAction == DefenderActionLeaveTokyo {
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s avoided damage by leaving Tokyo ({{b}}%s{{_b}})",
			game.RenderName(player),
			c.Name(),
		)))
		return 0
	}
	return damage
}
