package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardFireBreathing struct{}

func (c CardFireBreathing) Name() string {
	return "Fire Breathing"
}

func (c CardFireBreathing) Description() string {
	return "{{b}}Your neighbours take 1 extra damage{{_b}} when you deal damage."
}

func (c CardFireBreathing) Cost() int {
	return 4
}

func (c CardFireBreathing) Kind() int {
	return CardKindKeep
}

func (c CardFireBreathing) ModifyAttackDamageForPlayer(
	game *Game,
	player, attacked, damage int,
) int {
	l := len(game.Players)
	leftNeighbour := player - 1
	if leftNeighbour < 0 {
		leftNeighbour += l
	}
	rightNeighbour := (player + 1) % l
	switch attacked {
	case leftNeighbour, rightNeighbour:
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s deals 1 extra damage to their neighbour %s ({{b}}%s{{_b}})",
			game.RenderName(player),
			game.RenderName(attacked),
			c.Name(),
		)))
		damage += 1
	}
	return damage
}
