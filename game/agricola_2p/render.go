package agricola_2p

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

const (
	TileWidth = 12
)

var EmptyTileLine = strings.Repeat(" ", TileWidth)

var Corner = `{{bg "green"}}  {{_bg}}`

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	return g.PBoards[pNum].Render(), nil
}

func (pb *PBoard) Render() string {
	lines := []string{}
	for y := 0; y < 4; y++ {
		for tLine := 0; tLine < 6; tLine++ {
			if y == 3 && tLine > 0 {
				// Bottom border, we don't care about the tile.
				break
			}
			row := bytes.Buffer{}
			for x := pb.XStart; x <= pb.XEnd+1; x++ {
				l := Loc{x, y}
				if tLine == 0 {
					// Top border
					row.WriteString(Corner)
					if x <= pb.XEnd {
						borderCol := render.Gray
						if pb.Tiles.Border(l, Up) {
							borderCol = render.Yellow
						}
						row.WriteString(fmt.Sprintf(
							`{{bg "%s"}}%s{{_bg}}`,
							borderCol,
							EmptyTileLine,
						))
					}
				} else {
					// Tile Body
					borderCol := render.Gray
					if pb.Tiles.Border(l, Left) {
						borderCol = render.Yellow
					}
					row.WriteString(fmt.Sprintf(`{{bg "%s"}}  {{_bg}}`, borderCol))
					if x <= pb.XEnd {
						t := pb.Tiles.At(l)
						tileCol := render.Green
						text := EmptyTileLine
						if t.Building != nil {
							tileCol = render.Black
							if tLine == 2 {
								text = render.Centre(
									render.Bold(t.Building.String()),
									TileWidth,
								)
							}
							capacity := t.Building.Capacity()
							if tLine == 5 && capacity > 0 {
								text = render.Right(render.Markup(
									strconv.Itoa(capacity),
									render.Red,
									true,
								), TileWidth)
							}
						}
						row.WriteString(fmt.Sprintf(
							`{{bg "%s"}}%s{{_bg}}`,
							tileCol,
							text,
						))
					}
				}
			}
			lines = append(lines, row.String())
		}
	}
	return strings.Join(lines, "\n")
}
