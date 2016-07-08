package lords_of_vegas

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

var CasinoColors = map[int]string{
	CasinoAlbion:  render.Magenta,
	CasinoSphynx:  render.Yellow,
	CasinoVega:    render.Green,
	CasinoTivoli:  render.Gray,
	CasinoPioneer: render.Red,
}

var renderStripLine = `  {{bg "black"}}  {{_bg}}  `
var renderStrip = strings.Join([]string{
	renderStripLine,
	renderStripLine,
	renderStripLine,
	renderStripLine,
}, "\n")

func (g *Game) RenderForPlayer(player string) (string, error) {
	cells := [][]interface{}{}
	for _, layoutRow := range BoardLayout {
		cellsRow := []interface{}{}
		for _, layoutCell := range layoutRow {
			var cell interface{}
			if layoutCell == "ST" {
				cell = renderStrip
			} else if bs, ok := BoardSpaceByLocation[layoutCell]; ok {
				cell = RenderSpace(bs, g.Board[layoutCell])
			} else {
				cell = "\n\n\n"
			}
			cellsRow = append(cellsRow, cell)
		}
		cells = append(cells, cellsRow)
	}
	return render.Table(cells, 0, 0), nil
}

var (
	renderSpaceBorderStr    = " "
	renderSpaceBorderWidth  = 2
	renderSpaceContentWidth = 5
)

func RenderCasinoBg(input string, casino int) string {
	if col, ok := CasinoColors[casino]; ok {
		return render.Bg(input, col)
	}
	return input
}

func RenderSpace(bs BoardSpace, bsState BoardSpaceState) string {
	edge := RenderCasinoBg(strings.Repeat(renderSpaceBorderStr, renderSpaceBorderWidth), bsState.Casino)
	header := []interface{}{
		edge,
		RenderCasinoBg(strings.Repeat(renderSpaceBorderStr, renderSpaceContentWidth), bsState.Casino),
		edge,
	}
	locText := bs.Location
	if bsState.Owned {
		locText = render.Markup(locText, render.PlayerColour(bsState.Owner), true)
	}
	cells := [][]interface{}{
		header,
		{edge, render.Centred(locText), edge},
	}
	contextualRow := ""
	if bsState.Casino == CasinoNone {
		contextualRow = fmt.Sprintf(
			"%s %d",
			RenderPrice(bs.BuildPrice),
			bs.Dice,
		)
	} else if bsState.Owned && bsState.Dice > 0 {
		contextualRow = render.Markup(
			fmt.Sprintf("%d", bsState.Dice),
			render.PlayerColour(bsState.Owner),
			true,
		)
	}
	cells = append(cells, []interface{}{
		edge,
		render.Centred(contextualRow),
		edge,
	})
	cells = append(cells, header)
	return render.Table(cells, 0, 0)
}

func RenderPrice(price int) string {
	return render.Markup(fmt.Sprintf("$%d", price), render.Yellow, true)
}
