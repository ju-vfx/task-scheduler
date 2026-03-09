package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type websocketMessage struct {
	message     []byte
	messageType int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func UpgradeConnection(w http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	// Allow all connections
	upgrader.CheckOrigin = func(req *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}
