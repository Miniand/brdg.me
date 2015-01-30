package alhambra

import (
	"bytes"
	"errors"
	"fmt"
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
	return g.Boards[pNum].Grid.Render(), nil
}

func RenderTileAbbr(tileType int) string {
	return render.Markup(TileAbbrs[tileType], TileColours[tileType], true)
}

func (g Grid) Render() string {
	var ok bool
	min, max := g.Bounds()
	output := bytes.NewBuffer([]byte{})
	for y := min.Y - 1; y <= max.Y+1; y++ {
		l1 := bytes.NewBuffer([]byte{})
		l2 := bytes.NewBuffer([]byte{})
		for x := min.X - 1; x <= max.X+1; x++ {
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
			// Centre pixel
			r = strings.Repeat(NoTileStr, TileAbbrLen)
			if ct.Type != TileTypeEmpty {
				r = RenderTileAbbr(ct.Type)
			}
			l2.WriteString(r)
		}
		output.WriteString(l1.String())
		output.WriteString(fmt.Sprintf("%s\n", NoTileStr))
		output.WriteString(l2.String())
		output.WriteString(fmt.Sprintf("%s\n", NoTileStr))
	}
	output.WriteString(strings.Repeat(NoTileStr, (max.X-min.X+3)*(TileAbbrLen+1)+1))
	return output.String()
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
