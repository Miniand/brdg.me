package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", RootHandler)
	addr := os.Getenv("BRDGME_WEB_SERVER_ADDRESS")
	if addr == "" {
		addr = ":9998"
	}
	log.Print("Running web server on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}

func Title() string {
	return " _             _                        \n" +
		"| |__  _ __ __| | __ _   _ __ ___   ___ \n" +
		"| '_ \\| '__/ _` |/ _` | | '_ ` _ \\ / _ \\\n" +
		"| |_) | | | (_| | (_| |_| | | | | |  __/\n" +
		"|_.__/|_|  \\__,_|\\__, (_)_| |_| |_|\\___|\n" +
		"                 |___/"
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><head><style>"+Style+"</style></head><body><header><h1>"+Title()+"\n\n"+
		"Play board games online via web and email</h1></header></body></html>")
}
