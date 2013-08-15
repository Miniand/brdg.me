package main

import (
	"github.com/Miniand/brdg.me/server/model"
)

func UserIsUnsubscribed(email string) (bool, error) {
	u, err := model.FirstUserByEmail(email)
	if err != nil {
		return false, err
	}
	return u != nil && u.Unsubscribed, nil
}
