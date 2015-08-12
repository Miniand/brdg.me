package age_of_war

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type AttackCommand struct{}

func (c AttackCommand) Name() string { return "attack" }

func (c AttackCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a castle to attack")
	}

	castleNames := []string{}
	for _, c := range Castles {
		castleNames = append(castleNames, c.Name)
	}
	castle, err := helper.MatchStringInStrings(args[0], castleNames)
	if err != nil {
		return "", err
	}

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
	if g.Conquered[castle] && g.CastleOwners[castle] == player {
		return errors.New("you have already conquered that castle")
	}
	if ok, _ := g.ClanConquered(Castles[castle].Clan); ok {
		return errors.New("that clan is already conquered")
	}
	g.CurrentlyAttacking = castle
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s is attacking:\n%s",
		g.PlayerName(player),
		g.RenderCastle(castle, []int{}),
	)))
	g.CheckEndOfTurn()
	return nil
}
