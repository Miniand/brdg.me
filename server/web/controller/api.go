package controller

import (
	"net/http"

	"github.com/Miniand/brdg.me/server/model"
)

func WriteCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
}

func ApiOptions(w http.ResponseWriter, r *http.Request) {
	WriteCorsHeaders(w, r)
}

func ApiMustAuthenticate(w http.ResponseWriter, r *http.Request) (bool, *model.UserModel) {
	authorizationToken, found := AuthorizationToken(r)
	if !found {
		ApiError(http.StatusUnauthorized,
			"You must provide an OAuth token using the 'Authorization: token OAUTH-TOKEN' header.",
			w, r)
		return false, nil
	}
	user, loggedIn, err := AuthenticateToken(authorizationToken)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return false, nil
	}
	if !loggedIn {
		ApiError(http.StatusUnauthorized, "You must be logged in", w, r)
		return false, nil
	}
	return loggedIn, user
}

func ApiBadRequest(text string, w http.ResponseWriter, r *http.Request) error {
	return ApiError(http.StatusBadRequest, text, w, r)
}

func ApiInternalServerError(text string, w http.ResponseWriter, r *http.Request) error {
	return ApiError(http.StatusInternalServerError, text, w, r)
}

func ApiUnprocessableEntity(text string, w http.ResponseWriter, r *http.Request) error {
	return ApiError(422, text, w, r)
}

func ApiError(status int, text string, w http.ResponseWriter, r *http.Request) error {
	return Json(status, map[string]string{
		"error": text,
	}, w, r)
}
