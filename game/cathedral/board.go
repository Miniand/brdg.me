package cathedral

import (
	"bytes"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

var (
	AllLocs   []Loc
	LocsByRow [][]Loc
)

func init() {
	AllLocs = make([]Loc, 100)
	LocsByRow = make([][]Loc, 10)
	for y := 0; y < 10; y++ {
		LocsByRow[y] = make([]Loc, 10)
		for x := 0; x < 10; x++ {
			l := Loc{x, y}
			AllLocs[y*10+x] = l
			LocsByRow[y][x] = l
		}
	}
}

type Board map[Loc]Tile

func (b Board) TileAt(loc Loc) (Tile, bool) {
	t, ok := b[loc]
	return t, ok
}

func (b Board) Render() string {
	buf := bytes.NewBuffer([]byte{})
	cells := [][]interface{}{}
	// Header
	header := []interface{}{}
	header = append(header, render.Bold(WallStrs[DirDown|DirRight]))
	for i := 0; i < len(LocsByRow[0]); i++ {
		header = append(header, render.Bold(strings.Repeat(
			WallStrs[DirLeft|DirRight],
			TileWidth,
		)))
	}
	header = append(header, render.Bold(WallStrs[DirDown|DirLeft]))
	cells = append(cells, header)
	// Body
	for _, xs := range LocsByRow {
		row := []interface{}{}
		row = append(row, SideWall)
		for _, l := range xs {
			rt, ok := RenderTile(b, l)
			if !ok {
				rt = RenderEmptyTile(l)
			}
			row = append(row, rt)
		}
		row = append(row, SideWall)
		cells = append(cells, row)
	}
	// Footer
	footer := []interface{}{}
	footer = append(footer, render.Bold(WallStrs[DirUp|DirRight]))
	for i := 0; i < len(LocsByRow[0]); i++ {
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
