package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lekkalraja/go-by-websockets/chat-application/handlers"
)

func main() {
	mux := routes()
	port := 8080
	log.Printf("Starting Web Server on : %d \n", port)

	go handlers.ListenPayloadChannel()

	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
