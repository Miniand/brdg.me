package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type LoseCommand struct{}

func (c LoseCommand) Name() string { return "lose" }

func (c LoseCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("you must specify which module to lose")
	}
	m, err := ParseModule(args[0])
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
