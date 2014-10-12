package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/server/email"
	"github.com/Miniand/brdg.me/server/model"
)

var tokenRegexp = regexp.MustCompile(`^\s*token\s+([^\s]+)\s*$`)

func isValidEmail(emailAddr string) bool {
	return regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$`).
		MatchString(emailAddr)
}

func AuthorizationToken(r *http.Request) (string, bool) {
	matches := tokenRegexp.FindStringSubmatch(r.Header.Get("Authorization"))
	if matches == nil {
		return "", false
	}
	return matches[1], true
}

func AuthenticateToken(authorizationToken string) (*model.UserModel, bool, error) {
	if isValidEmail(authorizationToken) {
		return &model.UserModel{
			Email: authorizationToken,
		}, true, nil
	}
	token, ok, err := model.FindAuthToken(authorizationToken)
	if err != nil || !ok ||
		token.CreatedAt.Before(time.Now().AddDate(0, 0, -30)) {
		return nil, false, err
	}
	um, err := model.LoadUser(token.UserId)
	return um, err == nil, err
}

func AuthRequest(w http.ResponseWriter, r *http.Request) {
	emailAddr := strings.TrimSpace(r.PostFormValue("email"))
	if emailAddr == "" {
		ApiBadRequest("Please pass an email parameter", w, r)
		return
	}
	if !isValidEmail(emailAddr) {
		ApiBadRequest(fmt.Sprintf("%s is not a valid email", emailAddr), w, r)
		return
	}
	user, ok, err := model.FirstUserByEmail(emailAddr)
	if err != nil {
		ApiBadRequest("Error finding user", w, r)
		return
	}
	if !ok {
		user = &model.UserModel{
			Email: emailAddr,
		}
	}
	user.AuthRequest = model.GenerateAuthRequestToken()
	user.AuthRequestedAt = time.Now()
	if err := user.Save(); err != nil {
		ApiBadRequest("Error creating auth request token", w, r)
		return
	}
	if err := email.SendMail([]string{emailAddr},
		fmt.Sprintf(`Your brdg.me confirmation is %s`, user.AuthRequest)); err != nil {
		ApiBadRequest("Error emailing auth request token", w, r)
		return
	}
}

func AuthConfirm(w http.ResponseWriter, r *http.Request) {
	emailAddr := strings.TrimSpace(r.PostFormValue("email"))
	confirmation := strings.TrimSpace(r.PostFormValue("confirmation"))
	if emailAddr == "" {
		ApiBadRequest("Please pass an email parameter", w, r)
		return
	}
	if !isValidEmail(emailAddr) {
		ApiBadRequest(fmt.Sprintf("%s is not a valid email", emailAddr), w, r)
		return
	}
	if confirmation == "" {
		ApiBadRequest("Please pass a confirmation parameter, which should be found in the email sent during the auth request", w, r)
		return
	}
	user, ok, err := model.FirstUserByEmail(emailAddr)
	if err != nil {
		ApiBadRequest("Error finding user", w, r)
		return
	}
	if !ok || user.AuthRequest != confirmation || user.AuthRequestedAt.Before(
		time.Now().Add(-30*time.Minute)) {
		ApiBadRequest("Invalid confirmation", w, r)
		return
	}
	user.AuthRequest = ""
	if err := user.Save(); err != nil {
		ApiBadRequest("Error generating token", w, r)
		return
	}
	token, err := model.NewAuthToken(user.Id)
	if err != nil {
		ApiBadRequest("Error generating token", w, r)
		return
	}
	if err := token.Save(); err != nil {
		ApiBadRequest("Error generating token", w, r)
		return
	}
	Json(http.StatusOK, token.Token, w, r)
}
