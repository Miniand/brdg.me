package controller

import (
	"fmt"
	"net/http"

	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/api"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/web/view"
	gameView "github.com/Miniand/brdg.me/server/web/view/game"
	"github.com/gorilla/mux"
)

func GameIndex(w http.ResponseWriter, r *http.Request) {
	view.LoggedInUser = GetEmail(r)
	gameView.Index(w)
}

func ApiGameIndex(w http.ResponseWriter, r *http.Request) {
	api.Json("LOL", w, r)
}

func ApiGameShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gm, err := model.LoadGame(vars["id"])
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	g, err := gm.ToGame()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	output, err := g.RenderForPlayer("beefsack@gmail.com")
	if err != nil {
		fmt.Fprint(w, err)
	}
	html, err := render.RenderHtml(output)
	if err != nil {
		fmt.Fprint(w, err)
	}
	api.Json(map[string]interface{}{
		"game": html,
	}, w, r)
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
