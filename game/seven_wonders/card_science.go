package seven_wonders

import (
	"encoding/gob"

	"github.com/Miniand/brdg.me/game/cost"
)

var AllFields = []int{
	FieldMathematics,
	FieldEngineering,
	FieldTheology,
}

func init() {
	gob.Register(CardScience{})
}

type CardScience struct {
	Card
	Fields []int
}

func NewCardScience(
	name string,
	cost cost.Cost,
	field int,
	freeWith, makesFree []string,
) CardScience {
	return CardScience{
		NewCard(name, CardKindScientific, cost, freeWith, makesFree),
		[]int{field},
	}
}

func (c CardScience) SuppString() string {
	return RenderResourceList(c.Fields, "/")
}

func (c CardScience) ScienceField(player int, g *Game) int {
	if len(c.Fields) == 1 {
		return c.Fields[0]
	}
	panic("implement")
}
