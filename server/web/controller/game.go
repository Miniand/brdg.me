package controller

import (
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/web/view"
	gameView "github.com/Miniand/brdg.me/server/web/view/game"
	"github.com/gorilla/mux"
	"net/http"
)

func GameIndex(w http.ResponseWriter, r *http.Request) {
	view.LoggedInUser = GetEmail(r)
	gameView.Index(w)
}

func GameShow(w http.ResponseWriter, r *http.Request) {
	view.LoggedInUser = GetEmail(r)
	vars := mux.Vars(r)
	g := game.RawCollection()[vars["id"]]
	if g == nil {
		g = game.RawCollection()["acquire"]
	}
	g.Start([]string{"Mick", "Steve"})
	gm, err := model.GameToGameModel(g)
	if err != nil {
		panic(err.Error())
	}
	gameView.Show(w, gameView.ShowScope{
		GameModel: gm,
	})
}

func GameNew(w http.ResponseWriter, r *http.Request) {
	view.LoggedInUser = GetEmail(r)
	vars := mux.Vars(r)
	g := game.RawCollection()[vars["identifier"]]
	if g == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	gameView.New(w, gameView.NewScope{
		Game: g,
	})
}
