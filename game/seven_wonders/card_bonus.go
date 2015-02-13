package seven_wonders

import (
	"encoding/gob"
	"strings"

	"github.com/Miniand/brdg.me/game/cost"
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

func (c CardBonus) HandlePostBuild(player int, g *Game) {
	if c.Coins == 0 {
		return
	}
}
