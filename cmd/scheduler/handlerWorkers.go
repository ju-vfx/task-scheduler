package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

func (srv *server) handlerGetWorkers(w http.ResponseWriter, req *http.Request) {
	workers, err := srv.cfg.db.GetWorkers(req.Context())
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Error loading workers")
		return
	}
	type workerResp struct {
		ID          string `json:"id"`
		Host        string `json:"host"`
		ConnectedAt string `json:"connected_at"`
		LastSeenAt  string `json:"last_seen_at"`
		Status      string `json:"status"`
	}

	data := make([]workerResp, 0)

	for _, worker := range workers {
		w := workerResp{
			ID:          worker.ID.String(),
			Host:        worker.Host,
			ConnectedAt: worker.ConnectedAt.String(),
			LastSeenAt:  worker.LastSeenAt.String(),
			Status:      utils.ObjectStatus(worker.Status).String(),
		}

		data = append(data, w)
	}
	requests.RespondWithJSON(w, http.StatusOK, data)
}

func (srv *server) handlerRegisterWorker(w http.ResponseWriter, req *http.Request) {

	type workerParams struct {
		ID   *string `json:"id"`
		Host string  `json:"host"`
		Port string  `json:"port"`
	}
	reqParams, err := requests.DecodeRequest(req, workerParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	if reqParams.ID != nil {
		err = srv.cfg.db.UpdateLastSeen(req.Context(), uuid.MustParse(*reqParams.ID))
		if err != nil {
			requests.RespondWithError(w, http.StatusInternalServerError, "Error logging in worker")
			return
		}
	} else {
		worker, err := srv.cfg.db.CreateWorker(req.Context(), database.CreateWorkerParams{Host: reqParams.Host, Port: reqParams.Port, Status: int32(utils.StatusWaiting)})
		if err != nil {
			requests.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v+", err))
			return
		}

		type respParams struct {
			ID string `json:"id"`
		}
		requests.RespondWithJSON(w, http.StatusOK, respParams{ID: worker.ID.String()})
	}
}

func (srv *server) handlerDeleteWorkers(w http.ResponseWriter, req *http.Request) {
	err := srv.cfg.db.DeleteWorkers(req.Context())
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Error deleting workers")
		return
	}
	requests.RespondWithJSON(w, http.StatusOK, "")
}
