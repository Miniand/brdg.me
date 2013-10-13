package controller

import (
	"github.com/Miniand/brdg.me/server/web/view"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	view.Root(w)
}
