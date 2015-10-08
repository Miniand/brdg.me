package seven_wonders_duel

import (
	"bytes"
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
)

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
	output := &bytes.Buffer{}
	for _, c := range Cards {
		output.WriteString(c.RenderMultiline())
		output.WriteString("\n\n")
	}
	return output.String(), nil
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
		render.Bold(fmt.Sprintf("%s %s", RenderCardType(c.Type), c.Name)),
	}
	if !c.Cost.IsZero() {
		rows = append(rows, RenderCost(c.Cost))
	}
	rows = append(rows, c.RenderSummary())
	return render.CentreLayout(rows, 0)
}

func RenderCost(c cost.Cost) string {
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
