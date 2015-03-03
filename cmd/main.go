package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/davecgh/go-spew/spew"
)

const FILE = ".game"

var renderer = render.RenderTerminal

func Actions() map[string](func([]string) error) {
	return map[string](func([]string) error){
		"new":  NewAction,
		"play": PlayAction,
		"bot":  BotAction,
		"view": ViewAction,
		"dump": DumpAction,
	}
}

func main() {
	var (
		html, image, profile bool
		f                    *os.File
		err                  error
	)
	flag.BoolVar(&html, "html", false, "output html")
	flag.BoolVar(&image, "image", false, "output image (PNG)")
	flag.BoolVar(&profile, "profile", false, "run profiler")
	flag.Parse()
	if html {
		renderer = render.RenderHtml
	}
	if image {
		renderer = render.RenderImage
	}
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
	if profile {
		f, err = os.OpenFile("profile", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
		fmt.Println("starting profile")
		pprof.StartCPUProfile(f)
	}
	err = action(flag.Args()[1:])
	if profile {
		pprof.StopCPUProfile()
		f.Close()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

func RenderForPlayer(g game.Playable, p string) (string, error) {
	logOutput := "\n\n{{b}}Since last time:{{_b}}:\n" +
		log.RenderMessages(g.GameLog().NewMessagesFor(p))
	g.GameLog().MarkReadFor(p)
	rawOutput, err := g.RenderForPlayer(p)
	if err != nil {
		return "", err
	}
	commandsOutput := ""
	usages := command.CommandUsages(p, g, command.AvailableCommands(
		p, g, g.Commands()))
	if len(usages) > 0 {
		commandsOutput = fmt.Sprintf("\n\n%s", strings.Join(usages, "\n"))
	}
	return renderer(logOutput + "\n\n" + rawOutput + commandsOutput)
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
	if err := saveGame(g); err != nil {
		return err
	}
	fmt.Println("--- OUTPUT FOR " + args[0] + " ---")
	fmt.Println(output)
	err = OutputGameForPlayingPlayers(g)
	if err != nil {
		return err
	}
	return nil
}

func OutputGameForPlayingPlayers(g game.Playable) error {
	for _, p := range g.WhoseTurn() {
		output, err := RenderForPlayer(g, p)
		if err != nil {
			return err
		}
		fmt.Println("--- OUTPUT FOR " + p + " ---")
		fmt.Println(output)
	}
	// Save again in case logs were marked as read
	if err := saveGame(g); err != nil {
		return err
	}
	if g.IsFinished() {
		fmt.Println("Game finished!  Winners: " + strings.Join(g.Winners(), ", "))
	} else {
		fmt.Println("Current turn: " + strings.Join(g.WhoseTurn(), ", "))
	}
	return nil
}

func BotAction(args []string) error {
	if len(args) != 1 {
		return errors.New("Specify a player name to simulate")
	}
	player := args[0]
	rawG, err := loadGame()
	if err != nil {
		return err
	}
	g, ok := rawG.(game.Botter)
	if !ok {
		return errors.New("Game does not have bot support")
	}
	err = g.BotPlay(player)
	if err != nil {
		return err
	}
	if err := saveGame(rawG); err != nil {
		return err
	}
	return OutputGameForPlayingPlayers(rawG)
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
	return err
}

func DumpAction(args []string) error {
	g, err := loadGame()
	if err != nil {
		return err
	}
	spew.Dump(g)
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
