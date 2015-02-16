package seven_wonders

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/cost"
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

	WonderStage: render.Yellow,

	VP: render.Green,
}

const (
	CardSymbol = "##"
)

var CanBuySymbol = render.Markup("✔", render.Green, true)
var CanBuyWithCoinSymbol = render.Markup(
	ResourceSymbols[GoodCoin],
	render.Yellow,
	true,
)
var CannotBuySymbol = render.Markup("✘", render.Red, true)

var ResourceSymbols = map[int]string{
	GoodCoin:    "●",
	GoodWood:    "Wo",
	GoodStone:   "St",
	GoodOre:     "Or",
	GoodClay:    "Cl",
	GoodPapyrus: "Pa",
	GoodTextile: "Te",
	GoodGlass:   "Gl",

	CardKindRaw:          CardSymbol,
	CardKindManufactured: CardSymbol,
	CardKindCivilian:     CardSymbol,
	CardKindScientific:   CardSymbol,
	CardKindCommercial:   CardSymbol,
	CardKindMilitary:     CardSymbol,
	CardKindGuild:        CardSymbol,
	CardKindWonder:       "▲",

	FieldMathematics: "Ma",
	FieldEngineering: "En",
	FieldTheology:    "Th",

	AttackStrength: "X",
	TokenVictory:   "V",
	TokenDefeat:    "X",

	WonderStage: "▲",

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
		fmt.Sprintf("%s%s", text, ResourceSymbols[resource]),
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
	output := bytes.NewBuffer([]byte{})

	// Action output
	if g.Actions[pNum] != nil {
		actOut := g.Actions[pNum].Output(pNum, g)
		if actOut != "" {
			output.WriteString(actOut)
			output.WriteString("\n\n")
		}
	}

	// Resolver output
	if len(g.ToResolve) > 0 {
		resOut := g.ToResolve[0].String(pNum, g)
		if resOut != "" {
			output.WriteString(resOut)
			output.WriteString("\n\n")
		}
	}

	// Hand
	output.WriteString(render.Bold("Your hand:\n\n"))
	output.WriteString(
		g.RenderCardList(pNum, g.Hands[pNum], g.CanAction(pNum), true))

	// Wonders
	output.WriteString(render.Bold("\n\nRemaining wonder stages:\n\n"))
	remaining := g.RemainingWonderStages(pNum)
	if len(remaining) > 0 {
		output.WriteString(
			g.RenderCardList(pNum, remaining, false, true))
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
	return output.String(), nil
}

func (g *Game) RenderCardList(
	player int,
	cards card.Deck,
	showNums, showAfford bool,
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
		affordStr := ""
		if showAfford {
			canBuild, coins := g.CanBuildCard(player, crd)
			afford := ""
			switch {
			case !canBuild:
				afford = fmt.Sprintf("%s ", CannotBuySymbol)
			case canBuild && len(coins) == 0:
				afford = fmt.Sprintf("%s ", CanBuySymbol)
			default:
				sum := 0
				for _, coin := range coins[0] {
					sum += coin
				}
				afford = RenderMoney(sum)
			}
			affordStr = fmt.Sprintf(
				"%s ",
				afford,
			)
		}
		output.WriteString(fmt.Sprintf(
			"%s%s%s\n          Cost: %s\n",
			numStr,
			affordStr,
			RenderCard(crd),
			strings.Join(costStrs, render.Colour("  or  ", render.Gray)),
		))
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
			switch n {
			case GoodCoin:
				parts = append(parts, RenderMoney(amount))
			default:
				for i := 0; i < amount; i++ {
					parts = append(parts, RenderResourceSymbol(n))
				}
			}
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
			[]int{GoodCoin, VP, WonderStage},
		},
		{
			"Goods",
			Goods,
		},
		{
			"Cards",
			CardKinds,
		},
		{
			"Military",
			[]int{AttackStrength, TokenVictory, TokenDefeat},
		},
		{
			"Science",
			Fields,
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
				render.CellSpan{
					render.Centred(render.Markup(s.Heading, render.Gray, true)),
					len(g.Players) + 1,
				},
			})
		}
		for _, r := range s.Resources {
			row := []interface{}{RenderResourceSymbol(r)}
			for i := 0; i < pLen; i++ {
				p := g.NumFromPlayer(fromPlayer, i)
				colour := render.Gray
				bold := false
				if g.IsNeighbour(player, p) {
					colour = render.Black
				}
				if p == player {
					colour = ResourceColours[r]
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
