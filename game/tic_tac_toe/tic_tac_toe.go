package tic_tac_toe

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
)

type Game struct {
	Players         []string
	CurrentlyMoving []string
	StartPlayer     string
	Board           [3][3]int // 0 = empty cell, 1 = first player, 2 = second player
}

// Create a new game for specified players.  We return a pointer to make sure it
// confirms to interfaces.
func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Must be 2 players")
	}
	startPlayer := players[rand.Int()%2]
	g.Players = players
	g.StartPlayer = startPlayer
	g.CurrentlyMoving = []string{startPlayer}
	return nil
}

// Make an action for the specified player
func (g *Game) PlayerAction(player string, action string, args []string) error {
	if g.CurrentlyMoving[0] != player {
		return errors.New("Not your turn")
	}
	if !regexp.MustCompile("^[abcdefghi]$").MatchString(action) {
		return errors.New("Your action must be a letter between a - i")
	}
	switch action {
	case "a":
		g.MarkCellForPlayer(player, 0, 0)
	case "b":
		g.MarkCellForPlayer(player, 1, 0)
	case "c":
		g.MarkCellForPlayer(player, 2, 0)
	case "d":
		g.MarkCellForPlayer(player, 0, 1)
	case "e":
		g.MarkCellForPlayer(player, 1, 1)
	case "f":
		g.MarkCellForPlayer(player, 2, 1)
	case "g":
		g.MarkCellForPlayer(player, 0, 2)
	case "h":
		g.MarkCellForPlayer(player, 1, 2)
	case "i":
		g.MarkCellForPlayer(player, 2, 2)
	default:
		fmt.Println(action, "how did this get here...")
	}
	// @todo Convert a-i to cell coordinates (both 0-2) and call
	// MarkCellForPlayer
	g.NextPlayer()
	return nil
}

func (g *Game) NextPlayer() {
	// @todo Flip g.CurrentlyMoving[0] to the other player
}

// Marks the specified cell for the current player and changes the currently
// moving player to the next one.  It shouldn't let you mark a cell that's
// already marked.
func (g *Game) MarkCellForPlayer(player string, x, y int) error {
	// @todo implement
	return errors.New("Not implemented yet")
}

// Render an ascii representation of the game for a player
func (g *Game) RenderForPlayer(player string) (error, string) {
	output := "This is an example\n"
	output += "of some constructed output"
	// @todo implement.
	return errors.New("Not implemented yet"), output
}

// Check if there is a winner, if there is a line of 3 all 1s or 2s.  First
// argument is false if there isn't a winner yet
func (g *Game) CheckWinner() (bool, string) {
	// @todo implement
	return false, ""
}

// Check if the game is finished, i.e. if there is a winner or if there is no
// empty cells
func (g *Game) IsFinished() bool {
	won, _ := g.CheckWinner()
	if won {
		return true
	}
	// @todo check if there are any empty cells and return false if there are
	// any, otherwise return true
	return false
}

// We use human name for output.
func (g *Game) Name() string {
	return "Tic-tac-toe"
}

// We use machine name for referencing.
func (g *Game) Identifier() string {
	return "tic_tac_toe"
}

// Encode to a string
func (g *Game) Encode() ([]byte, error) {
	return json.Marshal(g)
}

// Decode from a string
func (g *Game) Decode(data []byte) error {
	return json.Unmarshal(data, g)
}
