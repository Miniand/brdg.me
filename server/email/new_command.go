package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/beefsack/brdg.me/game"
	"github.com/beefsack/brdg.me/server/model"
	"regexp"
	"strings"
)

type NewCommand struct{}

func (nc NewCommand) Parse(input string) []string {
	return regexp.MustCompile(
		`(?im)^\s*new\s+(\S+)((\s+\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}\b)+)\s*$`).
		FindStringSubmatch(input)
}

func (nc NewCommand) CanCall(player string, context interface{}) bool {
	unsubscribed, err := UserIsUnsubscribed(player)
	if err == nil && unsubscribed {
		return false
	}
	return context == nil || context.(game.Playable).IsFinished()
}

func (nc NewCommand) Call(player string, context interface{}, args []string) error {
	if len(args) < 2 {
		errors.New("Could not find game name and email addresses")
	}
	gType := game.Collection()[args[1]]
	if gType == nil {
		return errors.New(fmt.Sprintf(
			`Sorry, could not find a game called "%s", please see below for available game IDs`,
			args[1]))
	}
	players := append([]string{player}, regexp.MustCompile(`\s+`).Split(
		strings.TrimSpace(args[2]), -1)...)
	g, err := gType(players)
	if err != nil {
		return err
	}
	gm, err := model.SaveGame(g)
	if err != nil {
		return err
	}
	return CommunicateGameTo(gm.Id, g, g.PlayerList(),
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
