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
				cell = RenderSpace(bs, rnd.Int()%CasinoTheStrip, rnd.Int()%5-1, rnd.Int()%7)
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

func RenderSpace(bs BoardSpace, casino, owner, dice int) string {
	edge := RenderCasinoBg(strings.Repeat(renderSpaceBorderStr, renderSpaceBorderWidth), casino)
	header := []interface{}{
		edge,
		RenderCasinoBg(strings.Repeat(renderSpaceBorderStr, renderSpaceContentWidth), casino),
		edge,
	}
	locText := render.Bold(bs.Location)
	if owner != -1 {
		locText = render.Colour(locText, render.PlayerColour(owner))
	}
	cells := [][]interface{}{
		header,
		{edge, render.Centred(locText), edge},
	}
	contextualRow := ""
	if casino == CasinoNone {
		contextualRow = RenderPrice(bs.BuildPrice)
	} else if owner != -1 && dice > 0 {
		contextualRow = render.Markup(
			fmt.Sprintf("%d", dice),
			render.PlayerColour(owner),
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
