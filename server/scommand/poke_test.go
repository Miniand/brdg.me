package scommand

import (
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/tic_tac_toe"
	"testing"
)

func TestPokeCall(t *testing.T) {
	g := &tic_tac_toe.Game{}
	err := g.Start([]string{"mick", "steve"})
	if err != nil {
		t.Fatal(err)
	}
	g.CurrentlyMoving = "mick"
	output, err := command.CallInCommands("steve", g, "   poke  ",
		[]command.Command{
			PokeCommand{},
		})
	if err != nil {
		t.Fatal(err)
	}
	if output != "You poked the current turn players" {
		t.Fatal("Expected output to be 'You poked mick' but got:", output)
	}
}
