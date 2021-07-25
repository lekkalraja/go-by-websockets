package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"github.com/lekkalraja/go-by-websockets/chat-application/models"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(), // remove in production
)

// upgrader websocket conneciton configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Home Landing endpoint
func Home(w http.ResponseWriter, r *http.Request) {
	err := render(w, "home", nil)
	if err != nil {
		log.Printf("Something went wrong while rendering : %v", err)
	}
}

// WsEndpoint upgrades to websocket and write response back
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Something went wrong while Upgrading to WebSocket : %v", err)
	}
	log.Printf("Successfully Upgraded to WS")
	response := models.WsResponse{
		Message: `<em><small>Connected to server</small></em>`,
	}

	err = conn.WriteJSON(response)

	if err != nil {
		log.Printf("Something went wrong while Wring Ws response : %v", err)
	}
}

// render jet template
func render(w http.ResponseWriter, path string, variables jet.VarMap) error {
	template, err := views.GetTemplate(path)

	if err != nil {
		return err
	}

	err = template.Execute(w, variables, nil)

	if err != nil {
		return err
	}

	return nil
}
