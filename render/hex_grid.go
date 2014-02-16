package render

import (
	"bytes"
	"fmt"
	"github.com/Miniand/brdg.me/game/grid"
	"github.com/Miniand/brdg.me/game/grid/hex"
	"math"
	"strings"
)

const (
	SIZE_SMALL = 2
)

func HorizontalSideLength(size int) int {
	return int(math.Floor(float64(size) * 5 / 2))
}

func Width(size int) int {
	return HorizontalSideLength(size) + size*2
}

func PixelToLocs(x, y, size int) (locs []grid.Loc) {
	locs = append(locs, grid.Loc{0, 0})
	return
}

func RenderHexGrid(g hex.Grid, size int) string {
	if size < SIZE_SMALL {
		size = SIZE_SMALL
	}
	lower, upper := g.Bounds()
	if lower.X == 0 && lower.Y == 0 && upper.X == 0 && upper.Y == 0 &&
		g.Tile(lower) == nil {
		return ""
	}
	// FROM
	// If 0, -1
	// If -1, -5
	// If 1, 3
	// UNTIL
	// If 0, 3
	// If -1, -1
	// If 1, 7
	lines := []string{}
	horiz := HorizontalSideLength(size)
	for y := lower.Y*size*2 - 1; y <= (upper.Y+1)*size*2+size-1; y++ {
		line := bytes.Buffer{}
		for x := lower.X * (size + horiz); x < (upper.X+1)*(size+horiz)+size; x++ {
			line.WriteString(fmt.Sprintf("%#v ", PixelToLocs(x, y, size)))
		}
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n")
}
