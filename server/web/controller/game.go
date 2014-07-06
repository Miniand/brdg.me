package controller

import (
	"fmt"
	"net/http"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/web/view"
	gameView "github.com/Miniand/brdg.me/server/web/view/game"
	"github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
)

func GameIndex(w http.ResponseWriter, r *http.Request) {
	view.LoggedInUser = GetEmail(r)
	gameView.Index(w)
}

func ApiGameIndex(w http.ResponseWriter, r *http.Request) {
	loggedIn, authUser := ApiMustAuthenticate(w, r)
	if !loggedIn {
		return
	}
	games := []map[string]interface{}{}
	var (
		err error
		gm  *model.GameModel
		res *gorethink.Cursor
	)
	switch r.URL.Query().Get("gameState") {
	case "all":
		res, err = model.GamesForPlayer(authUser.Email)
	case "active":
		res, err = model.ActiveGamesForPlayer(authUser.Email)
	default:
		res, err = model.CurrentTurnGamesForPlayer(authUser.Email)
	}
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	for res.Next(&gm) {
		g, err := gm.ToGame()
		if err != nil {
			continue
		}
		games = append(games, map[string]interface{}{
			"id":         gm.Id,
			"name":       g.Name(),
			"identifier": g.Identifier(),
			"isFinished": gm.IsFinished,
			"playerList": gm.PlayerList,
			"whoseTurn":  gm.WhoseTurn,
			"winners":    gm.Winners,
		})
	}
	Json(http.StatusOK, map[string]interface{}{
		"games": games,
	}, w, r)
}

func ApiGameShow(w http.ResponseWriter, r *http.Request) {
	loggedIn, authUser := ApiMustAuthenticate(w, r)
	if !loggedIn {
		return
	}
	vars := mux.Vars(r)
	gm, err := model.LoadGame(vars["id"])
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	g, err := gm.ToGame()
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	gameOutput, err := g.RenderForPlayer(authUser.Email)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	gameHtml, err := render.RenderHtml(gameOutput)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	logHtml, err := render.RenderHtml(
		log.RenderMessages(g.GameLog().MessagesFor(authUser.Email)))
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	commandHtml, err := render.RenderHtml(
		render.CommandUsages(command.CommandUsages(
			authUser.Email, g,
			command.AvailableCommands(authUser.Email, g, g.Commands()))))
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	Json(http.StatusOK, map[string]interface{}{
		"identifier": g.Identifier(),
		"name":       g.Name(),
		"isFinished": g.IsFinished(),
		"whoseTurn":  g.WhoseTurn(),
		"playerList": g.PlayerList(),
		"winners":    g.Winners(),
		"game":       gameHtml,
		"log":        logHtml,
		"commands":   commandHtml,
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

func ApiGameCreate(w http.ResponseWriter, r *http.Request) {
	loggedIn, authUser := ApiMustAuthenticate(w, r)
	if !loggedIn {
		return
	}
	identifier := r.PostFormValue("identifier")
	g := game.RawCollection()[identifier]
	if g == nil {
		ApiBadRequest(fmt.Sprintf(
			"Could not find a game with the identifier '%s'", identifier),
			w, r)
		return
	}
	players := r.PostForm["opponents[]"]
	if players == nil {
		players = []string{}
	}
	players = append(players, authUser.Email)
	gm, err := model.StartNewGame(g, players)
	if err != nil {
		ApiBadRequest(err.Error(), w, r)
		return
	}
	Json(http.StatusOK, map[string]interface{}{
		"id": gm.Id,
	}, w, r)
}
