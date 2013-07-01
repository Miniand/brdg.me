package main

import (
	// "encoding/gob"
	"errors"
	"flag"
	"fmt"
	"github.com/beefsack/boredga.me/game"
	"os"
)

func Actions() map[string](func([]string) error) {
	return map[string](func([]string) error){
		"new":  NewAction,
		"play": PlayAction,
	}
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Put an action")
		fmt.Println()
		os.Exit(1)
	}
	actionName := flag.Args()[0]
	action := Actions()[actionName]
	if action == nil {
		fmt.Fprintf(os.Stderr, "Invalid action: %s", actionName)
		fmt.Println()
		os.Exit(2)
	}
	err := action(flag.Args()[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		fmt.Println()
		os.Exit(3)
	}
}

func NewAction(args []string) error {
	if len(args) < 2 {
		return errors.New("You must specify a game name (use list to list games) and at least one player")
	}
	gameName := args[0]
	players := args[1:]
	newGame := game.Collection()[gameName]
	if newGame == nil {
		return errors.New("Could not find game " + gameName)
	}
	err, game := newGame(players)
	if err != nil {
		return err
	}
	fmt.Println(game)
	return nil
}

func PlayAction(args []string) error {
	if len(args) < 2 {
		return errors.New("You must specify the player name first, followed by the plays to make")
	}
	return nil
}
