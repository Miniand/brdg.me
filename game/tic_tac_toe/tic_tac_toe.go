package tic_tac_toe

import (
	"errors"
	"math/rand"
	"regexp"
)

type Game struct {
	Players         []string
	CurrentlyMoving []string
	StartPlayer     string
	Board           [3][3]int // 0 = empty cell, 1 = first player, 2 = second player
}

// Create a new game for specified players
func NewGame(players []string) (error, Game) {
	if len(players) != 2 {
		return errors.New("Must be 2 players"), Game{}
	}
	startPlayer := players[rand.Int()%2]
	return nil, Game{
		Players:         players,
		StartPlayer:     startPlayer,
		CurrentlyMoving: []string{startPlayer},
	}
}

// Make an action for the specified player
func (g *Game) PlayerAction(player string, action string, args []string) error {
	if g.CurrentlyMoving[0] != player {
		return errors.New("Not your turn")
	}
	if !regexp.MustCompile("^[abcdefghi]$").MatchString(action) {
		return errors.New("Your action must be a letter between a - i")
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
