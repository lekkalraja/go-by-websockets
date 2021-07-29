package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"github.com/lekkalraja/go-by-websockets/chat-application/models"
)

var conns = make(map[*websocket.Conn]string)
var payChan = make(chan *models.WsPayload)

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
	conns[conn] = ""
	if err != nil {
		log.Printf("Something went wrong while Wring Ws response : %v", err)
	}
	go listenWS(conn)
}

func listenWS(conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Something went wrong, Recovering for the error : %v", r)
			delete(conns, conn)
		}
	}()
	var payload *models.WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err == nil {
			payload.Conn = *conn
			log.Printf("Payload : %v \n", payload.Action)
			payChan <- payload
		}
	}
}

func ListenPayloadChannel() {
	for {
		payload := <-payChan
		switch payload.Action {
		case "open":
			brodcast(getWsResponse("list_users"))
		case "username":
			conns[&payload.Conn] = payload.UserName
			brodcast(getWsResponse("list_users"))
		case "left":
			delete(conns, &payload.Conn)
			brodcast(getWsResponse("list_users"))
		case "message":
			response := &models.WsResponse{
				ConnectedUsers: getUsers(),
				Action:         "broadcast",
				Message:        fmt.Sprintf("<strong>%s</strong> : %s", payload.UserName, payload.Message),
			}
			brodcast(response)
		}

	}
}

func getWsResponse(action string) *models.WsResponse {
	return &models.WsResponse{
		ConnectedUsers: getUsers(),
		Action:         action,
	}
}

func getUsers() []string {
	var users []string
	for _, user := range conns {
		if user != "" {
			users = append(users, user)
		}
	}
	sort.Strings(users)
	log.Printf("Users : %v \n", users)
	return users
}

func brodcast(response *models.WsResponse) {
	log.Printf("Response : %v \n", response)
	for conn, user := range conns {
		if user != "" {
			log.Printf("Sending Response : %v  to %s\n", response, user)
			err := conn.WriteJSON(response)
			if err != nil {
				log.Printf("Something went wrong while sending message to : %v", conn.RemoteAddr())
				conn.Close()
				delete(conns, conn)
			}
		}
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
