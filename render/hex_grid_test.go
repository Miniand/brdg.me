package render

import (
	"fmt"
	"github.com/Miniand/brdg.me/game/grid"
	"github.com/Miniand/brdg.me/game/grid/hex"
	"testing"
)

type T struct {
	message, colour string
	colourPriority  int
}

func (t T) Message() string {
	return t.message
}
func (t T) Colour() string {
	return t.colour
}
func (t T) ColourPriority() int {
	return t.colourPriority
}

func TestHexGrid(t *testing.T) {
	g := hex.Grid{}
	// g.SetTile(grid.Loc{-3, -4}, T{"fart", "blue", 0})
	g.SetTile(grid.Loc{0, 0}, T{"fart", "red", 0})
	fmt.Println(RenderHexGrid(g, 2))
}
