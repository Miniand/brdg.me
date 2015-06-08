package alhambra

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type RemoveCommand struct{}

func (c RemoveCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("remove", 1, input)
}

func (c RemoveCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanRemove(pNum)
}

func (c RemoveCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", ErrCouldNotFindPlayer
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) != 1 {
		return "", errors.New("you must specify a coordinate")
	}

	v, err := g.Boards[pNum].Grid.ParseCoord(a[0])
	if err != nil {
		return "", err
	}

	return "", g.Remove(pNum, v)
}

func (c RemoveCommand) Usage(player string, context interface{}) string {
	return "{{b}}remove ##{{_b}} to remove a tile from your Alhambra, eg. {{b}}remove b4{{_b}}"
}

func (g *Game) CanRemove(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseAction
}

func (g *Game) Remove(player int, v Vect) error {
	if !g.CanRemove(player) {
		return errors.New("unable to remove a tile right now")
	}
	if g.Boards[player].Grid.TileAt(v).Type == TileTypeEmpty {
		return errors.New("there isn't any tile there")
	}

	t := g.Boards[player].Grid[v]
	testG := g.Boards[player].Grid.Clone()
	delete(testG, v)
	if valid, reason := testG.IsValid(); !valid {
		return errors.New(reason)
	}
	g.Boards[player].Grid = testG

	g.Boards[player].Reserve = append(g.Boards[player].Reserve, t)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s removed %s from their Alhambra and added it to their reserve",
		g.PlayerName(player),
		RenderTileAbbr(t.Type),
	)))
	g.NextPhase()
	return nil
}
