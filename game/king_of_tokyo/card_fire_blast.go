package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardFireBlast struct{}

func (c CardFireBlast) Name() string {
	return "Fire Blast"
}

func (c CardFireBlast) Description() string {
	return "{{b}}Deal 2 damage to all other monsters.{{_b}}"
}

func (c CardFireBlast) Cost() int {
	return 3
}

func (c CardFireBlast) Kind() int {
	return CardKindDiscard
}

func (c CardFireBlast) PostCardBuy(game *Game, player int, card CardBase, cost int) {
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s dealt 2 damage to all other monsters ({{b}}%s{{_b}})",
		game.RenderName(player),
		c.Name(),
	)))
	for p, _ := range game.Players {
		if p == player {
			continue
		}
		game.TakeDamage(p, 2)
	}
}
