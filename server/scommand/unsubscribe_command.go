package scommand

import (
	"errors"
	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/server/model"
)

type UnsubscribeCommand struct{}

func (uc UnsubscribeCommand) Parse(input string) []string {
	return comm.ParseNamedCommandNArgs("unsubscribe", 0, input)
}

func (uc UnsubscribeCommand) CanCall(player string, context interface{}) bool {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil {
		return false
	}
	return !ok || !u.Unsubscribed
}

func (uc UnsubscribeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil {
		return "", errors.New("Could not find you in the database")
	}
	if !ok {
		u = &model.UserModel{
			Email: player,
		}
	}
	u.Unsubscribed = true
	if err := u.Save(); err != nil {
		return "", errors.New("Could not mark you as unsubscribed in the database")
	}
	return "You have unsubscribed from brdg.me", nil
}

func (uc UnsubscribeCommand) Usage(player string, context interface{}) string {
	// We don't want it to show in the usage section
	return ""
}