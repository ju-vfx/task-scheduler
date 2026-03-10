package main

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type workerResp struct {
	ID          string `json:"id"`
	Host        string `json:"host"`
	ConnectedAt string `json:"connected_at"`
	LastSeenAt  string `json:"last_seen_at"`
	Status      string `json:"status"`
}

func (conf *appConfig) handlerRegisterWorker(w http.ResponseWriter, req *http.Request) {

	ws, err := UpgradeConnection(w, req)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Could not connect to websocket")
		return
	}

	addrPieces := strings.Split(ws.RemoteAddr().String(), ":")
	host := addrPieces[0]
	port := addrPieces[1]

	id, _ := uuid.NewUUID()
	wrk := &worker{
		conf:        conf,
		id:          id,
		conn:        ws,
		host:        host,
		port:        port,
		connectedAt: time.Now(),
		lastSeenAt:  time.Now(),
		status:      utils.StatusWaiting,
		task_id:     nil,
	}
	conf.workers = append(conf.workers, wrk)

	log.Printf("Worker connected: %s:%s", wrk.host, wrk.port)
	conf.broadcastWorkers()
	conf.ScheduleTasks()
	go wrk.ReadWorkerWebsocketMessage()

}

func (conf *appConfig) deleteWorker(w *worker) {
	conf.mu.Lock()
	conf.workers = slices.DeleteFunc(conf.workers, func(worker *worker) bool { return worker.id == w.id })
	conf.mu.Unlock()
	conf.broadcastWorkers()
	conf.ScheduleTasks()
}
