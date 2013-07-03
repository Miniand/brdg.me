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
		"view": ViewAction,
		"dump": DumpAction,
	}
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Available actions are:")
		for aName, _ := range Actions() {
			fmt.Println(aName)
		}
		os.Exit(1)
	}
	actionName := flag.Args()[0]
	action := Actions()[actionName]
	if action == nil {
		fmt.Printf("Invalid action: %s", actionName)
		fmt.Println()
		os.Exit(2)
	}
	err := action(flag.Args()[1:])
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Println()
		os.Exit(3)
	}
}

func NewAction(args []string) error {
	if len(args) < 2 {
		lines := []string{
			"You must specify a game name and at least one player.  Available games are:",
		}
		for _, rawG := range game.RawCollection() {
			lines = append(lines, rawG.Identifier()+" ("+rawG.Name()+")")
		}
		return errors.New(strings.Join(lines, "\n"))
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

func ViewAction(args []string) error {
	if len(args) < 1 {
		return errors.New("You must specify a player to view for.")
	}
	err, g := loadGame()
	if err != nil {
		return err
	}
	err, output := g.RenderForPlayer(args[0])
	if err == nil {
		fmt.Println(output)
	}
	return err
}

func DumpAction(args []string) error {
	err, g := loadGame()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", g)
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
	writer.Flush()
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
