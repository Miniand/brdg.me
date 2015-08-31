package alhambra

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type PlaceCommand struct{}

func (c PlaceCommand) Name() string { return "place" }

func (c PlaceCommand) Call(
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
		return "", errors.New("you must specify a tile and a coordinate")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("the first argument must the the tile number to place")
	}
	n-- // zero index

	v, err := g.Boards[pNum].Grid.ParseCoord(args[1])
	if err != nil {
		return "", err
	}

	return "", g.Place(pNum, n, v)
}

func (c PlaceCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	fromText := ""
	if g.Phase == PhaseAction {
		fromText = " from your reserved tiles"
	}
	return fmt.Sprintf(
		"{{b}}place # ##{{_b}} to place a tile%s, eg. {{b}}place 1 b3{{_b}}",
		fromText,
	)
}

func (g *Game) CanPlace(player int) bool {
	return g.CurrentPlayer == player &&
		(g.Phase == PhaseAction && len(g.Boards[player].Reserve) > 0 ||
			(g.Phase == PhasePlace || g.Phase == PhaseFinalPlace) &&
				len(g.Boards[player].Place) > 0)
}

func (g *Game) Place(player, n int, v Vect) error {
	if !g.CanPlace(player) {
		return errors.New("unable to place a tile right now")
	}
	if g.Boards[player].Grid.TileAt(v).Type != TileTypeEmpty {
		return errors.New("there is already a tile there")
	}

	t := Tile{}
	switch g.Phase {
	case PhaseAction:
		if n < 0 || n >= len(g.Boards[player].Reserve) {
			return errors.New(
				"this is not a valid tile number for your reserve tiles")
		}
		t = g.Boards[player].Reserve[n]
	default:
		if n < 0 || n >= len(g.Boards[player].Place) {
			return errors.New(
				"this is not a valid tile number for your placeable tiles")
		}
		t = g.Boards[player].Place[n]
	}

	testG := g.Boards[player].Grid.Clone()
	testG[v] = t
	if valid, reason := testG.IsValid(); !valid {
		return errors.New(reason)
	}
	g.Boards[player].Grid = testG
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s placed %s in their Alhambra",
		g.PlayerName(player),
		RenderTileAbbr(t.Type),
	)))

	switch g.Phase {
	case PhaseAction:
		g.Boards[player].Reserve = append(
			g.Boards[player].Reserve[:n],
			g.Boards[player].Reserve[n+1:]...,
		)
		g.NextPhase()
	default:
		g.Boards[player].Place[n] = Tile{}
		if len(NotEmpty(g.Boards[player].Place)) == 0 {
			g.NextPhase()
		}
	}
	return nil
}
