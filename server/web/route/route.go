package route

import (
	"github.com/gorilla/mux"

	"github.com/Miniand/brdg.me/server/web/controller"
)

var router *mux.Router

func Router() *mux.Router {
	if router == nil {
		router = mux.NewRouter()
		router.HandleFunc("/", controller.Root).Methods("GET").Name("root")
		router.HandleFunc("/session/sign-in", controller.SessionSignIn).
			Methods("POST").Name("sessionSignIn")
		router.HandleFunc("/session/sign-out", controller.SessionSignOut).
			Methods("POST").Name("sessionSignOut")
		router.HandleFunc("/", controller.Root).Methods("GET").
			Name("sessionSignOut")
		router.HandleFunc("/game", controller.GameIndex).Methods("GET").
			Name("gameIndex")
		router.HandleFunc("/game", controller.GameCreate).Methods("POST").
			Name("gameCreate")
		router.HandleFunc("/game/{id}", controller.GameShow).Methods("GET").
			Name("gameShow")
		router.HandleFunc("/game/new/{identifier}", controller.GameNew).
			Methods("GET").Name("gameNew")
		router.HandleFunc("/ws", controller.Websocket)
		router.HandleFunc("/auth/request", controller.AuthRequest).Methods("POST").
			Name("authRequest")
		router.HandleFunc("/auth/confirm", controller.AuthConfirm).Methods("POST").
			Name("authConfirm")
		api := router.PathPrefix("/api/").Subrouter()
		api.HandleFunc("/{url:.*}", controller.ApiOptions).Methods("OPTIONS").
			Name("apiOptions")
		api.HandleFunc("/game", controller.ApiGameIndex).Methods("GET").
			Name("apiGameIndex")
		api.HandleFunc("/game", controller.ApiGameCreate).Methods("POST").
			Name("apiGameCreate")
		api.HandleFunc("/game/{id}", controller.ApiGameShow).Methods("GET").
			Name("apiGameShow")
		api.HandleFunc("/game/{id}/command", controller.ApiGameCommand).Methods("POST").
			Name("apiGameCommand")
	}
	return router
}
