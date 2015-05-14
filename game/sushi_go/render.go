package sushi_go

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

const (
	CardColumnTempura = iota
	CardColumnSashimi
	CardColumnDumpling
	CardColumnMaki
	CardColumnNigiri
	CardColumnPudding
	CardColumnChopsticks
)

var CardColumns = []int{
	CardColumnTempura,
	CardColumnSashimi,
	CardColumnDumpling,
	CardColumnMaki,
	CardColumnNigiri,
	CardColumnPudding,
	CardColumnChopsticks,
}

var CardColumnMap = map[int]int{
	CardTempura:      CardColumnTempura,
	CardSashimi:      CardColumnSashimi,
	CardDumpling:     CardColumnDumpling,
	CardMakiRoll3:    CardColumnMaki,
	CardMakiRoll2:    CardColumnMaki,
	CardMakiRoll1:    CardColumnMaki,
	CardSalmonNigiri: CardColumnNigiri,
	CardSquidNigiri:  CardColumnNigiri,
	CardEggNigiri:    CardColumnNigiri,
	CardPudding:      CardColumnPudding,
	CardWasabi:       CardColumnNigiri,
	CardChopsticks:   CardColumnChopsticks,
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	buf := bytes.Buffer{}
	buf.WriteString("{{b}}Hand:{{_b}}\n\n")
	explained := map[int]bool{}
	cells := [][]interface{}{}
	for i, c := range g.Hands[pNum] {
		row := []interface{}{
			render.Colour(fmt.Sprintf("(%d)", i+1), render.Gray),
			RenderCard(c),
		}
		if !explained[c] && CardExplanations[c] != "" {
			row = append(row, render.Colour(
				"  "+CardExplanations[c],
				render.Gray,
			))
			explained[c] = true
			// Only explain for the first maki roll
			if c == CardMakiRoll1 || c == CardMakiRoll2 || c == CardMakiRoll3 {
				explained[CardMakiRoll1] = true
				explained[CardMakiRoll2] = true
				explained[CardMakiRoll3] = true
			}
		}
		cells = append(cells, row)
	}
	buf.WriteString(render.Table(cells, 0, 2))
	buf.WriteString(render.Bold("\n\nYour cards:\n"))
	buf.WriteString(render.Table(CardsCells(g.Played[pNum]), 0, 3))
	buf.WriteString(render.Bold("\n\n\nAll players:"))
	for p := range g.AllPlayers {
		buf.WriteString(fmt.Sprintf(
			"\n\n%s ({{b}}%d{{_b}} points)\n",
			g.RenderName(p),
			g.Points[p],
		))
		buf.WriteString(render.Table(CardsCells(g.Played[p]), 0, 3))
	}
	return buf.String(), nil
}

func CardsCells(cards []int) [][]interface{} {
	columns := map[int][]string{}
	for _, c := range CardColumns {
		columns[c] = []string{}
	}
	unusedWasabi := 0
	for _, c := range cards {
		col := CardColumnMap[c]
		switch c {
		case CardWasabi:
			columns[col] = append(columns[col], RenderCard(c))
			unusedWasabi++
		case CardSalmonNigiri, CardSquidNigiri, CardEggNigiri:
			if unusedWasabi > 0 {
				columns[col][len(columns[col])-unusedWasabi] = fmt.Sprintf(
					"%s + %s",
					RenderCard(c),
					RenderCard(CardWasabi),
				)
				unusedWasabi--
			} else {
				columns[col] = append(columns[col], RenderCard(c))
			}
		default:
			columns[col] = append(columns[col], RenderCard(c))
		}
	}
	cells := [][]interface{}{}
	y := 0
	for {
		row := []interface{}{}
		hadContent := false
		for _, col := range CardColumns {
			l := len(columns[col])
			if l == 0 {
				// Skip empty columns
				continue
			}
			cell := ""
			if l > y {
				cell = columns[col][y]
				hadContent = true
			}
			row = append(row, cell)
		}
		if !hadContent {
			break
		}
		cells = append(cells, row)
		y++
	}
	return cells
}

func RenderCard(c int) string {
	return render.Markup(CardStrings[c], CardColours[c], c != CardPlayed)
}

func RenderCards(cards []int) []string {
	cardStrs := make([]string, len(cards))
	for i, c := range cards {
		cardStrs[i] = RenderCard(c)
	}
	return cardStrs
}
