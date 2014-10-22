package controller

import (
	"fmt"
	"net/http"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	sgame "github.com/Miniand/brdg.me/server/game"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/scommand"

	"github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
)

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

func ApiGameTypeIndex(w http.ResponseWriter, r *http.Request) {
	gameTypes := []map[string]interface{}{}
	for _, g := range game.RawCollection() {
		gameTypes = append(gameTypes, map[string]interface{}{
			"identifier": g.Identifier(),
			"name":       g.Name(),
		})
	}
	Json(http.StatusOK, gameTypes, w, r)
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
			command.AvailableCommands(authUser.Email, g,
				append(g.Commands(), scommand.Commands(gm)...)))))
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

func ApiGameCommand(w http.ResponseWriter, r *http.Request) {
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
	if err := sgame.HandleCommandText(authUser.Email, gm.Id,
		r.FormValue("command")); err != nil {
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
