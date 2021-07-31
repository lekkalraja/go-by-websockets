package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lekkalraja/go-by-websockets/chat-application/handlers"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Home)
	r.HandleFunc("/ws", handlers.WsEndpoint)
	fileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
	return r
}
