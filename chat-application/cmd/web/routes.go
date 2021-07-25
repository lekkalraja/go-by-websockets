package main

import (
	"github.com/gorilla/mux"
	"github.com/lekkalraja/go-by-websockets/chat-application/handlers"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Home)
	r.HandleFunc("/ws", handlers.WsEndpoint)
	return r
}
