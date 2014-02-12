package email

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/server/model"
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

func (uc UnsubscribeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	u, err := model.FirstUserByEmail(player)
	if err != nil {
		return "", errors.New("Could not find you in the database")
	}
	if u == nil {
		u = &model.UserModel{
			Email: player,
		}
	}
	u.Unsubscribed = true
	err = u.Save()
	if err != nil {
		return "", errors.New("Could not mark you as unsubscribed in the database")
	}
	return "You have unsubscribed from brdg.me", nil
}

func (uc UnsubscribeCommand) Usage(player string, context interface{}) string {
	// We don't want it to show in the usage section
	return ""
}
