package models

import "github.com/gorilla/websocket"

type WsResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action   string         `json:"action"`
	UserName string         `json:"user_name"`
	Message  string         `json:"message"`
	Conn     websocket.Conn `json:"-"`
}
