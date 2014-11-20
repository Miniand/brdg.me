package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardAcidAttack struct{}

func (c CardAcidAttack) Name() string {
	return "Acid Attack"
}

func (c CardAcidAttack) Description() string {
	return "{{b}}Deal 1 extra damage each turn{{_b}} (even when you don't otherwise attack)."
}

func (c CardAcidAttack) Cost() int {
	return 6
}

func (c CardAcidAttack) Kind() int {
	return CardKindKeep
}

func (c CardAcidAttack) ModifyAttack(
	game *Game,
	player, damage int,
	attacked []int,
) (int, []int) {
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s deals 1 extra damage ({{b}}%s{{_b}})",
		game.RenderName(player),
		c.Name(),
	)))
	return damage + 1, attacked
}
