package scommand

import (
	"errors"
	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/server/model"
)

type SubscribeCommand struct{}

func (sc SubscribeCommand) Parse(input string) []string {
	return comm.ParseNamedCommandNArgs("subscribe", 0, input)
}

func (sc SubscribeCommand) CanCall(player string, context interface{}) bool {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || !ok {
		return false
	}
	return u.Unsubscribed
}

func (sc SubscribeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || !ok {
		return "", errors.New("Could not find you in the database")
	}
	if u == nil {
		u = &model.UserModel{
			Email: player,
		}
	}
	u.Unsubscribed = false
	err = u.Save()
	if err != nil {
		return "", errors.New("Could not mark you as subscribed in the database")
	}
	return "You are now able to access brdg.me again", nil
}

func (sc SubscribeCommand) Usage(player string, context interface{}) string {
	return "{{b}}subscribe{{_b}} to get access to brdg.me, you are currently unsubscribed from brdg.me"
}
