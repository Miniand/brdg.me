package cathedral

import (
	"bytes"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	cells := [][]interface{}{
		{"one\ntwo\nthree", render.Centred("four\nfiverlol")},
		{render.CellSpan{
			Content: render.Centred("This one should\nbe quite long"),
			Cols:    2,
		}},
		{"six\nseven", render.RightAligned("eight\nnine\nten\neleven")},
	}
	buf.WriteString(render.Table(cells, 0, 0))
	buf.WriteString("after")
	return buf.String(), nil
}
