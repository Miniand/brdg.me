package main

import (
	"github.com/beefsack/brdg.me/server/model"
)

func UserIsUnsubscribed(email string) (bool, error) {
	u, err := model.LoadUserByEmail(email)
	if err != nil {
		return false, err
	}
	return u != nil && u.Unsubscribed, nil
}
