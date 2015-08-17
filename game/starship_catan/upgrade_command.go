package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type UpgradeCommand struct{}

func (c UpgradeCommand) Name() string { return "upgrade" }

func (c UpgradeCommand) Call(
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
	if err != nil || len(args) != 1 {
		return "", errors.New("you must specify which module to upgrade")
	}
	m, err := ParseModule(args[0])
	if err != nil {
		return "", err
	}
	return "", g.UpgradeModule(p, m)
}

func (c UpgradeCommand) Usage(player string, context interface{}) string {
	return "{{b}}upgrade ##{{_b}} to upgrade a module.  Eg. {{b}}upgrade logistics{{_b}}"
}

func (g *Game) CanUpgrade(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseTradeAndBuild &&
		len(g.AvailableModuleUpgrades(player)) > 0
}

func (g *Game) UpgradeModule(player, module int) error {
	if !g.CanUpgrade(player) {
		return errors.New("you can't upgrade modules at the moment")
	}
	if !g.CanUpgradeModule(player, module) {
		return fmt.Errorf("you can't upgrade %s", ModuleNames[module])
	}
	newLevel := g.PlayerBoards[player].Modules[module] + 1
	t := ModuleTransaction(newLevel)
	if !g.PlayerBoards[player].CanAfford(t) {
		return t.CannotAffordError()
	}
	g.PlayerBoards[player].Transact(t)
	g.PlayerBoards[player].Modules[module] = newLevel
	g.LogTransaction(player, t)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s upgraded their {{b}}%s module{{_b}} to {{b}}level %d{{_b}}`,
		g.RenderName(player),
		ModuleNames[module],
		newLevel,
	)))
	return nil
}

func (g *Game) AvailableModuleUpgrades(player int) []int {
	modules := []int{}
	for _, m := range Modules {
		if g.CanUpgradeModule(player, m) {
			modules = append(modules, m)
		}
	}
	return modules
}

func (g *Game) CanUpgradeModule(player, module int) bool {
	opponent := (player + 1) % 2
	return g.PlayerBoards[player].Modules[module] == 0 ||
		g.PlayerBoards[player].Modules[module] == 1 &&
			g.PlayerBoards[opponent].Modules[module] < 2
}
