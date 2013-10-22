package route

import (
	"github.com/Miniand/brdg.me/server/web/controller"
	"github.com/gorilla/mux"
)

var router *mux.Router

func Router() *mux.Router {
	if router == nil {
		router = mux.NewRouter()
		router.HandleFunc("/", controller.Root).Name("root")
		router.HandleFunc("/game", controller.GameIndex).Name("gameIndex")
		router.HandleFunc("/game/{id}", controller.GameShow).Name("gameShow")
		router.HandleFunc("/game/new/{identifier}", controller.GameNew).Name(
			"gameNew")
	}
	return router
}