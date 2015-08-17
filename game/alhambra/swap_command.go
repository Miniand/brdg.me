package alhambra

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type SwapCommand struct{}

func (c SwapCommand) Name() string { return "swap" }

func (c SwapCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", ErrCouldNotFindPlayer
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 2 {
		return "", errors.New("you must specify the number tile from your reserve and a coordinate")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("the first argument must be a number")
	}
	n-- // zero index

	v, err := g.Boards[pNum].Grid.ParseCoord(args[1])
	if err != nil {
		return "", err
	}

	return "", g.Swap(pNum, n, v)
}

func (c SwapCommand) Usage(player string, context interface{}) string {
	return "{{b}}swap # ##{{_b}} to swap a tile between your reserve and your Alhambra, eg. {{b}}swap 2 b4{{_b}}"
}

func (g *Game) CanSwap(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseAction &&
		len(g.Boards[player].Reserve) > 0
}

func (g *Game) Swap(player, n int, v Vect) error {
	if !g.CanSwap(player) {
		return errors.New("unable to swap a tile right now")
	}
	if n < 0 || n >= len(g.Boards[player].Reserve) {
		return errors.New(
			"not a valid tile number for your reserve tiles")
	}
	if g.Boards[player].Grid.TileAt(v).Type == TileTypeEmpty {
		return errors.New("there isn't any tile there")
	}

	rt := g.Boards[player].Reserve[n]
	t := g.Boards[player].Grid[v]
	testG := g.Boards[player].Grid.Clone()
	testG[v] = rt
	if valid, reason := testG.IsValid(); !valid {
		return errors.New(reason)
	}
	g.Boards[player].Grid = testG
	g.Boards[player].Reserve = append(
		g.Boards[player].Reserve[:n],
		g.Boards[player].Reserve[n+1:]...,
	)
	g.Boards[player].Reserve = append(g.Boards[player].Reserve, t)

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s swapped %s in their reserve for %s in their Alhambra",
		g.PlayerName(player),
		RenderTileAbbr(rt.Type),
		RenderTileAbbr(t.Type),
	)))
	g.NextPhase()
	return nil
}
