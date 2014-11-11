package bang_dice

import "github.com/Miniand/brdg.me/render"

func (g *Game) RenderForPlayer(string) (string, error) {
	cells := [][]interface{}{}
	for _, c := range Chars {
		cells = append(cells, []interface{}{
			render.Bold(c.StartingLife()),
			render.Bold(c.Name()),
			render.Colour(c.Description(), "gray"),
		})
	}
	return render.Table(cells, 0, 2), nil
}
