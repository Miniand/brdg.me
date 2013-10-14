package controller

import (
	"github.com/Miniand/brdg.me/game"
	view "github.com/Miniand/brdg.me/server/web/view/game"
	"github.com/gorilla/mux"
	"net/http"
)

func GameIndex(w http.ResponseWriter, r *http.Request) {
	view.Index(w)
}

func GameShow(w http.ResponseWriter, r *http.Request) {
	view.Show(w)
}

func GameNew(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	g := game.RawCollection()[vars["identifier"]]
	if g == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	view.New(w, view.NewScope{
		Game: g,
	})
}
