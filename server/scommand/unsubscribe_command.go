package scommand

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/server/model"
)

type UnsubscribeCommand struct{}

func (uc UnsubscribeCommand) Name() string { return "unsubscribe" }

func (uc UnsubscribeCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
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

func CanUnsubscribe(player string, gm *model.GameModel) bool {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil {
		return false
	}
	return !ok || !u.Unsubscribed
}
