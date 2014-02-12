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
		router.HandleFunc("/", controller.Root).Name("root")
		router.HandleFunc("/session/sign-in", controller.SessionSignIn).
			Methods("POST").Name("sessionSignIn")
		router.HandleFunc("/session/sign-out", controller.SessionSignOut).
			Methods("POST").Name("sessionSignOut")
		router.HandleFunc("/", controller.Root).Name("sessionSignOut")
		router.HandleFunc("/game", controller.GameIndex).Name("gameIndex")
		router.HandleFunc("/game/{id}", controller.GameShow).Name("gameShow")
		router.HandleFunc("/game/new/{identifier}", controller.GameNew).Name(
			"gameNew")
		router.Handle("/ws", websocket.Handler(controller.Ws))
	}
	return router
}
