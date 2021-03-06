package seven_wonders

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/render"
)

var ResourceColours = map[int]string{
	GoodCoin:    render.Yellow,
	GoodWood:    render.Green,
	GoodStone:   render.Gray,
	GoodOre:     render.Black,
	GoodClay:    render.Red,
	GoodPapyrus: render.Cyan,
	GoodTextile: render.Magenta,
	GoodGlass:   render.Blue,

	CardKindRaw:          render.Black,
	CardKindManufactured: render.Gray,
	CardKindCivilian:     render.Blue,
	CardKindScientific:   render.Green,
	CardKindCommercial:   render.Yellow,
	CardKindMilitary:     render.Red,
	CardKindGuild:        render.Magenta,
	CardKindWonder:       render.Yellow,

	FieldMathematics: render.Red,
	FieldEngineering: render.Green,
	FieldTheology:    render.Blue,

	AttackStrength: render.Red,
	TokenVictory:   render.Green,
	TokenDefeat:    render.Gray,

	VP: render.Green,
}

const (
	CardSymbol = "###"
)

var CanBuySymbol = render.Markup("✔", render.Green, true)
var CanBuyWithCoinSymbol = render.Markup(
	ResourceSymbols[GoodCoin],
	render.Yellow,
	true,
)
var CannotBuySymbol = render.Markup("✘", render.Red, true)

var ResourceSymbols = map[int]string{
	GoodCoin:    "Coin",
	GoodWood:    "Wood",
	GoodStone:   "Ston",
	GoodOre:     "Ore",
	GoodClay:    "Clay",
	GoodPapyrus: "Papy",
	GoodTextile: "Text",
	GoodGlass:   "Glas",

	CardKindRaw:          CardSymbol,
	CardKindManufactured: CardSymbol,
	CardKindCivilian:     CardSymbol,
	CardKindScientific:   CardSymbol,
	CardKindCommercial:   CardSymbol,
	CardKindMilitary:     CardSymbol,
	CardKindGuild:        CardSymbol,
	CardKindWonder:       "▲",

	FieldMathematics: "Math",
	FieldEngineering: "Engi",
	FieldTheology:    "Theo",

	AttackStrength: "Str",
	TokenVictory:   "Vic",
	TokenDefeat:    "Def",

	VP: "VP",
}

func String(in interface{}) string {
	if str, ok := in.(string); ok {
		return str
	}
	if str, ok := in.(fmt.Stringer); ok {
		return str.String()
	}
	return fmt.Sprintf("%v", in)
}

func RenderResourceColour(str interface{}, resource int, bold bool) string {
	return render.Markup(String(str), ResourceColours[resource], bold)
}

func RenderResourceSymbol(resource int) string {
	return RenderResourceColour(ResourceSymbols[resource], resource, true)
}

func RenderResourceWithSymbol(text string, resource int) string {
	return RenderResourceColour(
		fmt.Sprintf("%s %s", text, ResourceSymbols[resource]),
		resource,
		true,
	)
}

func RenderCardName(c Carder) string {
	ca := c.GetCard()
	return strings.Join([]string{
		RenderResourceSymbol(ca.Kind),
		ca.Name,
	}, " ")
}

type SuppStringer interface {
	SuppString() string
}

func RenderCard(c Carder) string {
	parts := []string{RenderCardName(c.(Carder))}
	if ss, ok := c.(SuppStringer); ok {
		parts = append(parts, ss.SuppString())
	}
	return strings.Join(parts, "  ")
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	output := bytes.NewBufferString(fmt.Sprintf(
		"{{b}}Your city: %s{{_b}}  %s\n\n",
		g.Cities[pNum].Name,
		RenderResourceList(g.Cities[pNum].GoodsProduced()[0].Keys(), " "),
	))

	if len(g.ToResolve) > 0 {
		// Resolver output
		resOut := g.ToResolve[0].String(pNum, g)
		if resOut != "" {
			output.WriteString(resOut)
			output.WriteString("\n\n")
		}
	} else if g.Actions[pNum] != nil {
		// Action output
		actOut := g.Actions[pNum].Output(pNum, g)
		if actOut != "" {
			output.WriteString(render.Bold("Your action: "))
			output.WriteString(actOut)
			output.WriteString("\n\n")
		}
	}

	// Hand
	output.WriteString(render.Bold("Your hand:\n\n"))
	output.WriteString(
		g.RenderCardList(pNum, g.Hands[pNum], true, g.CanAction(pNum), true))

	// Wonders
	output.WriteString(render.Bold("\n\nRemaining wonder stages:\n\n"))
	remaining := g.RemainingWonderStages(pNum)
	if len(remaining) > 0 {
		output.WriteString(
			g.RenderCardList(pNum, remaining, true, false, true))
	} else {
		output.WriteString(
			render.Markup("All wonder stages complete.", render.Gray, false))
	}

	// Discard
	output.WriteString(fmt.Sprintf(
		"\n\n{{b}}Discard pile:{{_b}} %d",
		len(g.Discard),
	))

	// Stats table
	output.WriteString("\n\n")
	output.WriteString(g.RenderStatTable(pNum))

	// Played cards
	output.WriteString("\n\n{{b}}Your tableau:{{_b}}\n\n")
	output.WriteString(g.RenderCardList(pNum, g.Cards[pNum].Sort(), false, false, false))
	return output.String(), nil
}

func (g *Game) RenderCardList(
	player int,
	cards card.Deck,
	showCost, showNums, showAfford bool,
) string {
	output := bytes.NewBuffer([]byte{})
	for i, c := range cards {
		crd := c.(Carder)
		costStrs := []string{CostString(crd.GetCard().Cost)}
		for _, f := range crd.GetCard().FreeWith {
			costStrs = append(costStrs, RenderCardName(Cards[f]))
		}
		numStr := ""
		if showNums {
			numStr = fmt.Sprintf(
				"(%s) ",
				render.Markup(strconv.Itoa(i+1), render.Gray, true),
			)
		}
		afford := ""
		if showAfford {
			cardCoin := crd.GetCard().Cost[GoodCoin]
			canBuild, coins := g.CanBuildCard(player, crd)
			switch {
			case !canBuild:
				afford = CannotBuySymbol
			case canBuild && cardCoin == 0 && len(coins) == 0:
				afford = CanBuySymbol
			default:
				extraCoin := []int{}
				if len(coins) == 0 {
					extraCoin = append(extraCoin, cardCoin)
				}
				for _, perm := range coins {
					sum := cardCoin
					for _, coin := range perm {
						sum += coin
					}
					extraCoin = append(extraCoin, sum)
				}
				afford = RenderResourceColour(
					strconv.Itoa(helper.IntMin(extraCoin...)),
					GoodCoin,
					true,
				)
			}
			afford += " "
		}
		output.WriteString(fmt.Sprintf(
			"%s%s%s\n",
			numStr,
			afford,
			RenderCard(crd),
		))
		if showCost {
			output.WriteString(fmt.Sprintf(
				"          Cost: %s\n",
				strings.Join(costStrs, render.Colour("  or  ", render.Gray)),
			))
		}
		for fi, f := range crd.GetCard().MakesFree {
			prefix := "          Leads to: "
			if fi > 0 {
				prefix = "                    "
			}
			output.WriteString(fmt.Sprintf(
				"%s%s\n",
				prefix,
				RenderCard(Cards[f]),
			))
		}
		output.WriteString("\n")
	}
	return strings.TrimSpace(output.String())
}

func RenderMoney(n int) string {
	return RenderResourceWithSymbol(
		fmt.Sprintf("%d", n),
		GoodCoin,
	)
}

func RenderVP(n int) string {
	return RenderResourceWithSymbol(
		fmt.Sprintf("%d", n),
		VP,
	)
}

func RenderResourceList(goods []int, sep string) string {
	goodStrs := []string{}
	for _, g := range goods {
		goodStrs = append(goodStrs, RenderResourceSymbol(g))
	}
	return strings.Join(goodStrs, sep)
}

func RenderResources(resources cost.Cost, sep string) string {
	parts := []string{}
	for _, r := range resources.Trim().Keys() {
		parts = append(parts, RenderResourceWithSymbol(
			strconv.Itoa(resources[r]),
			r,
		))
	}
	return strings.Join(parts, sep)
}

func RenderDirections(directions []int) string {
	dirStrs := []string{}
	for _, d := range directions {
		dirStrs = append(dirStrs, DirStrings[d])
	}
	return render.Bold(strings.Join(dirStrs, " "))
}

func CostString(c cost.Cost) string {
	if len(c) == 0 {
		return render.Markup("free", render.Gray, true)
	}
	n := 0
	l := len(c)
	count := 0
	parts := []string{}
	for count < l {
		if amount, ok := c[n]; ok {
			part := ""
			switch n {
			case GoodCoin:
				part = RenderMoney(amount)
			default:
				if amount == 1 {
					part = RenderResourceSymbol(n)
				} else {
					part = RenderResourceWithSymbol(strconv.Itoa(amount), n)
				}
			}
			parts = append(parts, part)
			count++
		}
		n++
	}
	return strings.Join(parts, " ")
}

func (g *Game) RenderDeal(player int, deal map[int]int) string {
	parts := []string{}
	for _, dir := range DirNeighbours {
		if deal[dir] != 0 {
			parts = append(parts, fmt.Sprintf(
				"%s %s",
				g.PlayerName(g.NumFromPlayer(player, dir)),
				RenderMoney(deal[dir]),
			))
		}
	}
	return fmt.Sprintf("pay %s", render.CommaList(parts))
}

func (g *Game) RenderDeals(player int, deals []map[int]int) string {
	parts := []string{}
	for i, d := range deals {
		parts = append(parts, fmt.Sprintf(
			"({{b}}%d{{_b}}) %s",
			i+1,
			g.RenderDeal(player, d),
		))
	}
	return strings.Join(parts, "\n")
}

func (g *Game) RenderStatTable(player int) string {
	cells := [][]interface{}{}
	sections := []struct {
		Heading   string
		Resources []int
	}{
		{
			"",
			[]int{GoodCoin, VP, CardKindWonder},
		},
		{
			"Raw goods",
			RawGoods,
		},
		{
			"Manufactured goods",
			ManufacturedGoods,
		},
		{
			"Military",
			[]int{AttackStrength, TokenVictory, TokenDefeat},
		},
		{
			"Science",
			Fields,
		},
		{
			"Cards",
			CardKinds,
		},
	}
	pLen := len(g.Players)
	fromPlayer := g.NumFromPlayer(player, -(pLen-1)/2)
	heading := []interface{}{""}
	for i := 0; i < pLen; i++ {
		p := g.NumFromPlayer(fromPlayer, i)
		heading = append(heading, g.PlayerName(p))
	}
	cells = append(cells, heading)
	for _, s := range sections {
		if s.Heading != "" {
			cells = append(cells, []interface{}{
				"",
				render.CellSpan{
					render.Centred(render.Colour(s.Heading, render.Gray)),
					len(g.Players),
				},
			})
		}
		for _, r := range s.Resources {
			row := []interface{}{render.RightAligned(RenderResourceSymbol(r))}
			for i := 0; i < pLen; i++ {
				p := g.NumFromPlayer(fromPlayer, i)
				colour := render.Gray
				bold := false
				if g.IsNeighbour(player, p) {
					colour = render.Black
				}
				if p == player {
					colour = render.Black
					bold = true
				}
				row = append(row, render.Centred(render.Markup(
					strconv.Itoa(g.PlayerResourceCount(p, r)),
					colour,
					bold,
				)))
			}
			cells = append(cells, row)
		}
	}
	return render.Table(cells, 0, 2)
}
