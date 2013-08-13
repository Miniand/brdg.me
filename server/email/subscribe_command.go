package main

import (
	"errors"
	"github.com/beefsack/brdg.me/command"
	"github.com/beefsack/brdg.me/server/model"
)

type SubscribeCommand struct{}

func (sc SubscribeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("subscribe", 0, input)
}

func (sc SubscribeCommand) CanCall(player string, context interface{}) bool {
	unsubscribed, err := UserIsUnsubscribed(player)
	if err != nil {
		return false
	}
	return unsubscribed
}

func (sc SubscribeCommand) Call(player string, context interface{}, args []string) error {
	u, err := model.FirstUserByEmail(player)
	if err != nil {
		return errors.New("Could not find you in the database")
	}
	if u == nil {
		u = &model.UserModel{
			Email: player,
		}
	}
	u.Unsubscribed = false
	err = u.Save()
	if err != nil {
		return errors.New("Could not mark you as subscribed in the database")
	}
	return nil
}

func (sc SubscribeCommand) Usage(player string, context interface{}) string {
	return "{{b}}subscribe{{_b}} to get access to brdg.me, you are currently unsubscribed from brdg.me"
}
