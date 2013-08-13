package main

import (
	"errors"
	"github.com/beefsack/brdg.me/command"
	"github.com/beefsack/brdg.me/server/model"
)

type UnsubscribeCommand struct{}

func (uc UnsubscribeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("unsubscribe", 0, input)
}

func (uc UnsubscribeCommand) CanCall(player string, context interface{}) bool {
	unsubscribed, err := UserIsUnsubscribed(player)
	if err != nil {
		return false
	}
	return !unsubscribed
}

func (uc UnsubscribeCommand) Call(player string, context interface{}, args []string) error {
	u, err := model.LoadUserByEmail(player)
	if err != nil {
		return errors.New("Could not find you in the database")
	}
	if u == nil {
		u = &model.UserModel{
			Email: player,
		}
	}
	u.Unsubscribed = true
	err = u.Save()
	if err != nil {
		return errors.New("Could not mark you as unsubscribed in the database")
	}
	return nil
}

func (uc UnsubscribeCommand) Usage(player string, context interface{}) string {
	return "{{b}}unsubscribe{{_b}} to stop getting any future emails or game invites"
}
