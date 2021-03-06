package tic_tac_toe

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
)

func TestNew(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
}

func TestNewErrorsWithIncorrectPlayers(t *testing.T) {
	players := []string{"Mick"}
	game := &Game{}
	err := game.Start(players)
	if err == nil {
		t.Fail()
	}
}

func TestRenderForPlayer(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
	_, err = game.RenderForPlayer("Mick")
}

func TestPlayerAction(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
	// First lets see that a valid action works
	_, err = command.CallInCommands(game.CurrentlyMoving, game, "play a", game.Commands(game.CurrentlyMoving))
	if err != nil {
		t.Error(err)
	}
	if game.Board[0][0] == 0 {
		t.Error("The action didn't actually do anything")
	}
	// Now lets make an invalid action
	_, err = command.CallInCommands(game.CurrentlyMoving, game, "play moog", game.Commands(game.CurrentlyMoving))
	if err == nil {
		t.Error("It didn't actually error")
	}
}

func TestNextPlayer(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
	// Force the CurrentlyMoving[0] to what we want for testing sake
	game.CurrentlyMoving = "Mick"
	game.NextPlayer()
	if game.CurrentlyMoving != "Steve" {
		t.Error("Player didn't change to Steve")
	}
	game.NextPlayer()
	if game.CurrentlyMoving != "Mick" {
		t.Error("Player didn't change back to Mick")
	}
}

func TestPlaySameCell(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
	game.StartPlayer = "Mick"
	game.CurrentlyMoving = "Mick"
	// Play on a cell
	_, err = command.CallInCommands("Mick", game, "play a", game.Commands("Mick"))
	if err != nil {
		t.Error(err)
	}
	// Try to play again on the cell, it should return any errors from
	// MarkCellForPlayer
	_, err = command.CallInCommands("Steve", game, "play a", game.Commands("Steve"))
	if err == nil {
		t.Error("It allowed us to play on the same cell, it should return error if MarkCellForPlayer errors")
	}
}

func TestMarkCellForPlayer(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
	game.StartPlayer = "Mick"
	game.CurrentlyMoving = "Mick"
	// First lets just mark a cell and see if it worked
	err = game.MarkCellForPlayer("Mick", 1, 1)
	if err != nil {
		t.Error(err)
	}
	if game.Board[1][1] != 1 {
		t.Error("Didn't mark cell for Mick")
	}
	// Lets mark a different cell with the other player to see that works too
	err = game.MarkCellForPlayer("Steve", 1, 2)
	if err != nil {
		t.Error(err)
	}
	if game.Board[2][1] != 2 {
		t.Error("Didn't mark cell for Steve")
	}
	// Now lets try to remark a cell that's already marked and expect an error
	err = game.MarkCellForPlayer("Steve", 1, 1)
	if err == nil {
		t.Error("It let us change a cell that was already marked")
	}
}

func TestCheckWinner(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	game.StartPlayer = "Mick"
	game.CurrentlyMoving = "Mick"
	if err != nil {
		t.Error(err)
	}
	// This should be a winner for Mick
	game.Board[0] = [3]int{1, 1, 2}
	game.Board[1] = [3]int{1, 2, 1}
	game.Board[2] = [3]int{1, 1, 2}
	if !game.IsFinished() || game.Winner() != "Mick" {
		t.Error("Winner isn't Mick")
	}
	// This should be a winner for Steve
	game.Board[0] = [3]int{2, 1, 2}
	game.Board[1] = [3]int{1, 2, 1}
	game.Board[2] = [3]int{1, 1, 2}
	if !game.IsFinished() || game.Winner() != "Steve" {
		t.Error("Winner isn't Steve")
	}
	// This should be a winner for nobody
	game.Board[0] = [3]int{1, 1, 2}
	game.Board[1] = [3]int{2, 2, 1}
	game.Board[2] = [3]int{1, 1, 2}
	if !game.IsFinished() || len(game.Winners()) != 0 {
		t.Error("The game wasn't a draw as expected")
	}
}

func TestIsFinished(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	game.StartPlayer = "Mick"
	game.CurrentlyMoving = "Mick"
	if err != nil {
		t.Error(err)
	}
	// This shouldn't be finished
	game.Board[0] = [3]int{1, 1, 2}
	game.Board[1] = [3]int{0, 0, 1}
	game.Board[2] = [3]int{1, 1, 2}
	if game.IsFinished() {
		t.Error("Game shouldn't be finished")
	}
	// This should finished because all cells are full
	game.Board[0] = [3]int{1, 1, 2}
	game.Board[1] = [3]int{2, 2, 1}
	game.Board[2] = [3]int{1, 1, 2}
	if !game.IsFinished() {
		t.Error("Game should be finished because all cells are full")
	}
	// This should finished because Mick got a diagonal
	game.Board[0] = [3]int{0, 0, 1}
	game.Board[1] = [3]int{0, 1, 0}
	game.Board[2] = [3]int{1, 0, 0}
	if !game.IsFinished() {
		t.Error("Game should be finished because Mick won")
	}
}

func TestAllowUpperCase(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
	// Test to see if uppercase plays work
	_, err = command.CallInCommands(game.CurrentlyMoving, game, "play A", game.Commands(game.CurrentlyMoving))
	if err != nil {
		t.Error(err)
	}
	if game.Board[0][0] == 0 {
		t.Error("Using uppercase didn't mark the cell")
	}
}
