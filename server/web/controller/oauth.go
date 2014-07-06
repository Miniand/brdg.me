package controller

import (
	"net/http"
	"regexp"

	"github.com/Miniand/brdg.me/server/model"
)

var tokenRegexp = regexp.MustCompile(`^\s*token\s+([^\s]+)\s*$`)

func AuthorizationToken(r *http.Request) (string, bool) {
	matches := tokenRegexp.FindStringSubmatch(r.Header.Get("Authorization"))
	if matches == nil {
		return "", false
	}
	return matches[1], true
}

func AuthenticateToken(authorizationToken string) (*model.UserModel, bool, error) {
	return &model.UserModel{
		Email: authorizationToken,
	}, true, nil
}
