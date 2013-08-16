package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/render"
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

func RenderForPlayer(g game.Playable, p string) (string, error) {
	rawOutput, err := g.RenderForPlayer(p)
	if err != nil {
		return "", err
	}
	return render.RenderTerminal(rawOutput)
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
	g, err := newGame(players)
	if err != nil {
		return err
	}
	for _, p := range g.WhoseTurn() {
		output, err := RenderForPlayer(g, p)
		if err != nil {
			return err
		}
		fmt.Println("--- OUTPUT FOR " + p + " ---")
		fmt.Println(output)
	}
	fmt.Println("Current turn: " + strings.Join(g.WhoseTurn(), ", "))
	return saveGame(g)
}

func PlayAction(args []string) error {
	if len(args) < 2 {
		return errors.New("You must specify the player name first, followed by the plays to make")
	}
	g, err := loadGame()
	if err != nil {
		return err
	}
	commandOutput, err := command.CallInCommands(args[0],
		g, strings.Join(args[1:], " "), g.Commands())
	if err != nil {
		return err
	}
	output, err := RenderForPlayer(g, args[0])
	if err != nil {
		return err
	}
	if commandOutput != "" {
		output = commandOutput + "\n\n" + output
	}
	saveGame(g)
	fmt.Println("--- OUTPUT FOR " + args[0] + " ---")
	fmt.Println(output)
	usages := command.CommandUsages(args[0], g, g.Commands())
	if len(usages) > 0 {
		commandsOutput, err := render.RenderTerminal(render.CommandUsages(
			usages))
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println(commandsOutput)
	}
	for _, p := range g.WhoseTurn() {
		output, err = RenderForPlayer(g, p)
		if err != nil {
			return err
		}
		fmt.Println("--- OUTPUT FOR " + p + " ---")
		fmt.Println(output)
		usages = command.CommandUsages(p, g,
			command.AvailableCommands(p, g, g.Commands()))
		if len(usages) > 0 {
			commandsOutput, err := render.RenderTerminal(render.CommandUsages(
				usages))
			if err != nil {
				return err
			}
			fmt.Println()
			fmt.Println(commandsOutput)
		}
	}
	// Save again in case logs were marked as read
	saveGame(g)
	if g.IsFinished() {
		fmt.Println("Game finished!  Winners: " + strings.Join(g.Winners(), ", "))
	} else {
		fmt.Println("Current turn: " + strings.Join(g.WhoseTurn(), ", "))
	}
	return nil
}

func ViewAction(args []string) error {
	if len(args) < 1 {
		return errors.New("You must specify a player to view for.")
	}
	g, err := loadGame()
	if err != nil {
		return err
	}
	output, err := RenderForPlayer(g, args[0])
	if err == nil {
		fmt.Println(output)
	}
	usages := command.CommandUsages(args[0], g, g.Commands())
	if len(usages) > 0 {
		commandsOutput, err := render.RenderTerminal(render.CommandUsages(
			usages))
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println(commandsOutput)
	}
	return err
}

func DumpAction(args []string) error {
	g, err := loadGame()
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
	return ioutil.WriteFile(FILE, []byte(g.Identifier()+"\n"+string(data)),
		0666)
}

func loadGame() (game.Playable, error) {
	fi, err := os.Open(FILE)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(fi)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	gameType := strings.Trim(line, " \n")
	g := game.RawCollection()[gameType]
	if g == nil {
		return nil, errors.New("Could not match " + gameType + " to game type")
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = g.Decode(data)
	return g, err
}
