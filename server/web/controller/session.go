package controller

import (
	"encoding/json"
	"github.com/Miniand/brdg.me/server/web/config"
	"github.com/gorilla/securecookie"
	"github.com/sauerbraten/persona"
	"net/http"
	"time"
)

// initialize secure cookie storage; 16 byte ~ 128 bit AES encryption
var secCookie *securecookie.SecureCookie = securecookie.New(
	[]byte("jkf342hf3kj21l21kjch"), []byte{0x5d, 0xa5, 0xd3, 0x90, 0xc9, 0x54,
		0xa1, 0xc3, 0x70, 0x00, 0x8d, 0x6d, 0xa9, 0xd1, 0x07, 0x53})

func SessionSignIn(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	response, err := persona.VerifyAssertion(config.Get(config.SERVER_ADDRESS),
		r.FormValue("assertion"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.OK() {
		setSessionCookie(w, response.Email, response.Expires)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	enc.Encode(response)
}

func SessionSignOut(w http.ResponseWriter, r *http.Request) {
	revokeSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}

// looks for a session cookie and returns the user's email address, or "" if something failed.
// this example has no error handling!
func GetEmail(r *http.Request) (email string) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	err = secCookie.Decode("session", cookie.Value, &email)
	return
}

// sets a secure session cookie containing the user's email address, so we can recognize him later
func setSessionCookie(w http.ResponseWriter, email string, expires int64) error {
	encoded, err := secCookie.Encode("session", email)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 336), // session is valid for 2 weeks
	})

	return nil
}

// overwrites the secure session cookie
func revokeSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session", // same cookie name → overwrites session cookie
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // browser deletes this cookie immediatly
	})
}
