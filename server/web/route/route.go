package route

import (
	"code.google.com/p/go.net/websocket"
	"github.com/Miniand/brdg.me/server/web/controller"
	"github.com/gorilla/mux"
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
		router.Handle("/ws", websocket.Handler(controller.Ws))
		api := router.PathPrefix("/api/").Subrouter()
		api.HandleFunc("/game", controller.ApiGameIndex).Methods("GET").
			Name("apiGameIndex")
		api.HandleFunc("/game", controller.ApiGameCreate).Methods("POST").
			Name("apiGameCreate")
		api.HandleFunc("/game/{id}", controller.ApiGameShow).Methods("GET").
			Name("apiGameShow")
	}
	return router
}
