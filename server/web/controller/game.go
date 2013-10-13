package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GameIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Game index")
}

func GameShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Game: "+vars["id"])
}
