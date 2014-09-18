package scommand

import (
	"bytes"
	"errors"
	"regexp"
	"strings"

	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/communicate"
	"github.com/Miniand/brdg.me/server/model"
)

type NewCommand struct{}

func (nc NewCommand) Parse(input string) []string {
	return regexp.MustCompile(
		`(?im)^\s*new\s+(\S+)((\s+\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}\b)+)\s*$`).
		FindStringSubmatch(input)
}

func (nc NewCommand) CanCall(player string, context interface{}) bool {
	u, err := model.FirstUserByEmail(player)
	if err != nil || u != nil && u.Unsubscribed {
		return false
	}
	return context == nil || context.(game.Playable).IsFinished()
}

func (nc NewCommand) Call(player string, context interface{},
	args []string) (string, error) {
	if len(args) < 2 {
		errors.New("Could not find game name and email addresses")
	}
	g := game.RawCollection()[args[1]]
	players := append([]string{player}, regexp.MustCompile(`\s+`).Split(
		strings.ToLower(strings.TrimSpace(args[2])), -1)...)
	gm, err := model.StartNewGame(g, players)
	if err != nil {
		return "", err
	}
	return "", communicate.Game(gm.Id, g, g.PlayerList(),
		append(g.Commands(), Commands(gm.Id)...),
		"You have been invited by "+player+" to play "+g.Name()+" by email!",
		true)
}

func (nc NewCommand) Usage(player string, context interface{}) string {
	usage := bytes.NewBufferString(
		"{{b}}new (game ID) (email addresses){{_b}} start a new game with friends\n")
	usage.WriteString("   Available games:")
	for gName, g := range game.RawCollection() {
		usage.WriteString("\n   ")
		usage.WriteString(gName)
		usage.WriteString(" (")
		usage.WriteString(g.Name())
		usage.WriteString(")")
	}
	return usage.String()
}
