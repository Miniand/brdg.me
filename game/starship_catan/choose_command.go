package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type ChooseCommand struct{}

func (c ChooseCommand) Name() string { return "choose" }

func (c ChooseCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("please specify a module")
	}
	m, err := ParseModule(args[0])
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
