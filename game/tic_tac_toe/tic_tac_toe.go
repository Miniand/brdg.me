package tic_tac_toe

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type Game struct {
	Players         []string
	CurrentlyMoving string
	StartPlayer     string
	Board           [3][3]int // 0 = empty cell, 1 = first player, 2 = second player
}

// Create a new game for specified players.  We return a pointer to make sure it
// confirms to interfaces.
func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Must be 2 players")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	startPlayer := players[r.Int()%2]
	g.Players = players
	g.StartPlayer = startPlayer
	g.CurrentlyMoving = startPlayer
	return nil
}

// Make an action for the specified player
func (g *Game) PlayerAction(player string, action string, args []string) error {
	if g.CurrentlyMoving != player {
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
	if g.CurrentlyMoving == g.Players[0] {
		g.CurrentlyMoving = g.Players[1]
	} else {
		g.CurrentlyMoving = g.Players[0]
	}

}

// Marks the specified cell for the current player and changes the currently
// moving player to the next one.  It shouldn't let you mark a cell that's
// already marked.
func (g *Game) MarkCellForPlayer(player string, x, y int) error {
	if g.Board[y][x] != 0 {
		return errors.New("cell not empty, bro")
	} else {
		if g.CurrentlyMoving == g.StartPlayer {
			g.Board[y][x] = 1
		} else {
			g.Board[y][x] = 2
		}
	}
	return nil
}

// Render an ascii representation of the game for a player
func (g *Game) RenderForPlayer(player string) (error, string) {
	output := fmt.Sprintf("%#v", g)
	// @todo implement.
	//return errors.New("Not implemented yet"), output
	return nil, output

}

// Gets a full list of players.
func (g *Game) PlayerList() []string {
	return g.Players
}

// Check if there is a winner, if there is a line of 3 all 1s or 2s
func (g *Game) Winner() string {
	for i := 0; i < 3; i++ {
		if g.Board[i][0] == g.Board[i][1] && g.Board[i][0] == g.Board[i][2] {
			return g.Players[g.Board[i][0]]
		} else if g.Board[0][i] == g.Board[1][i] && g.Board[0][i] == g.Board[2][i] {
			return g.Players[g.Board[i][0]]
		}
	}
	return ""
}

// Wrapper of Winner to match game interface (some games can have more than 1 winner)
func (g *Game) Winners() []string {
	winner := g.Winner()
	if winner != "" {
		return []string{winner}
	}
	return []string{}
}

// Check if the game is finished, i.e. if there is a winner or if there is no
// empty cells
func (g *Game) IsFinished() bool {
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if g.Board[x][y] == 0 {
				return false
			}
		}
	}
	return true
}

// Returns all the users whose turn it is.
func (g *Game) WhoseTurn() []string {
	return []string{g.CurrentlyMoving}
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
