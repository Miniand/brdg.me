package cathedral

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

const (
	TileWidth  = 5
	TileHeight = 3
)

var NoTileStr = `▒`

var WallStrs = map[int]string{
	DirUp | DirDown | DirLeft | DirRight: render.Bold("╬"),
	DirUp | DirDown | DirLeft:            render.Bold("╣"),
	DirUp | DirDown | DirRight:           render.Bold("╠"),
	DirUp | DirLeft | DirRight:           render.Bold("╩"),
	DirDown | DirLeft | DirRight:         render.Bold("╦"),
	DirUp | DirLeft:                      render.Bold("╝"),
	DirUp | DirRight:                     render.Bold("╚"),
	DirDown | DirLeft:                    render.Bold("╗"),
	DirDown | DirRight:                   render.Bold("╔"),
	DirLeft | DirRight:                   render.Bold("═"),
	DirLeft:                              render.Bold("═"),
	DirRight:                             render.Bold("═"),
	DirUp | DirDown:                      render.Bold("║"),
	DirUp:                                render.Bold("║"),
	DirDown:                              render.Bold("║"),
}

var SideWall = render.Bold(strings.TrimSpace(strings.Repeat(fmt.Sprintf(
	"%s\n",
	WallStrs[DirUp|DirDown],
), TileHeight)))

func (g *Game) RenderForPlayer(string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	cells := [][]interface{}{}
	// Header
	header := []interface{}{}
	header = append(header, render.Bold(WallStrs[DirDown|DirRight]))
	for i := 0; i < len(g.Board[0]); i++ {
		header = append(header, render.Bold(strings.Repeat(
			WallStrs[DirLeft|DirRight],
			TileWidth,
		)))
	}
	header = append(header, render.Bold(WallStrs[DirDown|DirLeft]))
	cells = append(cells, header)
	// Body
	for y, xs := range g.Board {
		row := []interface{}{}
		row = append(row, SideWall)
		for x, _ := range xs {
			row = append(row, g.RenderTile(x, y))
		}
		row = append(row, SideWall)
		cells = append(cells, row)
	}
	// Footer
	footer := []interface{}{}
	footer = append(footer, render.Bold(WallStrs[DirUp|DirRight]))
	for i := 0; i < len(g.Board[0]); i++ {
		footer = append(footer, render.Bold(strings.Repeat(
			WallStrs[DirLeft|DirRight],
			TileWidth,
		)))
	}
	footer = append(footer, render.Bold(WallStrs[DirUp|DirLeft]))
	cells = append(cells, footer)
	buf.WriteString(render.Table(cells, 0, 0))
	return buf.String(), nil
}

var (
	emptyAbove = (TileHeight - 1) / 2
	emptyBelow = TileHeight / 2
)

func (g *Game) RenderTile(x, y int) string {
	t, ok := g.TileAt(x, y)
	if !ok || t.Player == NoPlayer {
		return RenderEmptyTile(x, y)
	}
	return RenderPlayerTile(t, g.OpenSides(x, y))
}

func RenderPlayerTile(tile Tile, open map[int]bool) string {
	// Top row
	buf := bytes.NewBufferString(RenderCorner(DirUp|DirLeft, open))
	c := WallStrs[DirLeft|DirRight]
	if open[DirUp] {
		c = " "
	}
	buf.WriteString(strings.Repeat(c, TileWidth-2))
	buf.WriteString(RenderCorner(DirUp|DirRight, open))
	buf.WriteString("\n")

	// Middle rows
	left := WallStrs[DirUp|DirDown]
	if open[DirLeft] {
		left = " "
	}
	right := WallStrs[DirUp|DirDown]
	if open[DirRight] {
		right = " "
	}
	middleRow := fmt.Sprintf(
		"%s%s%s\n",
		left,
		strings.Repeat(" ", TileWidth-2),
		right,
	)
	buf.WriteString(strings.Repeat(middleRow, TileHeight-2))

	// Bottom row
	buf.WriteString(RenderCorner(DirDown|DirLeft, open))
	c = WallStrs[DirLeft|DirRight]
	if open[DirDown] {
		c = " "
	}
	buf.WriteString(strings.Repeat(c, TileWidth-2))
	buf.WriteString(RenderCorner(DirDown|DirRight, open))

	return render.Markup(buf.String(), render.PlayerColour(tile.Player), true)
}

func RenderCorner(dir int, open map[int]bool) string {
	// If all three tiles in dir are open, then render nothing.
	numOpen := 0
	for _, d := range Dirs {
		if dir&d == d && open[d] {
			numOpen++
			if numOpen == 3 {
				return " "
			}
		}
	}

	// Map of one corner direction referencing the other.
	cornerMap := map[int]int{}
	first := -1
	for _, d := range Dirs {
		if dir&d != d {
			continue
		}
		if first == -1 {
			first = d
		} else {
			cornerMap[first] = d
			cornerMap[d] = first
			break
		}
	}

	var corner int
	for d, other := range cornerMap {
		if open[d] {
			corner = corner | d
		} else {
			corner = corner | DirInv(other)
		}
	}
	return WallStrs[corner]
}

func RenderEmptyTile(x, y int) string {
	buf := bytes.NewBufferString(strings.Repeat(fmt.Sprintf(
		"%s\n",
		strings.Repeat(NoTileStr, TileWidth),
	), emptyAbove))
	s := TileText(x, y)
	remainingWidth := TileWidth - len(s)
	buf.WriteString(strings.Repeat(NoTileStr, remainingWidth/2))
	buf.WriteString(render.Bold(s))
	buf.WriteString(strings.Repeat(NoTileStr, (remainingWidth+1)/2))
	buf.WriteByte('\n')
	buf.WriteString(strings.TrimSpace(strings.Repeat(fmt.Sprintf(
		"%s\n",
		strings.Repeat(NoTileStr, TileWidth),
	), emptyBelow)))
	return render.Colour(buf.String(), render.Gray)
}

func TileText(x, y int) string {
	return fmt.Sprintf("%c%d", 'A'+y, x+1)
}
