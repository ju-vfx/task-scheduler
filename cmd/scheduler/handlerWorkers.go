package main

import (
	"net/http"
	"slices"
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

func (conf *appConfig) handlerGetWorkers(w http.ResponseWriter, req *http.Request) {

	data := make([]workerResp, 0)

	for _, worker := range conf.workers {
		w := workerResp{
			ID:          worker.id.String(),
			Host:        worker.host,
			ConnectedAt: utils.TimeToString(worker.connectedAt),
			LastSeenAt:  utils.TimeToString(worker.lastSeenAt),
			Status:      utils.ObjectStatus(worker.status).String(),
		}

		data = append(data, w)
	}
	requests.RespondWithJSON(w, http.StatusOK, data)
}

func (conf *appConfig) handlerRegisterWorker(w http.ResponseWriter, req *http.Request) {

	ws, err := UpgradeConnection(w, req)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Could not connect to websocket")
		return
	}

	id, _ := uuid.NewUUID()
	wrk := &worker{
		conf:        conf,
		id:          id,
		conn:        ws,
		host:        "",
		port:        "",
		connectedAt: time.Now(),
		lastSeenAt:  time.Now(),
		status:      utils.StatusWaiting,
		task_id:     nil,
	}
	conf.workers = append(conf.workers, wrk)

	conf.UpdateState()
	go wrk.ReadWorkerWebsocketMessage()

}

func (conf *appConfig) deleteWorker(w *worker) {
	conf.mu.Lock()
	conf.workers = slices.DeleteFunc(conf.workers, func(worker *worker) bool { return worker.id == w.id })
	conf.mu.Unlock()
	conf.UpdateState()
}
