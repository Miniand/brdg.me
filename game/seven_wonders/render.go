package seven_wonders

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/card"
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

var ResourceSymbols = map[int]string{
	GoodCoin:    "$",
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
	output := bytes.NewBuffer([]byte{})
	for _, d := range []card.Deck{
		DeckAge1(7),
		DeckAge2(7),
		DeckAge3(7),
		DeckGuild,
	} {
		for _, c := range d {
			crd := c.(Carder)
			costStrs := []string{crd.GetCard().Cost.String()}
			for _, f := range crd.GetCard().FreeWith {
				costStrs = append(costStrs, RenderCardName(Cards[f]))
			}
			output.WriteString(fmt.Sprintf(
				"%s\n    Cost: %s\n",
				RenderCard(crd),
				strings.Join(costStrs, render.Colour("  or  ", render.Gray)),
			))
			for _, f := range crd.GetCard().MakesFree {
				output.WriteString(fmt.Sprintf(
					"    Makes free: %s\n",
					RenderCard(Cards[f]),
				))
			}
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}
	return output.String(), nil
}

func RenderMoney(n int) string {
	return RenderResourceColour(
		fmt.Sprintf("%s%d", ResourceSymbols[GoodCoin], n),
		GoodCoin,
		true,
	)
}

func RenderVP(n int) string {
	return RenderResourceWithSymbol(
		fmt.Sprintf("%d ", n),
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
		return "free"
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
