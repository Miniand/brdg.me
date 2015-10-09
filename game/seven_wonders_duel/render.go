package seven_wonders_duel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/render"
)

const (
	CardTypeText      = "##"
	ProgressTokenText = `{{b}}{{c "green"}}@{{_c}}{{_b}}`
	ExtraTurnText     = `{{b}}{{c "blue"}}&{{_c}}{{_b}}`
	WonderText        = `{{b}}{{c "yellow"}}WOND{{_c}}{{_b}}`
	CardWidth         = 14
	CardSpacing       = 2
)

var FaceDownCardLine = fmt.Sprintf(
	`{{bg "gray"}}%s{{_bg}}`,
	strings.Repeat(" ", CardWidth-4),
)
var FaceDownCard = render.CentreLines(strings.Join(
	[]string{
		FaceDownCardLine,
		FaceDownCardLine,
		FaceDownCardLine,
	},
	"\n",
), CardWidth)

var CardColours = map[int]string{
	CardTypeRaw:          render.Black,
	CardTypeManufactured: render.Gray,
	CardTypeCivilian:     render.Blue,
	CardTypeScientific:   render.Green,
	CardTypeCommercial:   render.Yellow,
	CardTypeMilitary:     render.Red,
	CardTypeGuild:        render.Magenta,
	CardTypeWonder:       render.Cyan,
}

var GoodColours = map[int]string{
	GoodCoin:    render.Yellow,
	GoodWood:    render.Green,
	GoodClay:    render.Red,
	GoodStone:   render.Gray,
	GoodGlass:   render.Cyan,
	GoodPapyrus: render.Yellow,
}

var GoodAbbr = map[int]string{
	GoodCoin:    "Co",
	GoodWood:    "Wo",
	GoodClay:    "Cl",
	GoodStone:   "St",
	GoodGlass:   "Gl",
	GoodPapyrus: "Pa",
}

var ScienceColours = map[int]string{
	ScienceCartography: render.Blue,
	ScienceLaw:         render.Red,
	ScienceAstronomy:   render.Yellow,
	ScienceMathematics: render.Cyan,
	ScienceMedicine:    render.Green,
	ScienceLiterature:  render.Magenta,
	ScienceEngineering: render.Gray,
}

var ScienceStrings = map[int]string{
	ScienceCartography: "Cart",
	ScienceLaw:         "Law",
	ScienceAstronomy:   "Astr",
	ScienceMathematics: "Math",
	ScienceMedicine:    "Medi",
	ScienceLiterature:  "Lite",
	ScienceEngineering: "Engi",
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	return g.RenderLayout(0, g.Layout), nil
}

func (g *Game) RenderLayout(player int, layout Layout) string {
	outputRows := []string{}
	for y, row := range layout {
		rowCells := []interface{}{}
		if y%2 == 1 {
			rowCells = append(rowCells, strings.Repeat(" ", CardWidth/2))
		}
		for x, card := range row {
			if card == 0 {
				rowCells = append(rowCells, strings.Repeat(" ", CardWidth))
			} else if !layout.IsVisible(Loc{x, y}) {
				rowCells = append(rowCells, FaceDownCard)
			} else {
				rowCells = append(rowCells, render.CentreLines(
					Cards[card].RenderMultiline(),
					CardWidth,
				))
			}
		}
		outputRows = append(outputRows, render.Table(
			[][]interface{}{rowCells},
			0,
			CardSpacing,
		))
	}
	return strings.Join(outputRows, "\n\n")
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func RenderCoins(amount int) string {
	return render.Markup(fmt.Sprintf("$%d", amount), render.Yellow, true)
}

func RenderVP(amount int) string {
	return render.Markup(fmt.Sprintf("%dVP", amount), render.Green, true)
}

func RenderCardType(cardType int) string {
	return render.Markup(CardTypeText, CardColours[cardType], true)
}

func (c Card) RenderMultiline() string {
	rows := []interface{}{
		fmt.Sprintf(
			`{{bg "%s"}}{{c "%s"}}{{b}} %s {{_b}}{{_c}}{{_bg}}`,
			CardColours[c.Type],
			render.ColourForBackground(CardColours[c.Type]),
			c.Name,
		),
		RenderCost(c.Cost),
		c.RenderSummary(),
	}
	return render.CentreLayout(rows, 0)
}

func RenderCost(c cost.Cost) string {
	if c.IsZero() {
		return render.Markup("Free", render.Gray, true)
	}
	parts := []string{}
	for _, k := range c.Keys() {
		v := c[k]
		switch k {
		case GoodCoin:
			parts = append(parts, RenderCoins(v))
		default:
			numStr := ""
			if v > 1 {
				numStr = strconv.Itoa(v)
			}
			parts = append(parts, render.Markup(
				fmt.Sprintf("%s%s", numStr, GoodAbbr[k]),
				GoodColours[k],
				true,
			))
		}
	}
	return strings.Join(parts, " ")
}

func RenderGoods(goods []int) []string {
	output := make([]string, len(goods))
	for k, g := range goods {
		output[k] = render.Markup(GoodAbbr[g], GoodColours[g], true)
	}
	return output
}

func RenderProvides(costs []cost.Cost) string {
	parts := []string{}
	for _, c := range costs {
		parts = append(parts, RenderCost(c))
	}
	return strings.Join(parts, " / ")
}

func RenderMilitary(amount int) string {
	return render.Markup(fmt.Sprintf("%dStr", amount), render.Red, true)
}

func RenderScience(science int) string {
	return render.Markup(ScienceStrings[science], ScienceColours[science], true)
}
