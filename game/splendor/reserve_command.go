package splendor

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type ReserveCommand struct{}

func (c ReserveCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("reserve", 1, input)
}

func (c ReserveCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	return err == nil && g.CanReserve(pNum)
}

func (c ReserveCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify a card")
	}
	row, col, err := ParseLoc(a[0])
	if err != nil {
		return "", err
	}
	return "", g.Reserve(pNum, row, col)
}

func (c ReserveCommand) Usage(player string, context interface{}) string {
	return "{{b}}reserve #{{_b}} to reserve a card for yourself, eg. {{b}}reserve 2B{{_b}}"
}

func (g *Game) CanReserve(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseMain &&
		len(g.PlayerBoards[player].Reserve) < 3
}

func (g *Game) Reserve(player, row, col int) error {
	if !g.CanReserve(player) {
		return errors.New("unable to reserve right now")
	}
	if row < 0 || row > 2 {
		return errors.New("that is not a valid row")
	}
	if col < 0 || col >= len(g.Board[row]) {
		return errors.New("that is not a valid card")
	}
	g.PlayerBoards[player].Reserve = append(
		g.PlayerBoards[player].Reserve,
		g.Board[row][col],
	)
	if g.Tokens[Gold] > 0 {
		g.PlayerBoards[player].Tokens[Gold] += 1
		g.Tokens[Gold] -= 1
	}
	if len(g.Decks[row]) > 0 {
		g.Board[row][col] = g.Decks[row][0]
		g.Decks[row] = g.Decks[row][1:]
	} else {
		g.Board[row] = append(
			g.Board[row][:col],
			g.Board[row][col+1:]...,
		)
	}
	g.NextPhase()
	return nil
}
