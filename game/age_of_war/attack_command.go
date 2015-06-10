package age_of_war

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type AttackCommand struct{}

func (c AttackCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("attack", 1, input)
}

func (c AttackCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanAttack(pNum)
}

func (c AttackCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) != 1 {
		return "", errors.New("you must specify a castle to attack")
	}

	castle := 0
	return "", g.Attack(pNum, castle)
}

func (c AttackCommand) Usage(player string, context interface{}) string {
	return "{{b}}attack #{{_b}} to attack a castle, eg. {{b}}attack kita{{_b}}"
}

func (g *Game) CanAttack(player int) bool {
	return g.CurrentPlayer == player && g.CurrentlyAttacking == -1
}

func (g *Game) Attack(player, castle int) error {
	if !g.CanAttack(player) {
		return errors.New("unable to attack a castle right now")
	}
	if castle < 0 || castle >= len(Castles) {
		return errors.New("that is not a valid castle")
	}
	return errors.New("not implemented")
}
