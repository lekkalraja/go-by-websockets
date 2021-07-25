package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	r := routes()
	port := 8080
	log.Printf("Starting Web Server on : %d \n", port)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
