package alhambra

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/render"
)

var NoTileStr = `{{c "gray"}}▒{{_c}}`

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

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	output := bytes.NewBuffer([]byte{})
	// Current player board
	output.WriteString(g.RenderPlayerGrid(pNum))
	if len(g.Boards[pNum].Place) > 0 {
		output.WriteString(render.Bold("\n\nTiles to be placed:\n"))
		output.WriteString(g.RenderTiles(
			g.Boards[pNum].Place,
			func(i int) string {
				t := g.Boards[pNum].Place[i]
				if t.Type == TileTypeEmpty {
					return ""
				}
				return render.Colour(strconv.Itoa(i+1), render.Gray)
			},
		))
	}
	if len(g.Boards[pNum].Reserve) > 0 {
		output.WriteString("\n\n")
		output.WriteString(g.RenderReserve(pNum))
	}
	output.WriteString(render.Bold("\n\nYour cards:\n"))
	output.WriteString(RenderCards(g.Boards[pNum].Cards))
	// Purchase tiles
	output.WriteString("\n\n")
	output.WriteString(render.Rule)
	output.WriteString(fmt.Sprintf(
		"\n\nIs is currently {{b}}round %d{{_b}}",
		g.Round,
	))
	output.WriteString(render.Bold(fmt.Sprintf(
		"\n\nTiles available for purchase (%d remaining):\n",
		len(g.TileBag),
	)))
	output.WriteString(g.RenderTiles(g.Tiles, func(i int) string {
		t := g.Tiles[i]
		if t.Type == TileTypeEmpty {
			return ""
		}
		c := Card{i, t.Cost}
		return c.String()
	}))
	// Draw cards
	output.WriteString(render.Bold("\n\nMoney available for taking:\n"))
	output.WriteString(RenderCards(g.Cards))
	// Player table
	header := []interface{}{render.Bold("Player")}
	for _, t := range ScoringTileTypes {
		header = append(header, RenderTileAbbr(t))
	}
	header = append(
		header,
		render.Bold("Wall"),
		render.Bold("Cards"),
		render.Bold("Pts"),
	)
	cells := [][]interface{}{header}
	for pNum, _ := range g.Players {
		row := []interface{}{g.PlayerName(pNum)}
		counts := g.Boards[pNum].TileCounts()
		for _, t := range ScoringTileTypes {
			row = append(row, render.Centred(strconv.Itoa(counts[t])))
		}
		row = append(
			row,
			render.Centred(strconv.Itoa(g.Boards[pNum].Grid.LongestExtWall())),
			render.Centred(strconv.Itoa(len(g.Boards[pNum].Cards))),
			render.Centred(strconv.Itoa(g.Boards[pNum].Points)),
		)
		cells = append(cells, row)
	}
	output.WriteString("\n\n")
	output.WriteString(render.Table(cells, 0, 1))
	// Other player boards
	for p := range g.Players {
		if p == pNum {
			continue
		}
		output.WriteString(fmt.Sprintf(
			"\n\n%s\n\nPlayer board for %s\n\n",
			render.Rule,
			g.PlayerName(p),
		))
		output.WriteString(g.RenderPlayerGrid(p))
		if len(g.Boards[p].Reserve) > 0 {
			output.WriteString("\n\n")
			output.WriteString(g.RenderReserve(p))
		}
	}
	return output.String(), nil
}

func RenderCards(cards card.Deck) string {
	cardStrs := []string{}
	for _, c := range cards.Sort() {
		cardStrs = append(cardStrs, c.(Card).String())
	}
	return strings.Join(cardStrs, "  ")
}

func (g *Game) RenderPlayerGrid(player int) string {
	return AddCoordsToGrid(g.Boards[player].Grid.Render(1))
}

func (g *Game) RenderReserve(player int) string {
	output := bytes.NewBufferString(render.Bold("Reserved tiles:\n"))
	output.WriteString(g.RenderTiles(
		g.Boards[player].Reserve,
		func(i int) string {
			return render.Markup(strconv.Itoa(i+1), render.Gray, true)
		},
	))
	return output.String()
}

func (g *Game) RenderTiles(tiles []Tile, footer func(i int) string) string {
	gr := Grid{}
	for i, t := range tiles {
		gr[Vect{i * 2, 0}] = t
	}
	output := bytes.NewBufferString(gr.Render(0))
	if footer != nil {
		output.WriteRune('\n')
		output.WriteString(HeaderRow(len(tiles)*2-1, func(i int) string {
			if i%2 == 1 {
				return ""
			}
			return footer(i / 2)
		}))
	}
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
			n := render.Markup(strconv.Itoa((i+1)/2), render.Gray, true)
			lines[i] = fmt.Sprintf(
				`%s %s %s`,
				render.Right(n, 3),
				l,
				render.Padded(n, 3),
			)
		}
	}
	// Top and bottom
	header := "    " + render.Markup(
		HeaderRowAlpha(width/TileWidth),
		render.Gray,
		true,
	)
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
