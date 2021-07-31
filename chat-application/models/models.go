package models

import "github.com/gorilla/websocket"

type WsResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsConnection struct {
	Conn *websocket.Conn
}
type WsPayload struct {
	Action   string       `json:"action"`
	UserName string       `json:"username"`
	Message  string       `json:"message"`
	Conn     WsConnection `json:"-"`
}
