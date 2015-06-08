package route

import (
	"github.com/gorilla/mux"

	"github.com/Miniand/brdg.me/server/web/controller"
)

var router *mux.Router

func Router() *mux.Router {
	if router == nil {
		router = mux.NewRouter()
		router.HandleFunc("/{url:.*}", controller.ApiOptions).Methods("OPTIONS").
			Name("apiOptions")
		router.HandleFunc("/game_type", controller.ApiGameTypeIndex).Methods("GET").
			Name("apiGameTypeIndex")
		router.HandleFunc("/game", controller.ApiGameIndex).Methods("GET").
			Name("apiGameIndex")
		router.HandleFunc("/game", controller.ApiGameCreate).Methods("POST").
			Name("apiGameCreate")
		router.HandleFunc("/game/summary", controller.ApiGameSummary).Methods("GET").
			Name("apiGameSummary")
		router.HandleFunc("/game/{id}", controller.ApiGameShow).Methods("GET").
			Name("apiGameShow")
		router.HandleFunc("/game/{id}/command", controller.ApiGameCommand).Methods("POST").
			Name("apiGameCommand")
		router.HandleFunc("/ws", controller.Websocket)
		router.HandleFunc("/auth/request", controller.AuthRequest).Methods("POST").
			Name("authRequest")
		router.HandleFunc("/auth/confirm", controller.AuthConfirm).Methods("POST").
			Name("authConfirm")
	}
	return router
}
