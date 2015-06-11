package age_of_war

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	Dice1Infantry = iota
	Dice2Infantry
	Dice3Infantry
	DiceArchery
	DiceCavalry
	DiceDaimyo
)

var DiceInfantry = map[int]int{
	Dice1Infantry: 1,
	Dice2Infantry: 2,
	Dice3Infantry: 3,
}

var DiceStrings = map[int]string{
	Dice1Infantry: "1 inf",
	Dice2Infantry: "2 inf",
	Dice3Infantry: "3 inf",
	DiceArchery:   "arch",
	DiceCavalry:   "cav",
	DiceDaimyo:    "dai",
}

var InfantryColour = render.Blue

var DiceColours = map[int]string{
	Dice1Infantry: InfantryColour,
	Dice2Infantry: InfantryColour,
	Dice3Infantry: InfantryColour,
	DiceArchery:   render.Magenta,
	DiceCavalry:   render.Green,
	DiceDaimyo:    render.Red,
}

func Roll() int {
	return rnd.Int() % 6
}

func RollN(n int) []int {
	if n <= 0 {
		return []int{}
	}
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = Roll()
	}
	return ints
}

func (g *Game) Roll(n int) {
	g.CurrentRoll = RollN(n)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s rolled  %s",
		g.PlayerName(g.CurrentPlayer),
		strings.Join(RenderDice(g.CurrentRoll), "  "),
	)))
}
