package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/beefsack/boredga.me/game"
	"io/ioutil"
	"os"
	"strings"
)

const FILE = ".game"

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
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	actionName := flag.Args()[0]
	action := Actions()[actionName]
	if action == nil {
		fmt.Fprintf(os.Stderr, "Invalid action: %s", actionName)
		fmt.Fprintln(os.Stderr)
		os.Exit(2)
	}
	err := action(flag.Args()[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr)
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
	err, g := newGame(players)
	if err != nil {
		return err
	}
	return saveGame(g)
}

func PlayAction(args []string) error {
	if len(args) < 2 {
		return errors.New("You must specify the player name first, followed by the plays to make")
	}
	err, g := loadGame()
	if err != nil {
		return err
	}
	err = g.PlayerAction(args[0], args[1], args[2:])
	if err != nil {
		return err
	}
	err, output := g.RenderForPlayer(args[0])
	if err != nil {
		return err
	}
	saveGame(g)
	fmt.Println(output)
	return nil
}

func saveGame(g game.Playable) error {
	data, err := g.Encode()
	if err != nil {
		return err
	}
	fi, err := os.OpenFile(FILE, os.O_WRONLY, 0666)
	if err != nil && os.IsNotExist(err) {
		fi, err = os.Create(FILE)
	}
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fi)
	_, err = writer.WriteString(g.Identifier() + "\n" + string(data))
	if err != nil {
		return err
	}
	return fi.Close()
}

func loadGame() (error, game.Playable) {
	fi, err := os.Open(FILE)
	if err != nil {
		return err, nil
	}
	reader := bufio.NewReader(fi)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err, nil
	}
	gameType := strings.Trim(line, " \n")
	g := game.RawCollection()[gameType]
	if g == nil {
		return errors.New("Could not match " + gameType + " to game type"),
			nil
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err, nil
	}
	err = g.Decode(data)
	return err, g
}
