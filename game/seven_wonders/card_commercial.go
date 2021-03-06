package seven_wonders

import (
	"encoding/gob"
	"strings"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/render"
)

const (
	DirLeft  = -1
	DirDown  = 0
	DirRight = 1
)

var DirStrings = map[int]string{
	DirLeft:  render.Bold("<"),
	DirDown:  render.Bold("v"),
	DirRight: render.Bold(">"),
}

var (
	DirAll        = []int{DirLeft, DirDown, DirRight}
	DirNeighbours = []int{DirLeft, DirRight}
	DirSelf       = []int{DirDown}
)

func init() {
	gob.Register(CardCommercialTrade{})
	gob.Register(CardCommercialTavern{})
}

type TradeDiscounter interface {
	TradeDiscount() (dirs []int, goods []int)
}

type CardCommercialTrade struct {
	Card
	Directions []int
	Goods      []int
}

func NewCardCommercialTrade(
	name string,
	cost cost.Cost,
	directions, goods []int,
	freeWith, makesFree []string,
) CardCommercialTrade {
	if directions == nil || len(directions) == 0 {
		panic("no directions")
	}
	if goods == nil || len(goods) == 0 {
		panic("no goods")
	}
	return CardCommercialTrade{
		NewCard(name, CardKindCommercial, cost, freeWith, makesFree),
		directions,
		goods,
	}
}

func (c CardCommercialTrade) SuppString() string {
	parts := []string{
		RenderMoney(1),
		"for",
		RenderResourceList(c.Goods, "/"),
		"from",
		RenderDirections(c.Directions),
	}
	return strings.Join(parts, " ")
}

func (c CardCommercialTrade) TradeDiscount() (dirs []int, goods []int) {
	return c.Directions, c.Goods
}

type CardCommercialTavern struct {
	Card
}

func NewCardCommercialTavern() CardCommercialTavern {
	return CardCommercialTavern{
		NewCard(CardTavern, CardKindCommercial, nil, nil, nil),
	}
}

func (c CardCommercialTavern) SuppString() string {
	return RenderMoney(5)
}

func (c CardCommercialTavern) HandlePostActionExecute(player int, g *Game) {
	g.Coins[player] += 5
}
