package seven_wonders

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

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

	FieldMathematics: render.Red,
	FieldEngineering: render.Green,
	FieldTheology:    render.Blue,

	AttackStrength: render.Red,
	TokenVictory:   render.Red,
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
	for i, c := range g.Hands[pNum].Sort() {
		crd := c.(Carder)
		costStrs := []string{crd.GetCard().Cost.String()}
		for _, f := range crd.GetCard().FreeWith {
			costStrs = append(costStrs, RenderCardName(Cards[f]))
		}
		afford := "  "
		switch rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 4 {
		case 0:
			afford = fmt.Sprintf("%s ", CanBuySymbol)
		case 1:
			afford = RenderMoney(1)
		case 2:
			afford = fmt.Sprintf("%s ", CannotBuySymbol)
		}
		output.WriteString(fmt.Sprintf(
			"(%s) %s %s\n          Cost: %s\n",
			render.Markup(strconv.Itoa(i+1), render.Gray, true),
			afford,
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
	return output.String(), nil
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

func RenderDirections(directions []int) string {
	dirStrs := []string{}
	for _, d := range directions {
		dirStrs = append(dirStrs, DirStrings[d])
	}
	return render.Bold(strings.Join(dirStrs, " "))
}

func (c Cost) String() string {
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
