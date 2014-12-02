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

func mergeMaps(in ...map[string]interface{}) map[string]interface{} {
	var first map[string]interface{}
	for _, i := range in {
		if first == nil {
			first = i
			continue
		}
		for k, v := range i {
			first[k] = v
		}
	}
	return first
}

func GameData(gm *model.GameModel, g game.Playable) map[string]interface{} {
	return map[string]interface{}{
		"id":         gm.Id,
		"name":       g.Name(),
		"identifier": g.Identifier(),
		"isFinished": gm.IsFinished,
		"finishedAt": gm.FinishedAt,
		"playerList": gm.PlayerList,
		"whoseTurn":  gm.WhoseTurn,
		"winners":    gm.Winners,
	}
}

func GameOutput(
	player string,
	gm *model.GameModel,
	g game.Playable,
) (map[string]interface{}, error) {
	gameOutput, err := g.RenderForPlayer(player)
	if err != nil {
		return nil, err
	}
	gameHtml, err := render.RenderHtml(gameOutput)
	if err != nil {
		return nil, err
	}
	logHtml, err := render.RenderHtml(
		log.RenderMessages(g.GameLog().MessagesFor(player)))
	if err != nil {
		return nil, err
	}
	commandHtml, err := render.RenderHtml(
		render.CommandUsages(command.CommandUsages(
			player, g,
			command.AvailableCommands(player, g,
				append(g.Commands(), scommand.Commands(gm)...)))))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"game":     gameHtml,
		"log":      logHtml,
		"commands": commandHtml,
	}, nil
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
	case "recentlyFinished":
		res, err = model.RecentlyFinishedGamesForPlayer(authUser.Email)
	default:
		res, err = model.CurrentTurnGamesForPlayer(authUser.Email)
	}
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	defer res.Close()
	for res.Next(&gm) {
		g, err := gm.ToGame()
		if err != nil {
			continue
		}
		games = append(games, GameData(gm, g))
	}
	Json(http.StatusOK, games, w, r)
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
	gameOutput, err := GameOutput(authUser.Email, gm, g)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	Json(http.StatusOK, mergeMaps(GameData(gm, g), gameOutput), w, r)
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
	gameId := vars["id"]
	// Do the command
	if err := sgame.HandleCommandText(authUser.Email, gameId,
		r.FormValue("command")); err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	// Load game to get changes
	gm, err := model.LoadGame(gameId)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	g, err := gm.ToGame()
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	// Mark log as read for player and resave
	g.GameLog().MarkReadFor(authUser.Email)
	if err := gm.UpdateState(g); err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	if err := gm.Save(); err != nil {
		ApiInternalServerError(err.Error(), w, r)
		return
	}
	// Get output for return
	gameOutput, err := GameOutput(authUser.Email, gm, g)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	Json(http.StatusOK, mergeMaps(GameData(gm, g), gameOutput), w, r)
}

func ApiGameSummary(w http.ResponseWriter, r *http.Request) {
	var gm *model.GameModel
	loggedIn, authUser := ApiMustAuthenticate(w, r)
	if !loggedIn {
		return
	}
	resp := map[string][]map[string]interface{}{
		"currentTurn":      []map[string]interface{}{},
		"recentlyFinished": []map[string]interface{}{},
	}
	// Current turn
	currentRes, err := model.CurrentTurnGamesForPlayer(authUser.Email)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	defer currentRes.Close()
	for currentRes.Next(&gm) {
		g, err := gm.ToGame()
		if err != nil {
			continue
		}
		resp["currentTurn"] = append(resp["currentTurn"], GameData(gm, g))
	}
	// Recently finished
	finishedRes, err := model.RecentlyFinishedGamesForPlayer(authUser.Email)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	defer finishedRes.Close()
	for finishedRes.Next(&gm) {
		g, err := gm.ToGame()
		if err != nil {
			continue
		}
		resp["recentlyFinished"] = append(resp["recentlyFinished"], GameData(gm, g))
	}
	Json(http.StatusOK, resp, w, r)
}
