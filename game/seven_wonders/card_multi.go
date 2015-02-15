package seven_wonders

import (
	"encoding/gob"
	"fmt"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/game/log"
)

func init() {
	gob.Register(CardMulti{})
}

type CardMulti struct {
	Card
	Resources cost.Cost
}

func (c CardMulti) SuppString() string {
	return RenderResources(c.Resources, " ")
}

func (c CardMulti) AttackStrength() int {
	return c.Resources[AttackStrength]
}

func (c CardMulti) VictoryPoints() int {
	return c.Resources[VP]
}

func (c CardMulti) HandlePostActionExecute(player int, g *Game) {
	if c.Resources[GoodCoin] == 0 {
		return
	}
	g.Coins[player] += c.Resources[GoodCoin]
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s",
		g.PlayerName(player),
		RenderMoney(c.Resources[GoodCoin]),
	)))
}
