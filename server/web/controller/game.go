package controller

import (
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
	view "github.com/Miniand/brdg.me/server/web/view/game"
	"github.com/gorilla/mux"
	"net/http"
)

func GameIndex(w http.ResponseWriter, r *http.Request) {
	view.Index(w)
}

func GameShow(w http.ResponseWriter, r *http.Request) {
	g := game.RawCollection()["acquire"]
	g.Start([]string{"Mick", "Steve"})
	gm, err := model.GameToGameModel(g)
	if err != nil {
		panic(err.Error())
	}
	view.Show(w, view.ShowScope{
		GameModel: gm,
	})
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
