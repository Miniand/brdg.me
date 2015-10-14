package seven_wonders_duel

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/render"
)

const (
	CardTypeText      = "##"
	ProgressTokenText = `{{b}}{{c "green"}}PR{{_c}}{{_b}}`
	ExtraTurnText     = `{{b}}{{c "blue"}}+1{{_c}}{{_b}}`
	WonderText        = `{{b}}{{c "yellow"}}WO{{_c}}{{_b}}`
	LinkBuildText     = `{{b}}{{c "cyan"}}->{{_c}}{{_b}}`
	CardWidth         = 14
	CardSpacing       = 2
)

var CardColours = map[int]string{
	CardTypeRaw:          render.Black,
	CardTypeManufactured: render.Gray,
	CardTypeCivilian:     render.Blue,
	CardTypeScientific:   render.Green,
	CardTypeCommercial:   render.Yellow,
	CardTypeMilitary:     render.Red,
	CardTypeGuild:        render.Magenta,
	CardTypeWonder:       render.Yellow,
	CardTypeProgress:     render.Green,
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

var GoodNames = map[int]string{
	GoodCoin:    "Coin",
	GoodWood:    "Wood",
	GoodClay:    "Clay",
	GoodStone:   "Ston",
	GoodGlass:   "Glas",
	GoodPapyrus: "Papy",
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

var AgeColours = map[int]string{
	1: render.Red,
	2: render.Blue,
	3: render.Magenta,
}

func SolidLine(colour string, width int) string {
	return fmt.Sprintf(
		`{{bg "%s"}}%s{{_bg}}`,
		colour,
		strings.Repeat(" ", width),
	)
}

func CardBack(colour string) string {
	width := CardWidth - 4
	return strings.Join([]string{
		SolidLine(render.Gray, width),
		strings.Join([]string{
			SolidLine(render.Gray, 2),
			SolidLine(colour, width-4),
			SolidLine(render.Gray, 2),
		}, ""),
		SolidLine(render.Gray, width),
	}, "\n")
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	oNum := Opponent(pNum)
	rows := []interface{}{}
	if g.Phase == PhaseChooseWonder {
		wonderOutputs := []interface{}{}
		for _, w := range g.AvailableWonders() {
			wonderOutputs = append(wonderOutputs, Cards[w].RenderMultiline())
		}
		rows = append(rows, render.Bold("Available wonders"), "", render.Table(
			[][]interface{}{wonderOutputs},
			0,
			3,
		))
	} else {
		rows = append(
			rows,
			render.Bold(fmt.Sprintf("Age %d", g.Age)),
			"",
			g.RenderLayout(0, g.Layout),
		)
	}

	// Progress tokens
	progCells := []interface{}{}
	for _, p := range g.ProgressTokens {
		progCells = append(progCells, Cards[p].RenderMultiline())
	}
	rows = append(
		rows,
		"",
		"",
		render.Bold(fmt.Sprintf("Available progress tokens (%s)", ProgressTokenText)),
		"",
		render.Table([][]interface{}{progCells}, 0, 3),
	)

	// Player table
	cells := [][]interface{}{
		{
			render.Centred("              \n" + g.RenderPlayerNotables(pNum)),
			render.Centred(g.RenderPlayerTable(pNum)),
			render.Centred("              \n" + g.RenderPlayerNotables(oNum)),
		},
	}
	rows = append(
		rows,
		"",
		"",
		render.Bold("Player tableaus"),
		render.Table(cells, 0, 5),
	)

	// Glossary
	rows = append(
		rows,
		"",
		render.Bold("Glossary"),
		render.Table([][]interface{}{
			{render.RightAligned(LinkBuildText), "Makes another card free (link build)"},
			{render.RightAligned(ExtraTurnText), "Extra turn after this one"},
			{render.RightAligned(fmt.Sprintf(
				"%s ^ %s",
				RenderVP(1),
				RenderCardType(CardTypeMilitary),
			)), fmt.Sprintf(
				"%s for each %s in the city with the most %s",
				RenderVP(1),
				RenderCardType(CardTypeMilitary),
				RenderCardType(CardTypeMilitary),
			)},
		}, 0, 3),
	)

	return render.CentreLayout(rows, 0), nil
}

func RenderCheapens(goods []int) string {
	return fmt.Sprintf(
		"%s costs %s",
		strings.Join(RenderGoods(goods), " "),
		RenderCoins(1),
	)
}

func (g *Game) RenderPlayerTable(player int) string {
	opp := Opponent(player)
	cells := [][]interface{}{
		{
			render.Centred(g.PlayerName(player)),
			"",
			render.Centred(g.PlayerName(opp)),
		},
		{
			render.Centred(render.Bold(strconv.Itoa(g.PlayerCoins[player]))),
			render.Centred(render.Markup("Coin", render.Yellow, true)),
			render.Centred(strconv.Itoa(g.PlayerCoins[opp])),
		},
		{
			render.Centred(render.Bold(strconv.Itoa(g.PlayerVP(player)))),
			render.Centred(render.Markup("VP", render.Green, true)),
			render.Centred(strconv.Itoa(g.PlayerVP(opp))),
		},
		{
			render.Centred(render.Bold(strconv.Itoa(g.PlayerCardTypeCount(player, CardTypeWonder)))),
			render.Centred(WonderText),
			render.Centred(strconv.Itoa(g.PlayerCardTypeCount(opp, CardTypeWonder))),
		},
		{
			render.Centred(render.Bold(strconv.Itoa(g.PlayerCardTypeCount(player, CardTypeProgress)))),
			render.Centred(ProgressTokenText),
			render.Centred(strconv.Itoa(g.PlayerCardTypeCount(opp, CardTypeProgress))),
		},
		{},
	}
	for _, good := range []int{
		GoodWood,
		GoodStone,
		GoodClay,
		GoodPapyrus,
		GoodGlass,
	} {
		cells = append(cells, []interface{}{
			render.Centred(render.Bold(g.RenderPlayerGoodCount(player, good))),
			render.Centred(RenderGoodName(good)),
			render.Centred(g.RenderPlayerGoodCount(opp, good)),
		})
	}
	cells = append(cells, []interface{}{})
	for _, ct := range []int{
		CardTypeRaw,
		CardTypeManufactured,
		CardTypeCivilian,
		CardTypeScientific,
		CardTypeCommercial,
		CardTypeMilitary,
		CardTypeGuild,
	} {
		cells = append(cells, []interface{}{
			render.Centred(render.Bold(strconv.Itoa(g.PlayerCardTypeCount(player, ct)))),
			render.Centred(RenderCardType(ct)),
			render.Centred(strconv.Itoa(g.PlayerCardTypeCount(opp, ct))),
		})
	}
	return render.Table(cells, 0, 2)
}

func (g *Game) RenderPlayerNotables(player int) string {
	rows := []interface{}{}
	if len(g.PlayerWonders[player]) > 0 {
		rows = append(rows, g.RenderUnbuiltWondersTable(g.PlayerWonders[player]))
	}
	ongoingEffects := []string{}
	for _, c := range g.PlayerCards[player] {
		if Cards[c].OngoingEffect != "" {
			ongoingEffects = append(ongoingEffects, Cards[c].OngoingEffect)
		}
	}
	if len(ongoingEffects) > 0 {
		rows = append(rows, render.Colour("Bonuses", render.Gray))
		rows = append(rows, strings.Join(ongoingEffects, "\n"))
	}
	return render.CentreLayout(rows, 1)
}

func (g *Game) RenderPlayerGoodCount(player, good int) string {
	base, extra := g.PlayerGoodCount(player, good)
	extraStr := ""
	if extra > 0 {
		extraStr += fmt.Sprintf("+%d", extra)
	}
	return fmt.Sprintf("%d%s", base, extraStr)
}

func (g *Game) RenderUnbuiltWondersTable(wonders []int) string {
	cells := [][]interface{}{{render.Centred(render.Colour("Unbuilt WO", render.Gray))}}
	for _, w := range wonders {
		cells = append(
			cells,
			[]interface{}{render.Centred(Cards[w].RenderMultiline())},
		)
	}
	return render.Table(cells, 1, 0)
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
				rowCells = append(rowCells, render.CentreLines(
					CardBack(AgeColours[g.Age]),
					CardWidth,
				))
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

func CardText(cardType int) string {
	switch cardType {
	case CardTypeWonder:
		return WonderText
	case CardTypeProgress:
		return ProgressTokenText
	default:
		return CardTypeText
	}
}

func RenderCardType(cardType int) string {
	return render.Markup(CardText(cardType), CardColours[cardType], true)
}

func (c Card) RenderMultiline() string {
	fg := CardColours[c.Type]
	bg := render.Black
	if c.Type != CardTypeWonder && c.Type != CardTypeProgress {
		fg = render.ColourForBackground(CardColours[c.Type])
		bg = CardColours[c.Type]
	}
	rows := []interface{}{fmt.Sprintf(
		`{{bg "%s"}}{{c "%s"}}{{b}} %s {{_b}}{{_c}}{{_bg}}`,
		bg,
		fg,
		c.Name,
	)}
	if c.Type != CardTypeProgress {
		rows = append(rows, RenderCost(c.Cost))
	}
	rows = append(rows, c.RenderSummary())
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

func RenderGoodName(good int) string {
	return render.Markup(GoodNames[good], GoodColours[good], true)
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
