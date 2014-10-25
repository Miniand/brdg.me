package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type LoseCommand struct{}

func (c LoseCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("lose", 1, input)
}

func (c LoseCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanLoseModule(p)
}

func (c LoseCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}

	if len(args) < 1 {
		return "", errors.New("you must specify which module to lose")
	}
	a := command.ExtractNamedCommandArgs(args)
	m, err := ParseModule(a[0])
	if err != nil {
		return "", err
	}

	return "", g.LoseModule(p, m)
}

func (c LoseCommand) Usage(player string, context interface{}) string {
	return "{{b}}lose ##{{_b}} to choose which module was destroyed, eg. {{b}}lose sensor{{_b}}"
}

func (g *Game) CanLoseModule(player int) bool {
	return g.CurrentPlayer == player || g.LosingModule
}

func (g *Game) LoseModule(player, module int) error {
	if !g.CanLoseModule(player) {
		return errors.New("you can't lose a module at the moment")
	}
	if g.PlayerBoards[player].Modules[module] <= 0 {
		return errors.New("you don't have that module")
	}
	g.PlayerBoards[player].Modules[module] -= 1
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s had their {{b}}%s module{{_b}} destroyed by the pirate`,
		g.RenderName(player),
		ModuleNames[module],
	)))
	g.LosingModule = false
	return g.EndFlight()
}
