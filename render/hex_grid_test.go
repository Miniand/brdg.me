package render

import (
	"fmt"
	"regexp"

	"github.com/Miniand/brdg.me/game/grid"
	"github.com/Miniand/brdg.me/game/grid/hex"
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
func prefixPeriods(in string) string {
	return regexp.MustCompile("(?m)^").ReplaceAllString(in, ".")
}

func ExampleHexGrid() {
	g := hex.Grid{}
	g.SetTile(grid.Loc{2, 1}, T{"egg", "blue", 0})
	g.SetTile(grid.Loc{0, 0}, T{"fart", "red", 0})
	fmt.Println(prefixPeriods(RenderHexGrid(g, 2)))
	// Output:
	// .  _____
	// . /     \
	// ./ fart  \
	// .\       /
	// . \_____/        _____
	// .               /     \
	// .              /  egg  \
	// .              \       /
	// .               \_____/
}
