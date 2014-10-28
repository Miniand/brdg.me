package category_5

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	// Board
	cells := [][]interface{}{}
	for i, b := range g.Board {
		row := []interface{}{
			fmt.Sprintf("{{b}}#%d{{_b}}  ", i+1),
		}
		for j := 0; j < 5; j++ {
			var cell interface{}
			cell = "  "
			if j < len(b) {
				cell = b[j]
			}
			row = append(row, cell)
		}
		row = append(row, fmt.Sprintf("  %d pts", CardsHeads(b)))
		cells = append(cells, row)
	}
	buf.WriteString(render.Table(cells, 0, 2))
	// Hand
	if len(g.Hands[pNum]) > 0 {
		buf.WriteString("\n\n")
		row := []interface{}{"{{b}}Your hand:{{_b}}"}
		for _, c := range g.Hands[pNum] {
			row = append(row, c)
		}
		buf.WriteString(render.Table([][]interface{}{row}, 0, 2))
	}
	// Legend
	buf.WriteString("\n\n")
	parts := []string{}
	for _, i := range []int{1, 2, 3, 5, 7} {
		parts = append(parts, render.Markup(
			fmt.Sprintf("%d pts", i),
			CardColours[i],
			true,
		))
	}
	buf.WriteString(fmt.Sprintf(
		"{{b}}Legend:{{_b}} %s",
		strings.Join(parts, ", "),
	))
	// Score table
	buf.WriteString("\n\n")
	cells = [][]interface{}{{
		render.Bold("Players"),
		render.Bold("Taken"),
		render.Bold("Pts"),
	}}
	for p, _ := range g.Players {
		cells = append(cells, []interface{}{
			g.RenderName(p),
			render.Centred(strconv.Itoa(len(g.PlayerCards[p]))),
			render.Centred(render.Bold(g.Points[p])),
		})
	}
	buf.WriteString(render.Table(cells, 0, 2))
	return buf.String(), nil
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
