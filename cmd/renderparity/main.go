// Render-parity harness: starts an old brdg.me game, renders plain text for
// every player, optionally applies text commands, and persists state to a
// file so hidden information (hands) can be inspected between commands.
//
// Usage:
//   renderparity new <statefile> <game> <player>...
//   renderparity play <statefile> <player> <command...>
//   renderparity render <statefile>
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/acquire"
	"github.com/Miniand/brdg.me/game/jaipur"
	"github.com/Miniand/brdg.me/game/lost_cities"
	"github.com/Miniand/brdg.me/game/tic_tac_toe"
	"github.com/Miniand/brdg.me/render"
)

func games() map[string]game.Playable {
	return map[string]game.Playable{
		"tic_tac_toe": &tic_tac_toe.Game{},
		"lost_cities": &lost_cities.Game{},
		"acquire":     &acquire.Game{},
		"jaipur":      &jaipur.Game{},
	}
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func save(path string, g game.Playable) {
	data, err := g.Encode()
	die(err)
	die(os.WriteFile(path, append([]byte(g.Identifier()+"\n"), data...), 0644))
}

func load(path string) game.Playable {
	raw, err := os.ReadFile(path)
	die(err)
	i := strings.IndexByte(string(raw), '\n')
	if i < 0 {
		die(fmt.Errorf("bad state file %s", path))
	}
	g, ok := games()[string(raw[:i])]
	if !ok {
		die(fmt.Errorf("unknown game %q", string(raw[:i])))
	}
	die(g.Decode(raw[i+1:]))
	return g
}

func renderAll(g game.Playable) {
	for _, p := range g.PlayerList() {
		fmt.Printf("=== render for %s ===\n", p)
		out, err := g.RenderForPlayer(p)
		die(err)
		fmt.Println(render.RenderPlain(out))
	}
	fmt.Printf("=== whose turn: %s | finished: %v",
		strings.Join(g.WhoseTurn(), ", "), g.IsFinished())
	if g.IsFinished() {
		fmt.Printf(" | winners: %s", strings.Join(g.Winners(), ", "))
	}
	fmt.Println(" ===")
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: renderparity new <statefile> <game> <player>...\n"+
			"       renderparity play <statefile> <player> <command...>\n"+
			"       renderparity render <statefile>")
		os.Exit(2)
	}
	statefile := os.Args[2]
	switch os.Args[1] {
	case "new":
		g, ok := games()[os.Args[3]]
		if !ok {
			die(fmt.Errorf("unknown game %q", os.Args[3]))
		}
		die(g.Start(os.Args[4:]))
		save(statefile, g)
		renderAll(g)
	case "play":
		g := load(statefile)
		player := os.Args[3]
		input := strings.Join(os.Args[4:], " ")
		output, err := command.CallInCommands(player, g, input, g.Commands(player))
		die(err)
		if output != "" {
			fmt.Println("--- command output ---")
			fmt.Println(render.RenderPlain(output))
		}
		save(statefile, g)
		renderAll(g)
	case "render":
		renderAll(load(statefile))
	default:
		die(fmt.Errorf("unknown subcommand %q", os.Args[1]))
	}
}
