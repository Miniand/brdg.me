package king_of_tokyo

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

const (
	DefenderActionNone = iota
	DefenderActionLeaveTokyo
)

type StayCommand struct{}

func (c StayCommand) Parse(input string) []string {
	return command.ParseNamedCommand("stay", input)
}

func (c StayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanStay(pNum)
}

func (c StayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	return "", g.Stay(pNum)
}

func (c StayCommand) Usage(player string, context interface{}) string {
	return "{{b}}stay{{_b}} to stay in Tokyo"
}

func (g *Game) CanStay(player int) bool {
	return g.Phase == PhaseAttack && g.AttackPlayers[0] == player
}

func (g *Game) Stay(player int) error {
	if !g.CanStay(player) {
		return errors.New("you can't call stay at the moment")
	}
	g.PostStayOrLeave(DefenderActionNone)
	return nil
}

func (g *Game) PostStayOrLeave(defenderAction int) {
	p := g.AttackPlayers[0]
	damage := g.AttackDamage
	g.DealDamage(g.CurrentPlayer, p, damage, defenderAction)
	g.NextAttackedPlayer()
}
