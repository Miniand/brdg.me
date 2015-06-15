package bang_dice

import (
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

const (
	DieArrow = iota
	DieDynamite
	Die1
	Die2
	DieBeer
	DieGatling
)

var DieColours = map[int]string{
	DieArrow:    render.Green,
	DieDynamite: render.Red,
	Die1:        render.Black,
	Die2:        render.Black,
	DieBeer:     render.Yellow,
	DieGatling:  render.Blue,
}

var DieStrings = map[int]string{
	DieArrow:    "A",
	DieDynamite: "D",
	Die1:        "1",
	Die2:        "2",
	DieBeer:     "B",
	DieGatling:  "G",
}

func RenderDie(die int) string {
	return render.Bold(fmt.Sprintf(
		"[%s]",
		render.Colour(DieStrings[die], DieColours[die]),
	))
}
