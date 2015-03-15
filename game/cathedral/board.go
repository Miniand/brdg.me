package cathedral

import (
	"bytes"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

type Board [10][10]Tile

func (b Board) TileAt(x, y int) (Tile, bool) {
	if x < 0 || y < 0 || x >= 10 || y >= 10 {
		return Tile{}, false
	}
	return b[y][x], true
}

func (b Board) Render() string {
	buf := bytes.NewBuffer([]byte{})
	cells := [][]interface{}{}
	// Header
	header := []interface{}{}
	header = append(header, render.Bold(WallStrs[DirDown|DirRight]))
	for i := 0; i < len(b[0]); i++ {
		header = append(header, render.Bold(strings.Repeat(
			WallStrs[DirLeft|DirRight],
			TileWidth,
		)))
	}
	header = append(header, render.Bold(WallStrs[DirDown|DirLeft]))
	cells = append(cells, header)
	// Body
	for y, xs := range b {
		row := []interface{}{}
		row = append(row, SideWall)
		for x, _ := range xs {
			rt, ok := RenderTile(b, x, y)
			if !ok {
				rt = RenderEmptyTile(x, y)
			}
			row = append(row, rt)
		}
		row = append(row, SideWall)
		cells = append(cells, row)
	}
	// Footer
	footer := []interface{}{}
	footer = append(footer, render.Bold(WallStrs[DirUp|DirRight]))
	for i := 0; i < len(b[0]); i++ {
		footer = append(footer, render.Bold(strings.Repeat(
			WallStrs[DirLeft|DirRight],
			TileWidth,
		)))
	}
	footer = append(footer, render.Bold(WallStrs[DirUp|DirLeft]))
	cells = append(cells, footer)
	buf.WriteString(render.Table(cells, 0, 0))
	return buf.String()
}
