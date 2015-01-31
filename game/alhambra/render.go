package alhambra

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

var NoTileStr = `{{c "gray"}}░{{_c}}`

var WallStrs = map[int]string{
	DirUp | DirDown | DirLeft | DirRight: "╬",
	DirUp | DirDown | DirLeft:            "╣",
	DirUp | DirDown | DirRight:           "╠",
	DirUp | DirLeft | DirRight:           "╩",
	DirDown | DirLeft | DirRight:         "╦",
	DirUp | DirLeft:                      "╝",
	DirUp | DirRight:                     "╚",
	DirDown | DirLeft:                    "╗",
	DirDown | DirRight:                   "╔",
	DirLeft | DirRight:                   "═",
	DirLeft:                              "═",
	DirRight:                             "═",
	DirUp | DirDown:                      "║",
	DirUp:                                "║",
	DirDown:                              "║",
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	output := bytes.NewBuffer([]byte{})
	output.WriteString(AddCoordsToGrid(g.Boards[pNum].Grid.Render(1)))
	output.WriteString(render.Bold("\n\nTiles for purchase\n\n"))
	output.WriteString(g.RenderTiles())
	return output.String(), nil
}

func (g *Game) RenderTiles() string {
	gr := Grid{}
	for i, t := range g.Tiles {
		gr[Vect{i * 2, 0}] = t
	}
	output := bytes.NewBufferString(EmptyBorder(gr.Render(0)))
	output.WriteString("\n ")
	output.WriteString(HeaderRow(len(g.Tiles)*2, func(i int) string {
		if i%2 == 1 {
			return ""
		}
		ti := i / 2
		t := g.Tiles[ti]
		if t.Type == TileTypeEmpty {
			return ""
		}
		c := Card{ti, t.Cost}
		return c.String()
	}))
	return output.String()
}

func RenderTileAbbr(tileType int) string {
	return render.Markup(TileAbbrs[tileType], TileColours[tileType], true)
}

func HeaderRow(n int, gen func(i int) string) string {
	output := bytes.NewBuffer([]byte{})
	for i := 0; i < n; i++ {
		output.WriteString(render.Centre(gen(i), TileWidth))
	}
	return output.String()
}

func HeaderRowAlpha(n int) string {
	return HeaderRow(n, func(i int) string {
		return fmt.Sprintf("%c", 'A'+i)
	})
}

func HeaderRowNum(n int) string {
	return HeaderRow(n, func(i int) string {
		return strconv.Itoa(i + 1)
	})
}

func EmptyBorder(grid string) string {
	lines := strings.Split(grid, "\n")
	width := 0
	for _, l := range lines {
		if ll := render.StrLen(l); ll > width {
			width = ll
		}
	}
	lineFormat := fmt.Sprintf("%%-%ds", width)
	top := strings.Repeat(NoTileStr, width+2)
	output := bytes.NewBufferString(top)
	for _, l := range lines {
		output.WriteRune('\n')
		output.WriteString(NoTileStr)
		output.WriteString(fmt.Sprintf(lineFormat, l))
		output.WriteString(NoTileStr)
	}
	output.WriteRune('\n')
	output.WriteString(top)
	return output.String()
}

func AddCoordsToGrid(grid string) string {
	if grid == "" {
		return grid
	}
	lines := strings.Split(grid, "\n")
	width := render.StrLen(lines[0])
	// Left and right
	for i, l := range lines {
		switch i % 2 {
		case 0:
			lines[i] = fmt.Sprintf("    %s", l)
		case 1:
			n := (i + 1) / 2
			lines[i] = fmt.Sprintf("%3d %s %-3d", n, l, n)
		}
	}
	// Top and bottom
	header := "    " + HeaderRowAlpha(width/TileWidth)
	lines = append([]string{header}, append(lines, header)...)
	return strings.Join(lines, "\n")
}

func (g Grid) Render(border int) string {
	var ok bool
	min, max := g.Bounds()
	output := bytes.NewBuffer([]byte{})
	firstRow := true
	for y := min.Y - border; y <= max.Y+border+1; y++ {
		if !firstRow {
			output.WriteRune('\n')
		}
		firstRow = false
		l1 := bytes.NewBuffer([]byte{})
		l2 := bytes.NewBuffer([]byte{})
		for x := min.X - border; x <= max.X+border+1; x++ {
			v := Vect{x, y}
			ct := g.TileAt(v)
			ut := g.TileAt(v.Add(VectUp))
			lt := g.TileAt(v.Add(VectLeft))
			dt := g.TileAt(v.Add(VectUpLeft))
			// Upper left pixel
			r := NoTileStr
			if ct.Type != TileTypeEmpty || ut.Type != TileTypeEmpty ||
				lt.Type != TileTypeEmpty || dt.Type != TileTypeEmpty {
				wallType := 0
				if dt.Walls[DirRight] || ut.Walls[DirLeft] {
					wallType |= DirUp
				}
				if lt.Walls[DirRight] || ct.Walls[DirLeft] {
					wallType |= DirDown
				}
				if dt.Walls[DirDown] || lt.Walls[DirUp] {
					wallType |= DirLeft
				}
				if ut.Walls[DirDown] || ct.Walls[DirUp] {
					wallType |= DirRight
				}
				r, ok = WallStrs[wallType]
				if !ok {
					r = " "
				}
			}
			l1.WriteString(r)
			if x <= max.X+border {
				// Upper pixel
				r = NoTileStr
				if ct.Type != TileTypeEmpty || ut.Type != TileTypeEmpty {
					wallType := 0
					if ct.Walls[DirUp] || ut.Walls[DirDown] {
						wallType = DirLeft | DirRight
					}
					r, ok = WallStrs[wallType]
					if !ok {
						r = " "
					}
				}
				l1.WriteString(strings.Repeat(r, TileAbbrLen))
			}
			if y <= max.Y+border {
				// Left pixel
				r = NoTileStr
				if ct.Type != TileTypeEmpty || lt.Type != TileTypeEmpty {
					wallType := 0
					if ct.Walls[DirLeft] || lt.Walls[DirRight] {
						wallType = DirUp | DirDown
					}
					r, ok = WallStrs[wallType]
					if !ok {
						r = " "
					}
				}
				l2.WriteString(r)
			}
			if x <= max.X+border && y <= max.Y+border {
				// Centre pixel
				r = strings.Repeat(NoTileStr, TileAbbrLen)
				if ct.Type != TileTypeEmpty {
					r = RenderTileAbbr(ct.Type)
				}
				l2.WriteString(r)
			}
		}
		output.WriteString(l1.String())
		if y <= max.Y+border {
			output.WriteRune('\n')
			output.WriteString(l2.String())
		}
	}
	return output.String()
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
