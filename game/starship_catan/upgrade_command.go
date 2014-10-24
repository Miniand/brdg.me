package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type UpgradeCommand struct{}

func (c UpgradeCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("upgrade", 1, -1, input)
}

func (c UpgradeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanUpgrade(p)
}

func (c UpgradeCommand) Call(player string, context interface{},
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
	return "", g.UpgradeModule(p, m)
}

func (c UpgradeCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return ""
	}
	cells := [][]string{{"{{b}}Module{{_b}}", "{{b}}Price{{_b}}"}}
	for _, m := range g.AvailableModuleUpgrades(p) {
		cells = append(cells, []string{
			ModuleNames[m],
			ModuleTransaction(g.PlayerBoards[p].Modules[m] + 1).LoseString(),
		})
	}
	table, err := render.Table(cells, 0, 2)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(
		`{{b}}upgrade ##{{_b}} to upgrade a module.  Eg. {{b}}upgrade logistics{{_b}}
%s`,
		table,
	)
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
