package main

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type worker struct {
	conf        *appConfig
	id          uuid.UUID
	conn        *websocket.Conn
	host        string
	port        string
	connectedAt time.Time
	lastSeenAt  time.Time
	status      utils.ObjectStatus
	task_id     *uuid.UUID
}

type workerMessage struct {
	Type    int               `json:"message_type"`
	Payload map[string]string `json:"payload"`
}

func (w *worker) SendWorkerWebsocketMessage(msg websocketMessage) {
	w.conn.WriteMessage(msg.messageType, msg.message)
}

func (w *worker) sendTaskMessage(taskID uuid.UUID, command string) error {
	payload := workerMessage{
		Type:    int(utils.TaskMessage),
		Payload: make(map[string]string, 0),
	}

	payload.Payload["task_id"] = taskID.String()
	payload.Payload["command"] = command

	message := requests.EncodeJSON(payload)

	err := w.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (w *worker) ReadWorkerWebsocketMessage() {
	defer w.conn.Close()
	for {
		messageType, msg, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected connection close from:", w.host)
				w.conf.deleteWorker(w)
				return
			}
			log.Println("Error reading message:", messageType, err)
			return
		}

		jsonMsg, err := requests.DecodeJSON(msg, workerMessage{})
		if err != nil {
			log.Println("Could not decode worker message")
		}

		switch jsonMsg.Type {
		case int(utils.ConnectMessage):
			w.handleConnectMessage(jsonMsg.Payload)
		case int(utils.StatusMessage):
			w.handleStatusMessage(jsonMsg.Payload)
		default:
			log.Println("Unknown Message type:", jsonMsg.Type)
		}

	}
}

func (w *worker) handleConnectMessage(payload map[string]string) {
	w.host = payload["host"]
	w.port = payload["port"]
	log.Printf("Worker connected: %s:%s", w.host, w.port)
}

func (w *worker) handleStatusMessage(payload map[string]string) {
	statusInt, err := strconv.Atoi(payload["status"])
	if err != nil {
		log.Println(err)
		return
	}

	taskStatus := utils.ObjectStatus(statusInt)

	switch taskStatus {
	case utils.StatusRunning:
		w.status = taskStatus
		log.Println("Task running:", *w.task_id)
	case utils.StatusFinished:
		w.status = utils.StatusWaiting
		log.Println("Task finished:", *w.task_id)
		w.task_id = nil
		// payload["output"]
		// payload["task_id"]
	case utils.StatusError:
		w.status = utils.StatusWaiting
		log.Println("Task error:", *w.task_id)
		w.task_id = nil
	}
	w.conf.UpdateState()
}
