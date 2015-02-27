package cathedral

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

const (
	TileWidth  = 6
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
		for x, tile := range xs {
			row = append(row, RenderEmptyTile(x, y))
			if tile.Player == NoPlayer {
			} else {
			}
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

func RenderEmptyTile(x, y int) string {
	buf := bytes.NewBufferString(strings.Repeat(fmt.Sprintf(
		"%s\n",
		strings.Repeat(NoTileStr, TileWidth),
	), emptyAbove))
	s := TileText(x, y)
	remainingWidth := TileWidth - len(s)
	buf.WriteString(strings.Repeat(NoTileStr, (remainingWidth+1)/2))
	buf.WriteString(render.Bold(s))
	buf.WriteString(strings.Repeat(NoTileStr, remainingWidth/2))
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
