package web

import (
	"github.com/Miniand/brdg.me/server/web/config"
	"github.com/Miniand/brdg.me/server/web/route"
	"log"
	"net/http"
)

func Run() error {
	http.Handle("/", route.Router())
	addr := config.Get(config.SERVER_ADDRESS)
	log.Print("Running web server on " + addr)
	return http.ListenAndServe(addr, nil)
}
