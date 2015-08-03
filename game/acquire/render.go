package acquire

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	output := bytes.NewBufferString("")
	output.WriteString("{{b}}Board:{{_b}}\n\n")

	boardRows := []string{}
	for _, r := range Rows() {
		row := bytes.Buffer{}
		for _, c := range Cols() {
			width := 2
			if c >= BOARD_COL_10 {
				width = 3
			}
			t := Tile{r, c}
			if _, n := g.PlayerTiles[pNum].Remove(t, 1); n > 0 {
				// Player has this tile.
				row.WriteString(
					fmt.Sprintf(`{{c "gray"}}{{b}}%s{{_b}}{{_c}} `, TileText(t)))
			} else {
				val, ok := g.TileAt(t)
				if !ok {
					return "", errors.New("somehow iterated to nonexistant tile")
				}
				horSpacer := " "
				switch val {
				case TILE_DISCARDED:
					row.WriteString(strings.Repeat(" ", width))
				case TILE_EMPTY:
					row.WriteString(
						render.Colour(strings.Repeat("-", width), render.Gray))
				case TILE_UNINCORPORATED:
					row.WriteString(`{{bg "gray"}}`)
					row.WriteString(strings.Repeat(" ", width))
					row.WriteString(`{{_bg}}`)
				default:
					row.WriteString(fmt.Sprintf(`{{bg "%s"}}`, CorpColours[val]))
					row.WriteString(strings.Repeat(" ", width))
					row.WriteString(`{{_bg}}`)
					if adjV, ok := g.TileAt(Tile{t.Row, t.Column + 1}); ok && adjV == val {
						horSpacer = fmt.Sprintf(
							`{{bg "%s"}} {{_bg}}`,
							CorpColours[val],
						)
					}
				}
				row.WriteString(horSpacer)
			}
		}
		boardRows = append(boardRows, row.String())
	}
	output.WriteString(strings.Join(boardRows, "\n"))
	output.WriteString("\n\n")

	// Board
	cells := [][]interface{}{}
	for _, r := range Rows() {
		row := []interface{}{}
		for _, c := range Cols() {
			cellOutput := g.RenderTile(Tile{r, c})
			// We embolden the tile if the player has it in their hand
			t := Tile{
				Row:    r,
				Column: c,
			}
			if _, n := g.PlayerTiles[pNum].Remove(t, 1); n > 0 {
				cellOutput = fmt.Sprintf(`{{c "gray"}}{{b}}%s{{_b}}{{_c}}`,
					TileText(t))
			}
			row = append(row, cellOutput)
		}
		cells = append(cells, row)
	}
	boardOutput := render.Table(cells, 0, 1)
	output.WriteString(boardOutput)
	// Hand
	handTiles := []string{}
	for _, tRaw := range g.PlayerTiles[pNum].Sort() {
		t := tRaw.(Tile)
		handTiles = append(handTiles, TileText(t))
	}
	output.WriteString(fmt.Sprintf(
		"\n\n{{b}}Your tiles: {{c \"gray\"}}%s{{_c}}{{_b}}\n",
		strings.Join(handTiles, " ")))
	output.WriteString(fmt.Sprintf(
		"{{b}}Your cash:  $%d{{_b}}\n", g.PlayerCash[pNum]))
	output.WriteString(fmt.Sprintf(
		"{{b}}Tiles left: %d{{_b}}", len(g.BankTiles)))
	// Corp table
	cells = [][]interface{}{
		[]interface{}{
			"{{b}}Corporation{{_b}}",
			"{{b}}Size{{_b}}",
			"{{b}}Value{{_b}}",
			"{{b}}Shares{{_b}}",
			"{{b}}Major{{_b}}",
			"{{b}}Minor{{_b}}",
		},
	}
	for _, c := range Corps() {
		cells = append(cells, []interface{}{
			fmt.Sprintf(`{{b}}%s{{_b}}`, RenderCorpWithShort(c)),
			fmt.Sprintf("%d", g.CorpSize(c)),
			fmt.Sprintf("$%d", g.CorpValue(c)),
			fmt.Sprintf("%d left", g.BankShares[c]),
			fmt.Sprintf("$%d", g.Corp1stBonus(c)),
			fmt.Sprintf("$%d", g.Corp2ndBonus(c)),
		})
	}
	corpOutput := render.Table(cells, 0, 2)
	output.WriteString("\n\n")
	output.WriteString(corpOutput)
	// Player table
	playerHeadings := []interface{}{
		"{{b}}Player{{_b}}",
		"{{b}}Cash{{_b}}",
	}
	for _, corp := range Corps() {
		playerHeadings = append(playerHeadings, fmt.Sprintf(
			"{{b}}%s{{_b}}", RenderCorpShort(corp)))
	}
	cells = [][]interface{}{
		playerHeadings,
	}
	for pNum, p := range g.Players {
		row := []interface{}{
			fmt.Sprintf("{{b}}%s{{_b}}", render.PlayerName(pNum, p)),
			fmt.Sprintf("$%d", g.PlayerCash[pNum]),
		}
		for _, corp := range Corps() {
			row = append(row, fmt.Sprintf("%d", g.PlayerShares[pNum][corp]))
		}
		cells = append(cells, row)
	}
	playerOutput := render.Table(cells, 0, 2)
	output.WriteString("\n\n")
	output.WriteString(playerOutput)
	return output.String(), nil
}

func (g *Game) RenderTile(t Tile) (output string) {
	val := g.Board[t.Row][t.Column]
	switch val {
	case TILE_DISCARDED:
		output = "  "
	case TILE_EMPTY:
		output = `{{c "gray"}}--{{_c}}`
	case TILE_UNINCORPORATED:
		output = `{{c "gray"}}##{{_c}}`
	default:
		output = fmt.Sprintf(`{{b}}%s{{_b}}`, RenderCorpShort(val))
	}
	return
}
