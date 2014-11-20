package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardCamouflage struct{}

func (c CardCamouflage) Name() string {
	return "Camouflage"
}

func (c CardCamouflage) Description() string {
	return fmt.Sprintf(
		"If you take damage roll a die for each damage point. {{b}}On a %s you do not take that damage point.{{_b}}",
		RenderDie(DieHeal),
	)
}

func (c CardCamouflage) Cost() int {
	return 3
}

func (c CardCamouflage) Kind() int {
	return CardKindKeep
}

func (c CardCamouflage) ModifyDamage(game *Game, player, attacker, damage int) int {
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s is rolling to avoid damage ({{b}}%s{{_b}})",
		game.RenderName(player),
		c.Name(),
	)))
	roll := RollDice(damage)
	game.LogRoll(player, roll, []int{})
	avoided := 0
	for _, d := range roll {
		if d == DieHeal {
			avoided += 1
		}
	}
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s avoided {{b}}%d{{_b}} damage",
		game.RenderName(player),
		avoided,
	)))
	damage -= avoided
	if damage < 0 {
		damage = 0
	}
	return damage
}
