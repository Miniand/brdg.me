package roll_through_the_ages

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type InvadeCommand struct{}

func (c InvadeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("invade", 1, input)
}

func (c InvadeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanInvade(pNum)
}

func (c InvadeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify how many spearheads to use")
	}
	amount, err := strconv.Atoi(a[0])
	if err != nil || amount < 1 {
		return "", errors.New("the amount must be a positive number")
	}

	return "", g.Invade(pNum, amount)
}

func (c InvadeCommand) Usage(player string, context interface{}) string {
	return "{{b}}invade #{{_b}} to use spearheads to inflict extra damage on other players, eg. {{b}}invade 2{{_b}}"
}

func (g *Game) CanInvade(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseInvade &&
		g.Boards[player].Developments[DevelopmentSmithing] &&
		g.Boards[player].Goods[GoodSpearhead] > 0
}

func (g *Game) Invade(player, amount int) error {
	if !g.CanInvade(player) {
		return errors.New("you can't invade at the moment")
	}
	if amount <= 0 {
		return errors.New("you must specify a positive amount of spearheads")
	}
	sh := g.Boards[player].Goods[GoodSpearhead]
	if amount > sh {
		return fmt.Errorf("you only have %d spearheads", sh)
	}

	g.Boards[player].Goods[GoodSpearhead] -= amount
	buf := bytes.NewBufferString(fmt.Sprintf(
		`%s used {{b}}%d{{_b}} spearheads to cause extra damage`,
		g.RenderName(player),
		amount,
	))
	for p, _ := range g.Players {
		if p == player {
			continue
		}
		if g.Boards[p].HasBuilt(MonumentGreatWall) {
			buf.WriteString(fmt.Sprintf(
				"\n  %s avoids the extra damage with their wall",
				g.RenderName(p),
			))
		} else {
			g.Boards[p].Disasters += amount
			buf.WriteString(fmt.Sprintf(
				"\n  %s takes {{b}}%d disaster points{{_b}}",
				g.RenderName(p),
				amount,
			))
		}
	}
	g.Log.Add(log.NewPublicMessage(buf.String()))

	g.BuildPhase()
	return nil
}
