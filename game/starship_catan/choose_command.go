package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type ChooseCommand struct{}

func (c ChooseCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("choose", 1, input)
}

func (c ChooseCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanChoose(p)
}

func (c ChooseCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	m, err := ParseModule(a[0])
	if err != nil {
		return "", err
	}
	return "", g.Choose(p, m)
}

func (c ChooseCommand) Usage(player string, context interface{}) string {
	return "{{b}}choose ##{{_b}} to choose which module to start with.  The logistics module is the most useful module for starting players.  Eg. {{b}}choose lo{{_b}}"
}

func (g *Game) CanChoose(player int) bool {
	return g.Phase == PhaseChooseModule &&
		len(g.PlayerBoards[player].Modules) == 0
}

func (g *Game) Choose(player, module int) error {
	if !g.CanChoose(player) {
		return errors.New("you can't choose a module at the moment")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s chose the {{b}}%s module{{_b}}`,
		g.RenderName(player), ModuleNames[module])))
	g.PlayerBoards[player].Modules[module] = 1
	if len(g.WhoseTurn()) == 0 {
		g.NewTurn()
	}
	return nil
}
