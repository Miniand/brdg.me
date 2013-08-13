package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/beefsack/brdg.me/command"
	"github.com/beefsack/brdg.me/game"
	"github.com/beefsack/brdg.me/server/model"
)

type NewCommand struct{}

func (nc NewCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("new", 1, -1, input)
}

func (nc NewCommand) CanCall(player string, context interface{}) bool {
	unsubscribed, err := UserIsUnsubscribed(player)
	if err == nil && unsubscribed {
		return false
	}
	return context == nil || context.(game.Playable).IsFinished()
}

func (nc NewCommand) Call(player string, context interface{}, args []string) error {
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return errors.New("Please also specify a game ID when starting a new game")
	}
	gType := game.Collection()[a[0]]
	if gType == nil {
		return errors.New(fmt.Sprintf(
			`Sorry, could not find a game called "%s", please see below for available game IDs`,
			a[0]))
	}
	players := append([]string{player}, a[1:]...)
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
