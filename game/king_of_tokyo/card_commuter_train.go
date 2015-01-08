package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardCommuterTrain struct{}

func (c CardCommuterTrain) Name() string {
	return "Commuter Train"
}

func (c CardCommuterTrain) Description() string {
	return RenderVPChange(2)
}

func (c CardCommuterTrain) Cost() int {
	return 4
}

func (c CardCommuterTrain) Kind() int {
	return CardKindDiscard
}

func (c CardCommuterTrain) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	game.Boards[player].ModifyVP(2)
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(2),
		c.Name(),
	)))
}
