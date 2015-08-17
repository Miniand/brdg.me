package scommand

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/server/model"
)

type SubscribeCommand struct{}

func (sc SubscribeCommand) Name() string { return "subscribe" }

func (sc SubscribeCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
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

func CanSubscribe(player string, gm *model.GameModel) bool {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || !ok {
		return false
	}
	return u.Unsubscribed
}
