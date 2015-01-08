package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardDedicatedNewsTeam struct{}

func (c CardDedicatedNewsTeam) Name() string {
	return "Dedicated News Team"
}

func (c CardDedicatedNewsTeam) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} whenever you buy a card.",
		RenderVP(1),
	)
}

func (c CardDedicatedNewsTeam) Cost() int {
	return 3
}

func (c CardDedicatedNewsTeam) Kind() int {
	return CardKindKeep
}

func (c CardDedicatedNewsTeam) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	game.Boards[player].ModifyVP(1)
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(1),
		c.Name(),
	)))
}
