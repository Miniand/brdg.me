package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardCornerStore struct{}

func (c CardCornerStore) Name() string {
	return "Corner Store"
}

func (c CardCornerStore) Description() string {
	return RenderVPChange(1)
}

func (c CardCornerStore) Cost() int {
	return 3
}

func (c CardCornerStore) Kind() int {
	return CardKindDiscard
}

func (c CardCornerStore) PostCardBuy(game *Game, player int, card CardBase, cost int) {
	game.Boards[player].ModifyVP(1)
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(1),
		c.Name(),
	)))
}
