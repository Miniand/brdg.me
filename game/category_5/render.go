package category_5

import (
	"bytes"
	"fmt"

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
		total := 0
		for j := 0; j < 5; j++ {
			var cell interface{}
			cell = "  "
			if j < len(b) {
				cell = b[j]
				total += b[j].Heads()
			}
			row = append(row, cell)
		}
		row = append(row, fmt.Sprintf("  %d pts", total))
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
	return buf.String(), nil
}
