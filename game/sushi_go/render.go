package sushi_go

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

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
	return buf.String(), nil
}

func RenderCard(c int) string {
	return render.Markup(CardStrings[c], CardColours[c], c != CardPlayed)
}
