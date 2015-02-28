package seven_wonders

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/game/log"
)

func init() {
	gob.Register(CardBonus{})
}

type CardBonus struct {
	Card
	TargetKinds []int
	Directions  []int
	VP          int
	Coins       int
}

func NewCardBonus(
	name string,
	kind int,
	cost cost.Cost,
	targetKinds, directions []int,
	vp, coins int,
	freeWith, makesFree []string,
) CardBonus {
	if targetKinds == nil || len(targetKinds) == 0 {
		panic("no targetKinds")
	}
	if directions == nil || len(directions) == 0 {
		panic("no directions")
	}
	return CardBonus{
		NewCard(name, kind, cost, freeWith, makesFree),
		targetKinds,
		directions,
		vp,
		coins,
	}
}

func (c CardBonus) SuppString() string {
	reward := []string{}
	if c.VP > 0 {
		reward = append(reward, RenderVP(c.VP))
	}
	if c.Coins > 0 {
		reward = append(reward, RenderMoney(c.Coins))
	}
	parts := []string{
		strings.Join(reward, " and "),
		"for each",
		RenderResourceList(c.TargetKinds, " "),
		"owned by",
		RenderDirections(c.Directions),
	}
	return strings.Join(parts, " ")
}

func (c CardBonus) SumTargetKinds(player int, g *Game) int {
	sum := 0
	for _, dir := range c.Directions {
		for _, kind := range c.TargetKinds {
			sum += g.PlayerResourceCount(g.NumFromPlayer(player, dir), kind)
		}
	}
	return sum
}

func (c CardBonus) HandlePostActionExecute(player int, g *Game) {
	if c.Coins == 0 {
		return
	}
	sum := c.SumTargetKinds(player, g)
	amt := sum * c.Coins
	g.Coins[player] += amt
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s for each %s owned by %s",
		g.PlayerName(player),
		RenderMoney(amt),
		RenderResourceList(c.TargetKinds, " "),
		RenderDirections(c.Directions),
	)))
}

func (c CardBonus) VictoryPoints(player int, g *Game) int {
	if c.VP == 0 {
		return 0
	}
	return c.SumTargetKinds(player, g) * c.VP
}
