package render

import (
	"fmt"
	"math"

	"github.com/Miniand/brdg.me/game/grid"
	"github.com/Miniand/brdg.me/game/grid/hex"
)

const (
	SIZE_SMALL     = 2
	EDGE_DRAW_NONE = iota
	EDGE_DRAW_PARTIAL
	EDGE_DRAW_FULL
)

func HorizontalSideLength(size int) int {
	return int(math.Floor(float64(size) * 5 / 2))
}

func Height(size int) int {
	return size * 2
}

func Width(size int) int {
	return OffsetWidth(size) + size
}

func OffsetWidth(size int) int {
	return HorizontalSideLength(size) + size
}

func PixelToLocs(x, y, size int) (locs []grid.Loc) {
	locs = append(locs, grid.Loc{0, 0})
	return
}

func RenderHexGrid(g hex.Grid, size int) string {
	if size < SIZE_SMALL {
		size = SIZE_SMALL
	}
	c := NewCanvas()
	horizontalSideLength := HorizontalSideLength(size)
	offsetHeight := Height(size)
	offsetWidth := OffsetWidth(size)
	fullWidth := Width(size)
	g.Each(func(l grid.Loc, t interface{}) {
		if t == nil {
			return
		}
		// Initialise vars
		drawX := l.X * offsetWidth
		drawY := l.Y*offsetHeight + size*(l.X&1)
		colour := "black"
		if tc, ok := t.(grid.Colourer); ok {
			colour = tc.Colour()
		}
		priority := 0
		if tcp, ok := t.(grid.ColourPrioritiser); ok {
			priority = tcp.ColourPriority()
		}
		// Tile message
		message := ""
		if tm, ok := t.(grid.Messager); ok {
			message = tm.Message()
		}
		if message != "" {
			c.Draw(drawX+(fullWidth-len(message))/2, drawY+size-1, fmt.Sprintf(
				`{{c "%s"}}%s{{_c}}`, colour, message))
		}
		// Loc text
		locString := fmt.Sprintf("%d %d", l.X, l.Y)
		c.Draw(drawX+(fullWidth-len(locString))/2, drawY+size, fmt.Sprintf(
			`{{c "gray"}}%s{{_c}}`, locString))
		// Horizontal lines
		drawModes := map[int]int{}
		for _, dir := range hex.Dirs() {
			drawModes[dir] = EDGE_DRAW_FULL
			neighbour := g.Tile(g.Neighbour(l, dir))
			if neighbour == nil {
				continue
			}
			nPriority := 0
			if cp, ok := neighbour.(grid.ColourPrioritiser); ok {
				nPriority = cp.ColourPriority()
			}
			switch {
			case priority < nPriority:
				drawModes[dir] = EDGE_DRAW_NONE
			case priority == nPriority:
				drawModes[dir] = EDGE_DRAW_PARTIAL
			}
		}
		for i := 0; i < horizontalSideLength; i++ {
			// Top
			if drawModes[hex.DIR_NORTH] == EDGE_DRAW_FULL ||
				(drawModes[hex.DIR_NORTH] == EDGE_DRAW_PARTIAL &&
					i%2 == 0) {
				c.Draw(drawX+size+i, drawY-1,
					fmt.Sprintf(`{{c "%s"}}_{{_c}}`, colour))
			}
			// Bottom
			if drawModes[hex.DIR_SOUTH] == EDGE_DRAW_FULL ||
				(drawModes[hex.DIR_SOUTH] == EDGE_DRAW_PARTIAL &&
					(i+1)%2 == 0) {
				c.Draw(drawX+size+i, drawY+size*2-1,
					fmt.Sprintf(`{{c "%s"}}_{{_c}}`, colour))
			}
		}
		// Diagonal lines
		for i := 0; i < size; i++ {
			// Top left
			if drawModes[hex.DIR_NORTH_WEST] == EDGE_DRAW_FULL ||
				(drawModes[hex.DIR_NORTH_WEST] == EDGE_DRAW_PARTIAL &&
					i%2 == 0) {
				c.Draw(drawX+i, drawY+size-1-i,
					fmt.Sprintf(`{{c "%s"}}/{{_c}}`, colour))
			}
			// Top right
			if drawModes[hex.DIR_NORTH_EAST] == EDGE_DRAW_FULL ||
				(drawModes[hex.DIR_NORTH_EAST] == EDGE_DRAW_PARTIAL &&
					i%2 == 0) {
				c.Draw(drawX+offsetWidth+i, drawY+i,
					fmt.Sprintf(`{{c "%s"}}\{{_c}}`, colour))
			}
			// Bottom right
			if drawModes[hex.DIR_SOUTH_EAST] == EDGE_DRAW_FULL ||
				(drawModes[hex.DIR_SOUTH_EAST] == EDGE_DRAW_PARTIAL &&
					i%2 == 0) {
				c.Draw(drawX+offsetWidth+size-1-i, drawY+size+i,
					fmt.Sprintf(`{{c "%s"}}/{{_c}}`, colour))
			}
			// Bottom left
			if drawModes[hex.DIR_SOUTH_WEST] == EDGE_DRAW_FULL ||
				(drawModes[hex.DIR_SOUTH_WEST] == EDGE_DRAW_PARTIAL &&
					i%2 == 0) {
				c.Draw(drawX+size-1-i, drawY+size*2-1-i,
					fmt.Sprintf(`{{c "%s"}}\{{_c}}`, colour))
			}
		}
	})
	return c.Render()
}
