package controller

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/render"
	sgame "github.com/Miniand/brdg.me/server/game"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/scommand"

	"github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
)

const (
	ParamRenderer = "renderer"
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

func findRenderer(rendererParam string) (render.Renderer, error) {
	switch rendererParam {
	case "", "html":
		return render.RenderHtml, nil
	case "ansi":
		return render.RenderTerminal, nil
	case "raw":
		return func(tmpl string) (string, error) {
			return tmpl, nil
		}, nil
	case "plain":
		return func(tmpl string) (string, error) {
			return render.RenderPlain(tmpl), nil
		}, nil
	default:
		return nil, errors.New(
			"unknown renderer, must be one in html, ansi, raw or plain")
	}
}

func GameData(
	gm *model.GameModel,
	g game.Playable,
	renderer render.Renderer,
) (data map[string]interface{}, err error) {
	playerList := gm.PlayerList
	if playerList != nil {
		if playerList, err = render.RenderTemplates(render.PlayerNamesInPlayers(
			playerList,
			gm.PlayerList,
		), renderer); err != nil {
			return
		}
	}
	whoseTurn := gm.WhoseTurn
	if whoseTurn != nil {
		if whoseTurn, err = render.RenderTemplates(render.PlayerNamesInPlayers(
			whoseTurn,
			gm.PlayerList,
		), renderer); err != nil {
			return
		}
	}
	winners := gm.Winners
	if winners != nil {
		if winners, err = render.RenderTemplates(render.PlayerNamesInPlayers(
			winners,
			gm.PlayerList,
		), renderer); err != nil {
			return
		}
	}
	data = map[string]interface{}{
		"id":         gm.Id,
		"name":       g.Name(),
		"identifier": g.Identifier(),
		"isFinished": gm.IsFinished,
		"finishedAt": gm.FinishedAt,
		"playerList": playerList,
		"whoseTurn":  whoseTurn,
		"winners":    winners,
	}
	return
}

func GameOutput(
	player string,
	gm *model.GameModel,
	g game.Playable,
	renderer render.Renderer,
) (map[string]interface{}, error) {
	gameOutput, err := g.RenderForPlayer(player)
	if err != nil {
		return nil, err
	}
	gameRender, err := renderer(gameOutput)
	if err != nil {
		return nil, err
	}
	logs := []map[string]interface{}{}
	for _, l := range g.GameLog().MessagesFor(player) {
		logRender, err := renderer(l.Text)
		if err != nil {
			return nil, err
		}
		t := time.Unix(l.Time/int64(math.Pow10(9)), 0)
		logs = append(logs, map[string]interface{}{
			"time": t.UTC().Format(time.RFC3339),
			"text": logRender,
		})
	}
	if err != nil {
		return nil, err
	}
	commandRender, err := renderer(
		render.CommandUsages(command.CommandUsages(
			player, g,
			command.AvailableCommands(player, g,
				scommand.CommandsForGame(gm, g)))))
	if err != nil {
		return nil, err
	}
	lastRead := time.Unix(g.GameLog().LastReadTimeFor[player]/int64(math.Pow10(9)), 0)
	return map[string]interface{}{
		"game":         gameRender,
		"log":          logs,
		"lastReadTime": lastRead.UTC().Format(time.RFC3339),
		"commands":     commandRender,
	}, nil
}

func ApiGameIndex(w http.ResponseWriter, r *http.Request) {
	loggedIn, authUser := ApiMustAuthenticate(w, r)
	if !loggedIn {
		return
	}
	var (
		err error
		gm  *model.GameModel
		res *gorethink.Cursor
	)
	query := r.URL.Query()
	renderer, err := findRenderer(query.Get(ParamRenderer))
	if err != nil {
		ApiBadRequest(err.Error(), w, r)
		return
	}
	games := []map[string]interface{}{}
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
		gd, err := GameData(gm, g, renderer)
		if err != nil {
			continue
		}
		games = append(games, gd)
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
	query := r.URL.Query()
	renderer, err := findRenderer(query.Get(ParamRenderer))
	if err != nil {
		ApiBadRequest(err.Error(), w, r)
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
	gameOutput, err := GameOutput(authUser.Email, gm, g, renderer)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	gd, err := GameData(gm, g, renderer)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	Json(http.StatusOK, mergeMaps(gd, gameOutput), w, r)
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
	renderer, err := findRenderer(r.FormValue(ParamRenderer))
	if err != nil {
		ApiBadRequest(err.Error(), w, r)
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
	gameOutput, err := GameOutput(authUser.Email, gm, g, renderer)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	gd, err := GameData(gm, g, renderer)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	Json(http.StatusOK, mergeMaps(gd, gameOutput), w, r)
}

func ApiGameSummary(w http.ResponseWriter, r *http.Request) {
	var gm *model.GameModel
	loggedIn, authUser := ApiMustAuthenticate(w, r)
	if !loggedIn {
		return
	}
	query := r.URL.Query()
	renderer, err := findRenderer(query.Get(ParamRenderer))
	if err != nil {
		ApiBadRequest(err.Error(), w, r)
		return
	}
	resp := map[string][]map[string]interface{}{
		"currentTurn":      []map[string]interface{}{},
		"recentlyFinished": []map[string]interface{}{},
		"otherActive":      []map[string]interface{}{},
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
		gd, err := GameData(gm, g, renderer)
		if err != nil {
			ApiInternalServerError(err.Error(), w, r)
		}
		resp["currentTurn"] = append(resp["currentTurn"], gd)
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
		gd, err := GameData(gm, g, renderer)
		if err != nil {
			continue
		}
		resp["recentlyFinished"] = append(resp["recentlyFinished"], gd)
	}
	// Other active
	otherRes, err := model.NotCurrentTurnGamesForPlayer(authUser.Email)
	if err != nil {
		ApiInternalServerError(err.Error(), w, r)
	}
	defer otherRes.Close()
	for otherRes.Next(&gm) {
		g, err := gm.ToGame()
		if err != nil {
			continue
		}
		gd, err := GameData(gm, g, renderer)
		if err != nil {
			continue
		}
		resp["otherActive"] = append(resp["otherActive"], gd)
	}
	Json(http.StatusOK, resp, w, r)
}
