package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsMessage struct {
	message     []byte
	messageType int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (w *worker) SendWsMessage(msg wsMessage) {
	w.conn.WriteMessage(msg.messageType, msg.message)
}

func (w *worker) ReadWsMessage() {
	defer w.conn.Close()
	for {
		messageType, p, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", messageType, err)
			return
		}
		log.Println(messageType, string(p))
	}
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
