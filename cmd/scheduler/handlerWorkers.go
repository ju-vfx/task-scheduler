package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	type workerParams struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
	// reqParams, err := requests.DecodeRequest(req, workerParams{})
	// if err != nil {
	// 	requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
	// 	return
	// }

	ws, err := UpgradeConnection(w, req)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Could not connect to websocket")
		return
	}

	id, _ := uuid.NewUUID()
	wrk := &worker{
		id:          id,
		conn:        ws,
		host:        "",
		port:        "",
		connectedAt: time.Now(),
		lastSeenAt:  time.Now(),
		status:      utils.StatusWaiting,
		task:        "",
	}
	conf.workers = append(conf.workers, wrk)

	wrkResp := workerResp{
		ID:          wrk.id.String(),
		Host:        wrk.host,
		ConnectedAt: utils.TimeToString(wrk.connectedAt),
		LastSeenAt:  utils.TimeToString(wrk.lastSeenAt),
		Status:      utils.ObjectStatus(wrk.status).String(),
	}

	wrk.SendWsMessage(wsMessage{message: requests.EncodeJSON(wrkResp), messageType: websocket.BinaryMessage})

	go wrk.ReadWsMessage()
	fmt.Println(conf.workers)
}

func (conf *appConfig) handlerDeleteWorkers(w http.ResponseWriter, req *http.Request) {
	err := conf.db.DeleteWorkers(req.Context())
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Error deleting workers")
		return
	}
	requests.RespondWithJSON(w, http.StatusOK, "")
}
