package controller

import (
	"fmt"
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
	gm, err := model.LoadGame(vars["id"])
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
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
		Game:    g,
		Players: r.URL.Query()["players[]"],
	})
}

func GameCreate(w http.ResponseWriter, r *http.Request) {
	view.LoggedInUser = GetEmail(r)
	g := game.RawCollection()[r.PostFormValue("identifier")]
	if g == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	players := r.PostForm["players[]"]
	players = append(players, view.LoggedInUser)
	gm, err := model.StartNewGame(g, players)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		fmt.Println(gm.Id)
		http.Redirect(w, r, "/game/"+gm.Id, http.StatusFound)
	}
}
